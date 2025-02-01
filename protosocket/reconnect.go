package protosocket

import (
	"time"
)

type ReconnectConfig struct {
	MaxAttempts       int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
}

func (p *Peer) EnableReconnect(config ReconnectConfig) {
	go func() {
		attempts := 0
		delay := config.InitialDelay

		for {
			if !p.IsConnected() && attempts < config.MaxAttempts {
				time.Sleep(delay)

				if err := p.Connect("localhost:8081"); err != nil {
					attempts++
					delay = time.Duration(float64(delay) * config.BackoffMultiplier)
					if delay > config.MaxDelay {
						delay = config.MaxDelay
					}
					continue
				}

				attempts = 0
				delay = config.InitialDelay
			}
			time.Sleep(time.Second)
		}
	}()
}
