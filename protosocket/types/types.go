package types

import (
	"time"
)

// Mensagem base
type Message struct {
	ID        string
	Type      string
	Data      []byte
	Metadata  map[string]string
	Timestamp time.Time
}

// Estados e Status
type ConnectionState int
type NodeStatus int
type NodeRole int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateReconnecting
)

const (
	StatusHealthy NodeStatus = iota
	StatusDegraded
	StatusUnhealthy
)

const (
	RolePrimary NodeRole = iota
	RoleSecondary
	RoleObserver
)

// Estatísticas
type ConnectionStats struct {
	BytesSent     uint64
	BytesReceived uint64
	MessagesSent  uint64
	MsgsReceived  uint64
	LastError     error
	LastActivity  time.Time
}
