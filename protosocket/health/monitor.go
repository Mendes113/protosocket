package health

import (
	"time"
)

type MetricsCollector interface {
	IncrementCounter(name string, value int64)
	SetGauge(name string, value float64)
	RecordHistogram(name string, value float64)
	AddLabel(key, value string)
}

type HealthMonitor struct {
	checks  []HealthCheck
	status  *Status
	alerter Alerter
	metrics MetricsCollector
}

type Status struct {
	Uptime      time.Duration
	Connections int
	MemoryUsage uint64
	CPUUsage    float64
	ErrorRate   float64
}

type Alerter interface {
	Alert(level string, message string, metadata map[string]interface{}) error
	ResolveAlert(alertID string) error
	GetActiveAlerts() ([]Alert, error)
}

type Alert struct {
	ID        string
	Level     string
	Message   string
	Timestamp time.Time
	Metadata  map[string]interface{}
	Resolved  bool
}
