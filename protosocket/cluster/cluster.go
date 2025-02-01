package cluster

import (
	"time"
)

type NodeRole int
type NodeStatus int

const (
	RolePrimary NodeRole = iota
	RoleSecondary
	RoleObserver
)

const (
	StatusHealthy NodeStatus = iota
	StatusDegraded
	StatusUnhealthy
)

type ClusterNode struct {
	ID          string
	Address     string
	Role        NodeRole
	Status      NodeStatus
	Connections int
	Load        float64
}

type ClusterManager struct {
	nodes       map[string]*ClusterNode
	coordinator *Coordinator
	balancer    *LoadBalancer
}

type Coordinator struct {
	LeaderID   string
	Term       int64
	LastUpdate time.Time
}

type LoadBalancer struct {
	Strategy     string
	MaxLoad      float64
	Distribution map[string]float64
}
