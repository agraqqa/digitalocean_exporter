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

func TestDomainCollector_NewDomainCollector(t *testing.T) {
	logger := log.NewNopLogger()
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_errors_total",
		Help: "Test errors",
	}, []string{"collector"})
	client := &godo.Client{}
	timeout := 5 * time.Second

	collector := NewDomainCollector(logger, errors, client, timeout)

	assert.NotNil(t, collector)
	assert.Equal(t, logger, collector.logger)
	assert.Equal(t, errors, collector.errors)
	assert.Equal(t, client, collector.client)
	assert.Equal(t, timeout, collector.timeout)
	
	// Verify metric descriptions are created
	assert.NotNil(t, collector.DomainRecordPort)
	assert.NotNil(t, collector.DomainRecordPriority)
	assert.NotNil(t, collector.DomainRecordWeight)
	assert.NotNil(t, collector.DomainTTL)
}

func TestDomainCollector_Describe(t *testing.T) {
	logger := log.NewNopLogger()
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_errors_total",
		Help: "Test errors",
	}, []string{"collector"})
	client := &godo.Client{}
	timeout := 5 * time.Second

	collector := NewDomainCollector(logger, errors, client, timeout)
	
	descCh := make(chan *prometheus.Desc, 10)
	collector.Describe(descCh)
	close(descCh)

	var descriptions []*prometheus.Desc
	for desc := range descCh {
		descriptions = append(descriptions, desc)
	}

	require.Len(t, descriptions, 4)
	assert.Contains(t, descriptions, collector.DomainRecordPort)
	assert.Contains(t, descriptions, collector.DomainRecordPriority)
	assert.Contains(t, descriptions, collector.DomainRecordWeight)
	assert.Contains(t, descriptions, collector.DomainTTL)
}