package protosocket

import (
	"time"
)

type RetryConfig struct {
	MaxAttempts       int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
}

func WithRetry(config RetryConfig, operation func() error) error {
	var err error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		err = operation()
		if err == nil {
			return nil
		}

		if attempt < config.MaxAttempts-1 {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * config.BackoffMultiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}
	}

	return err
}
