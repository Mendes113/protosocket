package protosocket

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

type CircuitBreaker struct {
	failures    int
	threshold   int
	timeout     time.Duration
	lastFailure time.Time
	isOpen      bool
	mutex       sync.RWMutex
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
	cb.mutex.Lock()
	if cb.isOpen {
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.isOpen = false
			cb.failures = 0
		} else {
			cb.mutex.Unlock()
			return ErrCircuitOpen
		}
	}
	cb.mutex.Unlock()

	err := operation()
	if err != nil {
		cb.recordFailure()
		return err
	}

	return nil
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()
	if cb.failures >= cb.threshold {
		cb.isOpen = true
	}
}
