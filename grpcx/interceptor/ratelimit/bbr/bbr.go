package bbr

import (
	"context"

	"google.golang.org/grpc"

	ratelimit "github.com/apus-run/sea-kit/ratelimit/bbr"
	"github.com/apus-run/sea-kit/zlog"
)

// Kind is the type of Interceptor
const Kind string = "BBR"

type InterceptorBuilder struct {
	log     zlog.Logger
	limiter ratelimit.Limiter
}

func NewRatelimitInterceptorBuilder(l ratelimit.Limiter, log zlog.Logger) *InterceptorBuilder {
	return &InterceptorBuilder{
		log:     log,
		limiter: l,
	}
}

// Kind return the name of interceptor
func (b *InterceptorBuilder) Kind() string {
	return Kind
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		done, err := b.limiter.Allow()
		if err != nil {
			return nil, err
		}

		reply, err := handler(ctx, req)
		done(ratelimit.DoneInfo{Err: err})
		return reply, err
	}
}

func (b *InterceptorBuilder) BuildStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		done, err := b.limiter.Allow()
		if err != nil {
			return err
		}

		err = handler(srv, ss)
		done(ratelimit.DoneInfo{Err: err})
		return err
	}
}
