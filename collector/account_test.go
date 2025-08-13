package collector

import (
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountCollector_NewAccountCollector(t *testing.T) {
	logger := log.NewNopLogger()
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_errors_total",
		Help: "Test errors",
	}, []string{"collector"})
	client := &godo.Client{}
	timeout := 5 * time.Second

	collector := NewAccountCollector(logger, errors, client, timeout)

	assert.NotNil(t, collector)
	assert.Equal(t, logger, collector.logger)
	assert.Equal(t, errors, collector.errors)
	assert.Equal(t, client, collector.client)
	assert.Equal(t, timeout, collector.timeout)
	
	// Verify all metric descriptions are created
	assert.NotNil(t, collector.DropletLimit)
	assert.NotNil(t, collector.FloatingIPLimit)
	assert.NotNil(t, collector.EmailVerified)
	assert.NotNil(t, collector.Active)
}

func TestAccountCollector_Describe(t *testing.T) {
	logger := log.NewNopLogger()
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_errors_total",
		Help: "Test errors",
	}, []string{"collector"})
	client := &godo.Client{}
	timeout := 5 * time.Second

	collector := NewAccountCollector(logger, errors, client, timeout)
	
	descCh := make(chan *prometheus.Desc, 10)
	collector.Describe(descCh)
	close(descCh)

	var descriptions []*prometheus.Desc
	for desc := range descCh {
		descriptions = append(descriptions, desc)
	}

	require.Len(t, descriptions, 4)
	assert.Contains(t, descriptions, collector.DropletLimit)
	assert.Contains(t, descriptions, collector.FloatingIPLimit)
	assert.Contains(t, descriptions, collector.EmailVerified)
	assert.Contains(t, descriptions, collector.Active)
}