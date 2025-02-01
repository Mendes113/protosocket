package protosocket

import (
	"context"
	"errors"
	"sync/atomic"
)

var (
	ErrTooManyRequests = errors.New("too many requests")
)

type Request interface {
	Process(ctx context.Context) error
}

type LoadController struct {
	activeRequests int64
	maxRequests    int64
	queueSize      int
	requestQueue   chan Request
}

func (lc *LoadController) HandleRequest(ctx context.Context, req Request) error {
	if atomic.LoadInt64(&lc.activeRequests) >= lc.maxRequests {
		select {
		case lc.requestQueue <- req:
			// Request enfileirado
		default:
			return ErrTooManyRequests
		}
	}

	atomic.AddInt64(&lc.activeRequests, 1)
	defer atomic.AddInt64(&lc.activeRequests, -1)

	return req.Process(ctx)
}
