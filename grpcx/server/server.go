package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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
		grpc.ChainUnaryInterceptor(unaryInterceptor...),
		grpc.ChainStreamInterceptor(streamInterceptor...),
	}
	if options.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(options.tlsConf)))
	}

	// other server option or middleware
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}

	// new grpc server
	options.Server = grpc.NewServer(grpcOpts...)
	if !options.isHealth {
		options.healthServer = health.NewServer()
		grpc_health_v1.RegisterHealthServer(options.Server, options.healthServer)
	}
	// register reflection and the interface can be debugged through the grpcurl tool
	// https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#enable-server-reflection
	// see https://github.com/fullstorydev/grpcurl
	reflection.Register(options.Server)

	return options
}

func (s *Options) Start() error {
	lis, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}
	s.healthServer.Resume()
	return s.Server.Serve(lis)
}

// Stop stop the gRPC server.
func (s *Options) Stop() {
	s.Server.Stop()
	s.healthServer.Shutdown()
}

// GracefulStop graceful stop the gRPC server.
func (s *Options) GracefulStop() {
	s.Server.GracefulStop()
	s.healthServer.Shutdown()
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	grpc://127.0.0.1:9000?isSecure=false
func (s *Options) Endpoint() (string, error) {
	addr, err := Extract(s.addr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc://%s", addr), nil
}
