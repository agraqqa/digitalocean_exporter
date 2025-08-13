package collector

import (
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExporterCollector_NewExporterCollector(t *testing.T) {
	logger := log.NewNopLogger()
	version := "1.0.0"
	revision := "abc123"
	buildDate := "2024-01-01"
	goVersion := "go1.21"
	startTime := time.Now()

	collector := NewExporterCollector(logger, version, revision, buildDate, goVersion, startTime)

	assert.NotNil(t, collector)
	assert.Equal(t, version, collector.version)
	assert.Equal(t, revision, collector.revision)
	assert.Equal(t, buildDate, collector.buildDate)
	assert.Equal(t, goVersion, collector.goVersion)
	assert.Equal(t, startTime, collector.startTime)

	// Verify all metric descriptions are created
	assert.NotNil(t, collector.StartTime)
	assert.NotNil(t, collector.BuildInfo)
}

func TestExporterCollector_Describe(t *testing.T) {
	logger := log.NewNopLogger()
	startTime := time.Now()
	collector := NewExporterCollector(logger, "1.0.0", "abc123", "2024-01-01", "go1.21", startTime)
	
	descCh := make(chan *prometheus.Desc, 10)
	collector.Describe(descCh)
	close(descCh)

	var descriptions []*prometheus.Desc
	for desc := range descCh {
		descriptions = append(descriptions, desc)
	}

	require.Len(t, descriptions, 1)
	assert.Contains(t, descriptions, collector.StartTime)
}