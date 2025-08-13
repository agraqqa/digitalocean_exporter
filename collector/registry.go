package collector

import (
	"context"
	"time"

	"github.com/digitalocean/godo"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// RegistryCollector collects metrics about container registries.
type RegistryCollector struct {
	logger  log.Logger
	errors  *prometheus.CounterVec
	client  *godo.Client
	timeout time.Duration

	IncludedStorageBytes *prometheus.Desc
	StorageUsageBytes    *prometheus.Desc
	StorageUsageUpdated  *prometheus.Desc
}

// NewRegistryCollector returns a new RegistryCollector.
func NewRegistryCollector(logger log.Logger, errors *prometheus.CounterVec, client *godo.Client, timeout time.Duration) *RegistryCollector {
	errors.WithLabelValues("registry").Add(0)

	labels := []string{"registry_name", "region"}
	return &RegistryCollector{
		logger:  logger,
		errors:  errors,
		client:  client,
		timeout: timeout,

		IncludedStorageBytes: prometheus.NewDesc(
			"digitalocean_registry_included_storage_bytes",
			"Registry's included storage in bytes",
			labels, nil,
		),
		StorageUsageBytes: prometheus.NewDesc(
			"digitalocean_registry_storage_usage_bytes",
			"Registry's current storage usage in bytes",
			labels, nil,
		),
		StorageUsageUpdated: prometheus.NewDesc(
			"digitalocean_registry_storage_usage_updated",
			"Unix timestamp when storage usage was last updated",
			labels, nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector.
func (c *RegistryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.IncludedStorageBytes
	ch <- c.StorageUsageBytes
	ch <- c.StorageUsageUpdated
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *RegistryCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Get registry information - there's typically only one registry per account
	registry, _, err := c.client.Registry.Get(ctx)
	if err != nil {
		c.errors.WithLabelValues("registry").Add(1)
		level.Debug(c.logger).Log(
			"msg", "can't get registry, likely no registry configured",
			"err", err,
		)
		return
	}

	// Get subscription information
	subscription, _, err := c.client.Registry.GetSubscription(ctx)
	if err != nil {
		c.errors.WithLabelValues("registry").Add(1)
		level.Warn(c.logger).Log(
			"msg", "can't get registry subscription",
			"registry", registry.Name,
			"err", err,
		)
		return
	}

	labels := []string{
		registry.Name,
		"",
	}

	// Included storage from subscription tier
	ch <- prometheus.MustNewConstMetric(
		c.IncludedStorageBytes,
		prometheus.GaugeValue,
		float64(subscription.Tier.IncludedStorageBytes),
		labels...,
	)

	// Current storage usage from registry
	ch <- prometheus.MustNewConstMetric(
		c.StorageUsageBytes,
		prometheus.GaugeValue,
		float64(registry.StorageUsageBytes),
		labels...,
	)

	// When storage usage was last updated
	ch <- prometheus.MustNewConstMetric(
		c.StorageUsageUpdated,
		prometheus.GaugeValue,
		float64(registry.StorageUsageBytesUpdatedAt.Unix()),
		labels...,
	)
}
