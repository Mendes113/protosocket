package protosocket

import (
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	TotalConnections  int
	ActiveConnections int
	MessagesReceived  int
	MessagesSent      int
	StartTime         time.Time
	mutex             sync.RWMutex
}

type MetricsCollector struct {
	messagesSent     uint64
	messagesReceived uint64
	bytesTransferred uint64
	errors           uint64
	latencies        []time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{
		StartTime: time.Now(),
	}
}

func (m *Metrics) ConnectionOpened() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.TotalConnections++
	m.ActiveConnections++
}

func (m *Metrics) ConnectionClosed() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ActiveConnections--
}

func (m *MetricsCollector) RecordMessage(size int, latency time.Duration) {
	atomic.AddUint64(&m.messagesSent, 1)
	atomic.AddUint64(&m.bytesTransferred, uint64(size))
	m.latencies = append(m.latencies, latency)
}

func (m *MetricsCollector) RecordError() {
	atomic.AddUint64(&m.errors, 1)
}
