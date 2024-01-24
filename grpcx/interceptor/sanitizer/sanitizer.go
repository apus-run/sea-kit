package sanitizer

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Sanitizer interface {
	Sanitize() error
}

type InterceptorBuilder struct{}

func NewSanitizerInterceptorBuilder() *InterceptorBuilder {
	return &InterceptorBuilder{}
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		if err := sanitize(ctx, req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func sanitize(_ context.Context, reqOrRes interface{}) (err error) {
	switch v := reqOrRes.(type) {
	case Sanitizer:
		err = v.Sanitize()
	}

	if err == nil {
		return nil
	}

	return status.Error(codes.InvalidArgument, err.Error())
}
