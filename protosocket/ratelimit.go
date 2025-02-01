package protosocket

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/protobuf/proto"
)

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

func (r *RateLimiter) Middleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, msg proto.Message) error {
			if err := r.limiter.Wait(ctx); err != nil {
				return err
			}
			return next(ctx, msg)
		}
	}
}
