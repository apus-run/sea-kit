package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func NewClient(ctx context.Context, opts ...Option) (*grpc.ClientConn, error) {
	var uints []grpc.UnaryClientInterceptor
	var sints []grpc.StreamClientInterceptor

	options := Apply(opts...)

	if len(options.unaryInts) > 0 {
		uints = append(uints, options.unaryInts...)
	}
	if len(options.streamInts) > 0 {
		sints = append(sints, options.streamInts...)
	}
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		}),

		grpc.WithChainUnaryInterceptor(uints...),
		grpc.WithChainStreamInterceptor(sints...),
	}

	if !options.secure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if options.tlsConf != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	}

	if len(options.dialOpts) > 0 {
		dialOpts = append(dialOpts, options.dialOpts...)
	}

	return grpc.DialContext(ctx, options.addr, dialOpts...)
}
