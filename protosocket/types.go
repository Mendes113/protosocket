package protosocket

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type HandlerFunc = func(ctx context.Context, msg proto.Message) error
type Middleware = func(next HandlerFunc) HandlerFunc
