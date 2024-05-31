package error

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/apus-run/sea-kit/zlog"
)

type InterceptorBuilder struct{
	log zlog.Logger
}

func NewErrorInterceptorBuilder(log zlog.Logger) *InterceptorBuilder {
	return &InterceptorBuilder{
		log: log,
	}
}


// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func (e *InterceptorBuilder) ErrorUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		res, err := handler(ctx, req)

		// if this is not a grpc error already, convert it to an internal grpc error
		if err != nil && status.Code(err) == codes.Unknown {
			e.log.Error(err)

			

			err = status.Errorf(codes.Internal, "An internal error occurred.")
		}

		return res, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for panic recovery.
func (e *InterceptorBuilder) ErrorStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		err = handler(srv, stream)

		// if this is not a grpc error already, convert it to an internal grpc error
		if err != nil && status.Code(err) == codes.Unknown {
			e.log.Error(err)
			

			err = status.Errorf(codes.Internal, "An internal error occurred.")
		}

		return err
	}
}
