package breaker

import (
	"context"

	"github.com/apus-run/sea-kit/grpcx/interceptor/breaker/circuitbreaker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/apus-run/sea-kit/grpcx/interceptor/breaker/circuitbreaker/sre"
	"github.com/apus-run/sea-kit/grpcx/interceptor/breaker/group"
)

// Kind is the type of Interceptor
const Kind string = "CircuitBreaker"

type InterceptorBuilder struct {
	group *group.Group

	// rpc code for circuit breaker, default already includes codes.Internal and codes.Unavailable
	validCodes map[codes.Code]struct{}
}

func NewBreakerInterceptorBuilder() *InterceptorBuilder {
	return &InterceptorBuilder{
		group: group.NewGroup(func() any {
			return sre.NewBreaker()
		}),
		validCodes: map[codes.Code]struct{}{
			codes.Internal:    {},
			codes.Unavailable: {},
		},
	}
}

// Kind return the name of interceptor
func (b *InterceptorBuilder) Kind() string {
	return Kind
}

func (b *InterceptorBuilder) ValidCode(codes ...codes.Code) *InterceptorBuilder {
	for _, c := range codes {
		b.validCodes[c] = struct{}{}
	}

	return b
}

// Group with circuit breaker group.
// NOTE: implements generics circuitbreaker.CircuitBreaker
func (b *InterceptorBuilder) Group(g *group.Group) *InterceptorBuilder {
	b.group = g
	return b
}

// CircuitBreaker with circuit breaker genFunc.
func (b *InterceptorBuilder) CircuitBreaker(genBreakerFunc func() circuitbreaker.CircuitBreaker) *InterceptorBuilder {
	b.group = group.NewGroup(func() any {
		return genBreakerFunc()
	})
	return b
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		breaker := b.group.Get(info.FullMethod).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally, keep adding let the drop ratio higher.
			breaker.MarkFailed()
			return nil, err
		}

		reply, err := handler(ctx, req)
		if err != nil {
			// 借助这个区判定是不是业务错误
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			_, isHit := b.validCodes[s.Code()]
			if ok && isHit {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}

		// 触发了熔断器
		return reply, err
	}
}

func (b *InterceptorBuilder) BuildStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		breaker := b.group.Get(info.FullMethod).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			breaker.MarkFailed()
			return err
		}
		err := handler(srv, ss)
		if err != nil {
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			_, isHit := b.validCodes[s.Code()]
			if ok && isHit {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}

		return err
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
		breaker := b.group.Get(method).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			breaker.MarkFailed()
			return err
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			_, isHit := b.validCodes[s.Code()]
			if ok && isHit {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}
		return err
	}
}

func (b *InterceptorBuilder) BuildStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		breaker := b.group.Get(method).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally, keep adding counter let the drop ratio higher.
			breaker.MarkFailed()
			return nil, err
		}

		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			_, isHit := b.validCodes[s.Code()]
			if ok && isHit {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}

		return clientStream, err
	}
}
