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

func TestDropletCollector_NewDropletCollector(t *testing.T) {
	logger := log.NewNopLogger()
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_errors_total",
		Help: "Test errors",
	}, []string{"collector"})
	client := &godo.Client{}
	timeout := 5 * time.Second

	collector := NewDropletCollector(logger, errors, client, timeout)

	assert.NotNil(t, collector)
	assert.Equal(t, logger, collector.logger)
	assert.Equal(t, errors, collector.errors)
	assert.Equal(t, client, collector.client)
	assert.Equal(t, timeout, collector.timeout)
	
	// Verify all metric descriptions are created
	assert.NotNil(t, collector.Up)
	assert.NotNil(t, collector.CPUs)
	assert.NotNil(t, collector.Memory)
	assert.NotNil(t, collector.Disk)
	assert.NotNil(t, collector.PriceHourly)
	assert.NotNil(t, collector.PriceMonthly)
}

func TestDropletCollector_Describe(t *testing.T) {
	logger := log.NewNopLogger()
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_errors_total",
		Help: "Test errors",
	}, []string{"collector"})
	client := &godo.Client{}
	timeout := 5 * time.Second

	collector := NewDropletCollector(logger, errors, client, timeout)
	
	descCh := make(chan *prometheus.Desc, 10)
	collector.Describe(descCh)
	close(descCh)

	var descriptions []*prometheus.Desc
	for desc := range descCh {
		descriptions = append(descriptions, desc)
	}

	require.Len(t, descriptions, 6)
	assert.Contains(t, descriptions, collector.Up)
	assert.Contains(t, descriptions, collector.CPUs)
	assert.Contains(t, descriptions, collector.Memory)
	assert.Contains(t, descriptions, collector.Disk)
	assert.Contains(t, descriptions, collector.PriceHourly)
	assert.Contains(t, descriptions, collector.PriceMonthly)
}