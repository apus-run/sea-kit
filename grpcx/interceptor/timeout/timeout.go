package timeout

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type InterceptorBuilder struct {
	timeout time.Duration

	// 忽略指定路由的超时控制
	ignoreMethods map[string]struct{}
}

func NewTimeoutInterceptorBuilder(t time.Duration) *InterceptorBuilder {
	return &InterceptorBuilder{
		timeout: t,
		ignoreMethods: map[string]struct{}{
			"/grpc.health.v1.Health/Check": {},
		},
	}
}

func (b *InterceptorBuilder) IgnoreMethods(fullMethodNames ...string) *InterceptorBuilder {
	for _, method := range fullMethodNames {
		b.ignoreMethods[method] = struct{}{}
	}

	return b
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		// 检查方法是否在忽略列表中
		if _, ok := b.ignoreMethods[info.FullMethod]; ok {
			// 如果在忽略列表中，直接调用原始处理程序
			return handler(ctx, req)
		}

		ctx, cancel := context.WithTimeout(ctx, b.timeout)
		defer cancel()

		var lock sync.Mutex
		done := make(chan struct{})
		// create channel with buffer size 1 to avoid goroutine leak
		panicChan := make(chan any, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// attach call stack to avoid missing in different goroutine
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			lock.Lock()
			defer lock.Unlock()
			resp, err = handler(ctx, req)
			close(done)
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-done:
			lock.Lock()
			defer lock.Unlock()
			return resp, err
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				err = status.Error(codes.Canceled, err.Error())
			} else if errors.Is(err, context.DeadlineExceeded) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, err
		}
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
		// 检查方法是否在忽略列表中
		if _, ok := b.ignoreMethods[method]; ok {
			// 如果在忽略列表中，直接调用原始处理程序
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		t := getTimeoutFromCallOptions(opts, b.timeout)
		if t <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// TimeoutCallOption is a call option that controls timeout.
type TimeoutCallOption struct {
	grpc.EmptyCallOption
	timeout time.Duration
}

// WithCallTimeout returns a call option that controls method call timeout.
func WithCallTimeout(timeout time.Duration) grpc.CallOption {
	return TimeoutCallOption{
		timeout: timeout,
	}
}

func getTimeoutFromCallOptions(opts []grpc.CallOption, defaultTimeout time.Duration) time.Duration {
	for _, opt := range opts {
		if o, ok := opt.(TimeoutCallOption); ok {
			return o.timeout
		}
	}

	return defaultTimeout
}
