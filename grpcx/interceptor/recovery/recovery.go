package recovery

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct{}

func NewRecoveryInterceptorBuilder() *InterceptorBuilder {
	return &InterceptorBuilder{}
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	// https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/recovery
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "triggered panic: %v", p)
	}
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(customFunc),
	}

	return recovery.UnaryServerInterceptor(opts...)
}

func (b *InterceptorBuilder) BuildStreamServerInterceptor() grpc.StreamServerInterceptor {
	// https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/recovery
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "triggered panic: %v", p)
	}
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(customFunc),
	}

	return recovery.StreamServerInterceptor(opts...)
}

func (b *InterceptorBuilder) BuildUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "triggered panic: %v", r)
			}
		}()

		err = invoker(ctx, method, req, reply, cc, opts...)
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
		opts ...grpc.CallOption) (s grpc.ClientStream, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "triggered panic: %v", r)
			}
		}()

		s, err = streamer(ctx, desc, cc, method, opts...)
		return s, err
	}
}
