package ratelimit

import (
	"context"
	"time"
)

type Request struct {
	IP       string
	UserID   string
	Resource string
	Method   string
}

type RateLimiter struct {
	windowSize  time.Duration
	maxRequests int
	buckets     map[string]*TokenBucket
	perIP       bool
	perUser     bool
	globalLimit bool
}

func (rl *RateLimiter) Allow(ctx context.Context, req *Request) bool {
	// Implementa rate limiting por IP, usuário e global
	return false
}
