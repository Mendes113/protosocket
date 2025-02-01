package protosocket

import (
	"time"
)

type PeerState struct {
	ID            string
	Connections   map[string]*ConnectionState
	LastSyncTime  time.Time
	Configuration map[string]interface{}
}

type ConnectionState struct {
	Connected    bool
	LastPing     time.Time
	MessageCount int64
	ErrorCount   int64
}

func (p *Peer) SaveState() error {
	state := &PeerState{
		ID:           p.ID,
		Connections:  make(map[string]*ConnectionState),
		LastSyncTime: time.Now(),
	}

	p.lock.RLock()
	for id := range p.clients {
		state.Connections[id] = &ConnectionState{
			Connected:    true,
			LastPing:     time.Now(),
			MessageCount: 0,
			ErrorCount:   0,
		}
	}
	p.lock.RUnlock()

	return nil
}

func (p *Peer) RecoverFromState() error {
	// Recupera estado após reinicialização
	return nil
}
