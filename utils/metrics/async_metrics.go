package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type AsyncMetrics struct {
	mux          sync.RWMutex
	wg           sync.WaitGroup
	blockMetrics bool

	metrics        chan MetricData
	metricStore    sync.Map
	customRegistry prometheus.Registerer // Custom Prometheus registry
}

func NewPrometheusMetrics(ctx context.Context) *AsyncMetrics {
	collector := &AsyncMetrics{
		metrics:        make(chan MetricData, 256),
		customRegistry: prometheus.NewRegistry(), // initialize a new registry
	}

	prometheus.MustRegister(collector)
	go collector.handleMetrics(ctx)
	return collector
}

func (col *AsyncMetrics) Emit(metric *MetricDefinition, value float64) error {
	col.withMetricsNotBlocked(func() {
		col.metrics <- MetricData{MetricDefinition: metric, Value: value}
	})
	return nil
}

// Describe sends the super-set of all possible descriptors of metrics
// produced by this Collector to the provided channel.
func (col *AsyncMetrics) Describe(ch chan<- *prometheus.Desc) {
	col.metricStore.Range(func(k, v interface{}) bool {
		if collector, ok := v.(prometheus.Collector); ok {
			collector.Describe(ch)
		}
		return true
	})
}

// Collect is called by the Prometheus registry when collecting metrics.
func (col *AsyncMetrics) Collect(ch chan<- prometheus.Metric) {
	col.metricStore.Range(func(k, v interface{}) bool {
		if collector, ok := v.(prometheus.Collector); ok {
			collector.Collect(ch)
		}
		return true
	})
}

func (col *AsyncMetrics) Flush() {
	for {
		select {
		case metric := <-col.metrics:
			col.writeMetric(metric)
		default:
			return
		}
	}
}

func (col *AsyncMetrics) handleMetrics(ctx context.Context) {
	for {
		select {
		case metric := <-col.metrics:
			col.writeMetric(metric)
		case <-ctx.Done():
			col.mux.Lock()
			defer col.mux.Unlock()

			col.Flush()
			col.blockMetrics = true
			close(col.metrics)
			return
		}
	}
}

func (col *AsyncMetrics) withMetricsNotBlocked(f func()) {
	col.mux.RLock()
	defer col.mux.RUnlock()

	if col.blockMetrics {
		return
	}

	f()
}

func (col *AsyncMetrics) writeMetric(metric MetricData) {
	metricKey := metric.Namespace + ":" + metric.Name
	metricInterface, ok := col.metricStore.Load(metricKey)

	if !ok {
		// If the metric doesn't exist, create it
		switch metric.MetricType {
		case Counter:
			metricInterface = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: metric.Namespace,
					Name:      metric.Name,
					Help:      metric.Help,
				},
				metric.LabelNames,
			)
		case Gauge:
			metricInterface = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: metric.Namespace,
					Name:      metric.Name,
					Help:      metric.Help,
				},
				metric.LabelNames,
			)
		case Histogram:
			metricInterface = prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Namespace: metric.Namespace,
					Name:      metric.Name,
					Help:      metric.Help,
					Buckets:   metric.Buckets,
				},
				metric.LabelNames,
			)
		case Summary:
			metricInterface = prometheus.NewSummaryVec(
				prometheus.SummaryOpts{
					Namespace:  metric.Namespace,
					Name:       metric.Name,
					Help:       metric.Help,
					Objectives: metric.Quantiles,
				},
				metric.LabelNames,
			)
		default:
			col.mux.Unlock()
			return // Ignore if it's an unknown type
		}

		// Try to register the new metric and store it only if the registration is successful
		err := col.customRegistry.Register(metricInterface.(prometheus.Collector))
		if err == nil {
			col.metricStore.Store(metricKey, metricInterface)
		}
	}

	// Assert the type of the metric and update it
	switch metric.MetricType {
	case Counter:
		if m, ok := metricInterface.(*prometheus.CounterVec); ok {
			m.WithLabelValues(metric.LabelValues...).Add(metric.Value)
		}
	case Gauge:
		if m, ok := metricInterface.(*prometheus.GaugeVec); ok {
			m.WithLabelValues(metric.LabelValues...).Set(metric.Value)
		}
	case Histogram:
		if m, ok := metricInterface.(*prometheus.HistogramVec); ok {
			m.WithLabelValues(metric.LabelValues...).Observe(metric.Value)
		}
	case Summary:
		if m, ok := metricInterface.(*prometheus.SummaryVec); ok {
			m.WithLabelValues(metric.LabelValues...).Observe(metric.Value)
		}
	}
}

func (col *AsyncMetrics) StartHTTPServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic("Could not start HTTP server: " + err.Error())
		}
	}()
}
