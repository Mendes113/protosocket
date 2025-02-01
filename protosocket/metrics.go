package protosocket

import (
	"sync"
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

