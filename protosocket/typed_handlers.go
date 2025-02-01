package protosocket

import (
	"google.golang.org/protobuf/proto"
)

type EventHandler[T proto.Message] struct {
	Handler func(data T, socket *Socket)
}

func OnTyped[T proto.Message](s *Socket, event string, handler EventHandler[T]) {
	s.On(event, func(data proto.Message, socket *Socket) {
		if msg, ok := data.(T); ok {
			handler.Handler(msg, socket)
		}
	})
}
