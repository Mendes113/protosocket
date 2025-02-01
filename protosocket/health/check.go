package health

import (
	"context"
	"time"
)

type HealthCheck interface {
	Name() string
	Check(ctx context.Context) error
	Interval() time.Duration
}

type HealthCheckResult struct {
	Name      string
	Status    bool
	Error     error
	Timestamp time.Time
	Duration  time.Duration
}

type SystemHealthCheck struct {
	name     string
	interval time.Duration
	check    func(context.Context) error
}

func NewSystemHealthCheck(name string, interval time.Duration, check func(context.Context) error) *SystemHealthCheck {
	return &SystemHealthCheck{
		name:     name,
		interval: interval,
		check:    check,
	}
}

func (hc *SystemHealthCheck) Name() string {
	return hc.name
}

func (hc *SystemHealthCheck) Check(ctx context.Context) error {
	return hc.check(ctx)
}

func (hc *SystemHealthCheck) Interval() time.Duration {
	return hc.interval
}
