package redis

import (
	"context"
	"strings"

	ratelimit "github.com/apus-run/sea-kit/ratelimit/redis"
	"github.com/apus-run/sea-kit/zlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Kind is the type of Interceptor
const Kind string = "Redis"

type InterceptorBuilder struct {
	limiter ratelimit.Limiter
	key     string

	log zlog.Logger
}

// NewRatelimitInterceptorBuilder key: user-service
// "limiter:service:user" 整个应用、集群限流
// "limiter:service:user:UserService" user 里面的 UserService 限流
func NewRatelimitInterceptorBuilder(l ratelimit.Limiter, key string, log zlog.Logger) *InterceptorBuilder {
	return &InterceptorBuilder{
		limiter: l,
		key:     key,
		log:     log,
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
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			b.log.Error("触发限流", zlog.Error(err))
			// 这里采用保守措施，在触发限流之后直接返回
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}
		if limited {
			b.log.Error("触发限流", zlog.Error(err))
			ctx = context.WithValue(ctx, "limited", "true")
		}
		return handler(ctx, req)
	}
}

// BuildUnaryServerInterceptorService 服务级别限流
func (b *InterceptorBuilder) BuildUnaryServerInterceptorService(prefix string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if strings.HasPrefix(info.FullMethod, prefix) {
			limited, err := b.limiter.Limit(ctx, b.key)
			if err != nil {
				b.log.Error("触发限流", zlog.Error(err))
				// 这里采用保守措施，在触发限流之后直接返回
				return nil, status.Errorf(codes.ResourceExhausted, "限流")
			}
			if limited {
				b.log.Error("触发限流", zlog.Error(err))
				ctx = context.WithValue(ctx, "limited", "true")
			}
		}
		return handler(ctx, req)
	}
}

func (b *InterceptorBuilder) BuildUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			b.log.Error("触发限流", zlog.Error(err))
			// 这里采用保守措施，在触发限流之后直接返回
			return status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		if limited {
			b.log.Error("触发限流", zlog.Error(err))
			ctx = context.WithValue(ctx, "limited", "true")
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
