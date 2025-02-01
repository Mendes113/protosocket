package observability

import (
	"time"

	"go.uber.org/zap"
)

type Span struct {
	ID        string
	TraceID   string
	ParentID  string
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Tags      map[string]string
	Events    []SpanEvent
}

type SpanEvent struct {
	Time    time.Time
	Name    string
	Details map[string]interface{}
}

type TraceExporter interface {
	ExportSpan(span *Span) error
}

type Tracer struct {
	spans    map[string]*Span
	metrics  *Metrics
	logger   *zap.Logger
	exporter TraceExporter
}

type MetricsExporter struct {
	prometheus *PrometheusExporter
	statsd     *StatsDExporter
	custom     []MetricsExporter
}

type Metrics struct {
	Counters   map[string]int64
	Gauges     map[string]float64
	Histograms map[string][]float64
	Labels     map[string]string
}

type PrometheusExporter struct {
	Namespace string
	Registry  interface{} // prometheus.Registry
}

type StatsDExporter struct {
	Prefix string
	Client interface{} // statsd.Client
}
