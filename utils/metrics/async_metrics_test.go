package metrics

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestComellect(t *testing.T) {
	ctx := context.Background()

	collector := NewPrometheusMetrics(ctx)

	// Test collecting a Counter metric
	counterMetric := &MetricDefinition{
		MetricType:  Counter,
		Namespace:   "test",
		Name:        "counter_metric",
		Help:        "This is a test counter metric",
		LabelNames:  []string{"label1", "label2"},
		LabelValues: []string{"value1", "value2"},
	}
	err := collector.Emit(counterMetric, 1.0)

	if err != nil {
		t.Errorf("Collecting metric failed: %s", err.Error())
		return
	}

	// Test collecting a Gauge metric
	gaugeMetric := &MetricDefinition{
		MetricType:  Gauge,
		Namespace:   "test",
		Name:        "gauge_metric",
		Help:        "This is a test gauge metric",
		LabelNames:  []string{"label1", "label2"},
		LabelValues: []string{"value1", "value2"},
	}
	err = collector.Emit(gaugeMetric, 2.5)

	if err != nil {
		t.Errorf("Collecting metric failed: %s", err.Error())
		return
	}

	// Test collecting a Histogram metric
	histogramMetric := &MetricDefinition{
		MetricType:  Histogram,
		Namespace:   "test",
		Name:        "histogram_metric",
		Help:        "This is a test histogram metric",
		LabelNames:  []string{"label1", "label2"},
		LabelValues: []string{"value1", "value2"},
		Buckets:     []float64{1, 2, 5, 10},
	}
	err = collector.Emit(histogramMetric, 4.0)

	if err != nil {
		t.Errorf("Collecting metric failed: %s", err.Error())
		return
	}

	// Test collecting a Summary metric
	summaryMetric := &MetricDefinition{
		MetricType:  Summary,
		Namespace:   "test",
		Name:        "summary_metric",
		Help:        "This is a test summary metric",
		LabelNames:  []string{"label1", "label2"},
		LabelValues: []string{"value1", "value2"},
		Quantiles:   map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}
	err = collector.Emit(summaryMetric, 0.8)

	if err != nil {
		t.Errorf("Collecting metric failed: %s", err.Error())
		return
	}

	// Sleep for a short period to ensure the collected metrics are processed
	time.Sleep(100 * time.Millisecond)

	// Use testutil package from Prometheus client library to check metric values
	if err = testutil.CollectAndCompare(collector, strings.NewReader(`
		# HELP test_counter_metric This is a test counter metric
		# TYPE test_counter_metric counter
		test_counter_metric{label1="value1",label2="value2"} 1

		# HELP test_gauge_metric This is a test gauge metric
		# TYPE test_gauge_metric gauge
		test_gauge_metric{label1="value1",label2="value2"} 2.5

		# HELP test_histogram_metric This is a test histogram metric
		# TYPE test_histogram_metric histogram
		test_histogram_metric_bucket{label1="value1",label2="value2",le="1"} 0
		test_histogram_metric_bucket{label1="value1",label2="value2",le="2"} 0
		test_histogram_metric_bucket{label1="value1",label2="value2",le="5"} 1
		test_histogram_metric_bucket{label1="value1",label2="value2",le="10"} 1
		test_histogram_metric_bucket{label1="value1",label2="value2",le="+Inf"} 1
		test_histogram_metric_sum{label1="value1",label2="value2"} 4
		test_histogram_metric_count{label1="value1",label2="value2"} 1

		# HELP test_summary_metric This is a test summary metric
		# TYPE test_summary_metric summary
		test_summary_metric{label1="value1",label2="value2",quantile="0.5"} 0.8
		test_summary_metric{label1="value1",label2="value2",quantile="0.9"} 0.8
		test_summary_metric{label1="value1",label2="value2",quantile="0.99"} 0.8
		test_summary_metric_sum{label1="value1",label2="value2"} 0.8
		test_summary_metric_count{label1="value1",label2="value2"} 1
	`), "test_counter_metric", "test_gauge_metric", "test_histogram_metric", "test_summary_metric"); err != nil {
		t.Errorf("Metrics do not match expected values: %s", err.Error())
	}
}

func BenchmarkCounterMetric(b *testing.B) {
	ctx := context.Background()
	collector := GetOrCreateGlobalMetrics(ctx)

	counterMetric := &MetricDefinition{
		MetricType:  Counter,
		Namespace:   "test",
		Name:        "counter",
		Help:        "counter help",
		LabelNames:  []string{"label1", "label2"},
		LabelValues: []string{"value1", "value2"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		collector.Emit(counterMetric, 1)
	}
}

func BenchmarkGaugeMetric(b *testing.B) {
	ctx := context.Background()
	collector := GetOrCreateGlobalMetrics(ctx)

	gaugeMetric := &MetricDefinition{
		MetricType:  Gauge,
		Namespace:   "test",
		Name:        "gauge",
		Help:        "gauge help",
		LabelNames:  []string{"label1", "label2"},
		LabelValues: []string{"value1", "value2"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		collector.Emit(gaugeMetric, 1)
	}
}

func BenchmarkHistogramMetric(b *testing.B) {
	ctx := context.Background()
	collector := GetOrCreateGlobalMetrics(ctx)

	histogramMetric := &MetricDefinition{
		MetricType:  Histogram,
		Namespace:   "test",
		Name:        "histogram",
		Help:        "histogram help",
		LabelNames:  []string{"label1", "label2"},
		LabelValues: []string{"value1", "value2"},
		Buckets:     []float64{1, 2, 5, 10, 20, 50},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		collector.Emit(histogramMetric, 1)
	}
}

func BenchmarkSummaryMetric(b *testing.B) {
	ctx := context.Background()
	collector := GetOrCreateGlobalMetrics(ctx)

	summaryMetric := &MetricDefinition{
		MetricType:  Summary,
		Namespace:   "test",
		Name:        "summary",
		Help:        "summary help",
		LabelNames:  []string{"label1", "label2"},
		LabelValues: []string{"value1", "value2"},
		Quantiles:   map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		collector.Emit(summaryMetric, 1)
	}
}
