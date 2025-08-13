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

func TestBalanceCollector_Describe(t *testing.T) {
	logger := log.NewNopLogger()
	errorsCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_errors_total",
		Help: "Test errors",
	}, []string{"collector"})

	collector := &BalanceCollector{
		logger:  logger,
		errors:  errorsCounter,
		timeout: 5 * time.Second,
		MonthToDateBalance: prometheus.NewDesc(
			"digitalocean_month_to_date_balance",
			"Balance as of the digitalocean_balance_generated_at time",
			nil, nil,
		),
		AccountBalance: prometheus.NewDesc(
			"digitalocean_account_balance", 
			"Current balance of your most recent billing activity",
			nil, nil,
		),
		MonthToDateUsage: prometheus.NewDesc(
			"digitalocean_month_to_date_usage",
			"Amount used in the current billing period as of the digitalocean_balance_generated_at time", 
			nil, nil,
		),
		BalanceGeneratedAt: prometheus.NewDesc(
			"digitalocean_balance_generated_at",
			"The time at which balances were most recently generated",
			nil, nil,
		),
	}

	descCh := make(chan *prometheus.Desc, 10)
	collector.Describe(descCh)
	close(descCh)

	var descriptions []*prometheus.Desc
	for desc := range descCh {
		descriptions = append(descriptions, desc)
	}

	require.Len(t, descriptions, 4)
	assert.Contains(t, descriptions, collector.MonthToDateBalance)
	assert.Contains(t, descriptions, collector.AccountBalance)
	assert.Contains(t, descriptions, collector.MonthToDateUsage)
	assert.Contains(t, descriptions, collector.BalanceGeneratedAt)
}

func TestBalanceCollector_NewBalanceCollector(t *testing.T) {
	logger := log.NewNopLogger()
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_errors_total",
		Help: "Test errors",
	}, []string{"collector"})
	client := &godo.Client{}
	timeout := 5 * time.Second

	collector := NewBalanceCollector(logger, errors, client, timeout)

	assert.NotNil(t, collector)
	assert.Equal(t, logger, collector.logger)
	assert.Equal(t, errors, collector.errors)
	assert.Equal(t, client, collector.client)
	assert.Equal(t, timeout, collector.timeout)
	
	// Verify all metric descriptions are created
	assert.NotNil(t, collector.MonthToDateBalance)
	assert.NotNil(t, collector.AccountBalance)
	assert.NotNil(t, collector.MonthToDateUsage)
	assert.NotNil(t, collector.BalanceGeneratedAt)
}