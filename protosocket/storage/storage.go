package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mendes113/protosocket/protosocket/types"
)

// Message representa uma mensagem armazenada
type Message struct {
	ID        string
	Type      string
	Data      []byte
	Metadata  map[string]string
	Timestamp time.Time
}

type MessageStore interface {
	Save(ctx context.Context, msg *types.Message) error
	GetByID(ctx context.Context, id string) (*types.Message, error)
	GetByTimeRange(ctx context.Context, start, end time.Time) ([]*types.Message, error)
	DeleteOlderThan(ctx context.Context, age time.Duration) error
}

// Implementação com Redis
type RedisMessageStore struct {
	client *redis.Client
	ttl    time.Duration
}
