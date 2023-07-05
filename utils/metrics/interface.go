package metrics

import (
	"context"
	"sync"
)

// MetricType represents the possible types of metrics (Counter, Gauge, Summary, Histogram)
type MetricType uint

const (
	// Counter is a cumulative metric that represents a single numerical value that only ever goes up
	Counter = iota

	// Gauge is a metric that represents a single numerical value that can arbitrarily go up and down
	Gauge

	// Summary samples observations (usually things like request durations) and provides configurable quantiles
	Summary

	// Histogram samples observations and counts them in configurable buckets
	Histogram
)

type MetricDefinition struct {
	MetricType  MetricType
	Namespace   string
	Name        string
	Help        string
	LabelNames  []string
	LabelValues []string
	Buckets     []float64           // Only used for Histograms
	Quantiles   map[float64]float64 // Only used for Summaries
}

type MetricData struct {
	*MetricDefinition
	Value float64
}

// Metrics represents a standard metrics interface
type Metrics interface {
	Emit(metric *MetricDefinition, value float64) error
	Flush()
}

var (
	globalMetrics Metrics
	once          sync.Once
)

func GetOrCreateGlobalMetrics(ctx context.Context) Metrics {
	once.Do(func() {
		globalMetrics = NewPrometheusMetrics(ctx)
	})
	return globalMetrics
}

func GetGlobalMetrics() Metrics {
	return globalMetrics
}
