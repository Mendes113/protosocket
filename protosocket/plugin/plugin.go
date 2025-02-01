package plugin

import (
	"context"

	"github.com/mendes113/protosocket/protosocket"
	"github.com/mendes113/protosocket/protosocket/types"
)

type Plugin interface {
	Name() string
	Init(ctx context.Context) error
	OnMessage(msg *types.Message) error
	OnConnect(client *protosocket.Client) error
	OnDisconnect(client *protosocket.Client) error
	Shutdown(ctx context.Context) error
}
