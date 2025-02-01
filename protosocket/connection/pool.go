package connection

import (
	"time"
)

type ConnectionPool struct {
	active      map[string]*Connection
	idle        []*Connection
	maxIdle     int
	maxActive   int
	idleTimeout time.Duration
}

type Connection struct {
	ID         string
	State      ConnectionState
	Stats      *ConnectionStats
	LastActive time.Time
}

type ConnectionState int

const (
	ConnectionStateActive ConnectionState = iota
	ConnectionStateIdle
	ConnectionStateClosing
	ConnectionStateClosed
)

type ConnectionStats struct {
	TotalConnections int
	ActiveConnections int
	IdleConnections   int
	ClosingConnections int
	ClosedConnections int
}
