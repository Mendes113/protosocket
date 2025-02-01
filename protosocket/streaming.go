package protosocket

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Stream struct {
	ctx     context.Context
	socket  *Socket
	event   string
	msgType proto.Message
}

func (s *Stream) Send(msg proto.Message) error {
	return s.socket.Emit(s.event, msg)
}

func (s *Stream) Recv() (proto.Message, error) {
	// Implementar lógica de recebimento de stream
	return nil, nil
}

func (s *Socket) CreateStream(ctx context.Context, event string, msgType proto.Message) *Stream {
	return &Stream{
		ctx:     ctx,
		socket:  s,
		event:   event,
		msgType: msgType,
	}
}
