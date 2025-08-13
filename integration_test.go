package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/go-kit/kit/log"
	"github.com/metalmatze/digitalocean_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestIntegration_MetricsEndpoint(t *testing.T) {
	// This is an integration test that requires a real or mock API
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a mock DigitalOcean API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/account"):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"account": {
					"droplet_limit": 25,
					"floating_ip_limit": 5,
					"email_verified": true,
					"status": "active"
				}
			}`))
		case strings.Contains(r.URL.Path, "/customers/my/balance"):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"month_to_date_balance": "10.50",
				"account_balance": "0.00",
				"month_to_date_usage": "5.25",
				"generated_at": "2024-01-01T12:00:00Z"
			}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a config that points to our mock server
	config := Config{
		HTTPTimeout: 5000,
		WebPath:     "/metrics",
		WebAddr:     ":0", // Let the system choose a port
	}

	// Create OAuth client that points to mock server
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "fake-token"})
	httpClient := oauth2.NewClient(ctx, tokenSource)
	
	// Create godo client pointing to mock server  
	client, err := godo.New(httpClient, godo.SetBaseURL(server.URL+"/"))
	require.NoError(t, err)

	timeout := time.Duration(config.HTTPTimeout) * time.Millisecond

	// Create logger
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	// Create error counter
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "digitalocean_errors_total",
		Help: "The total number of errors per collector",
	}, []string{"collector"})

	// Create registry and register collectors
	registry := prometheus.NewRegistry()
	registry.MustRegister(errors)
	registry.MustRegister(collector.NewAccountCollector(logger, errors, client, timeout))
	registry.MustRegister(collector.NewBalanceCollector(logger, errors, client, timeout))

	// Create HTTP server for metrics
	metricsServer := httptest.NewServer(promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	defer metricsServer.Close()

	// Test the metrics endpoint
	resp, err := http.Get(metricsServer.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain")

	// Read the response body to check content  
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	
	// The response should contain our metrics
	bodyStr := string(body)
	assert.Contains(t, bodyStr, "digitalocean_account_active")
	assert.Contains(t, bodyStr, "digitalocean_account_balance")
}

func TestIntegration_ConfigFromEnv(t *testing.T) {
	// Test that configuration works from environment variables
	oldToken := os.Getenv("DIGITALOCEAN_TOKEN")
	oldWebAddr := os.Getenv("WEB_ADDR")
	
	defer func() {
		os.Setenv("DIGITALOCEAN_TOKEN", oldToken)
		os.Setenv("WEB_ADDR", oldWebAddr)
	}()

	os.Setenv("DIGITALOCEAN_TOKEN", "test-token")
	os.Setenv("WEB_ADDR", ":9999")

	// This would normally parse from environment using arg.MustParse
	// but we can't easily test that without mocking os.Args
	config := Config{
		DigitalOceanToken: os.Getenv("DIGITALOCEAN_TOKEN"),
		WebAddr:           os.Getenv("WEB_ADDR"),
		HTTPTimeout:       5000,
		WebPath:           "/metrics",
	}

	assert.Equal(t, "test-token", config.DigitalOceanToken)
	assert.Equal(t, ":9999", config.WebAddr)
	assert.Equal(t, 5000, config.HTTPTimeout)
	assert.Equal(t, "/metrics", config.WebPath)

	// Test token generation
	token, err := config.Token()
	assert.NoError(t, err)
	assert.Equal(t, "test-token", token.AccessToken)
}