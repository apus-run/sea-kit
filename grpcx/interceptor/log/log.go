package log

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/apus-run/sea-kit/grpcx/interceptor"
	"github.com/apus-run/sea-kit/zlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Kind is the type of Interceptor
const Kind string = "Log"

type InterceptorBuilder struct {
	log zlog.Logger
	interceptor.Builder

	// 忽略指定路由的日志打印
	ignoreMethods map[string]struct{}
}

func NewLoggerInterceptorBuilder(log zlog.Logger) *InterceptorBuilder {
	return &InterceptorBuilder{
		log: log,
		ignoreMethods: map[string]struct{}{
			"/grpc.health.v1.Health/Check": {},
		},
	}
}

// Kind return the name of interceptor
func (b *InterceptorBuilder) Kind() string {
	return Kind
}

func (b *InterceptorBuilder) IgnoreMethods(fullMethodNames ...string) *InterceptorBuilder {
	for _, method := range fullMethodNames {
		b.ignoreMethods[method] = struct{}{}
	}

	return b
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		// ignore printing of the specified method
		if _, ok := b.ignoreMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		var start = time.Now()
		var fields = make([]zlog.Field, 0)
		var event = "normal"

		defer func() {
			cost := time.Since(start)
			if rec := recover(); rec != nil {
				switch recType := rec.(type) {
				case error:
					err = recType
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				event = "recover"
				err = status.New(codes.Internal, "panic, err "+err.Error()).Err()
			}
			st, _ := status.FromError(err)
			if st != nil {
				fields = append(
					fields,
					zlog.String("type", "UnaryServer"),
					zlog.String("peer_name", b.Builder.PeerName(ctx)),
					zlog.String("grpc_code", st.Code().String()),
					zlog.String("grpc_message", st.Message()),
					zlog.String("method", info.FullMethod),
					zlog.String("event", event),
					zlog.String("peer", b.PeerName(ctx)),
					zlog.String("peer_ip", b.Builder.PeerIP(ctx)),
					zlog.Duration("cost", cost),
				)
			}

			b.log.Info("grpc server", fields...)
		}()

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
		// ignore printing of the specified method
		if _, ok := b.ignoreMethods[method]; ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		var start = time.Now()
		var fields = make([]zlog.Field, 0)
		var event = "normal"

		err := invoker(ctx, method, req, reply, cc, opts...)

		defer func() {
			cost := time.Since(start)
			if rec := recover(); rec != nil {
				switch recType := rec.(type) {
				case error:
					err = recType
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				event = "recover"
				err = status.New(codes.Internal, "panic, err "+err.Error()).Err()
			}
			st, _ := status.FromError(err)
			if st != nil {
				fields = append(
					fields,
					zlog.String("type", "UnaryClient"),
					zlog.String("peer_name", b.Builder.PeerName(ctx)),
					zlog.String("grpc_code", st.Code().String()),
					zlog.String("grpc_message", st.Message()),
					zlog.String("method", method),
					zlog.String("event", event),
					zlog.String("peer", b.PeerName(ctx)),
					zlog.String("peer_ip", b.Builder.PeerIP(ctx)),
					zlog.Duration("cost", cost),
				)
			}

			b.log.Info("grpc server", fields...)
		}()

		return err
	}
}
