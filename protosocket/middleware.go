package protosocket

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type SocketMiddleware struct {
	middlewares []Middleware
	socket      *Socket
}

func NewMiddleware(socket *Socket) *SocketMiddleware {
	return &SocketMiddleware{
		socket: socket,
	}
}

func (m *SocketMiddleware) Use(middleware ...Middleware) {
	m.middlewares = append(m.middlewares, middleware...)
}

// Exemplo de middleware de logging
func LoggingMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, msg proto.Message) error {
			logger := GetLogger()

			logger.Info("processando mensagem",
				zap.String("type", fmt.Sprintf("%T", msg)),
				zap.Any("context", ctx))

			err := next(ctx, msg)
			if err != nil {
				logger.Error("erro no processamento",
					zap.Error(err),
					zap.String("type", fmt.Sprintf("%T", msg)))
			}

			return err
		}
	}
}
