package protosocket

import (
	"context"
	"errors"

	"google.golang.org/protobuf/proto"
)

type Authenticator interface {
	Authenticate(token string) (Claims, error)
}

type Claims map[string]interface{}

type AuthMiddleware struct {
	authenticator Authenticator
}

const TokenContextKey = "token"

func GetTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(TokenContextKey).(string); ok {
		return token
	}
	return ""
}

func NewAuthMiddleware(auth Authenticator) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, msg proto.Message) error {
			token := GetTokenFromContext(ctx)
			claims, err := auth.Authenticate(token)
			if err != nil {
				return errors.New("unauthorized")
			}

			ctx = context.WithValue(ctx, "claims", claims)
			return next(ctx, msg)
		}
	}
}
