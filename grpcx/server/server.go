package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

// NewServer returns new unsecured grpc server
func NewServer(opts ...Option) *Options {
	var unaryInterceptor []grpc.UnaryServerInterceptor
	var streamInterceptor []grpc.StreamServerInterceptor

	options := Apply(opts...)

	if len(options.unaryInts) > 0 {
		unaryInterceptor = append(unaryInterceptor, options.unaryInts...)
	}
	if len(options.streamInts) > 0 {
		streamInterceptor = append(streamInterceptor, options.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    time.Second * 30, // server initiated keep alive interval
			Timeout: time.Second * 30, // server initiated keep alive timeout
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Minute, // server enforcement for client keep alive
			PermitWithoutStream: true,        // allow connection without any active ongoing streams
		}),

		grpc.ChainUnaryInterceptor(unaryInterceptor...),
		grpc.ChainStreamInterceptor(streamInterceptor...),
	}
	if options.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(options.tlsConf)))
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}

	options.Server = grpc.NewServer(grpcOpts...)

	// register reflection and the interface can be debugged through the grpcurl tool
	// https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#enable-server-reflection
	// see https://github.com/fullstorydev/grpcurl
	reflection.Register(options.Server)

	return options
}

func (s *Options) Start() error {
	l, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}
	return s.Server.Serve(l)
}

// Stop stop the gRPC server.
func (s *Options) Stop() {
	s.Server.Stop()
}
