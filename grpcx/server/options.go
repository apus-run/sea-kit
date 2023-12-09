package server

import (
	"context"
	"crypto/tls"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

// Option is config option.
type Option func(*Options)

type Options struct {
	*grpc.Server
	err  error
	addr string
	// server listen network tcp/udp
	network string
	lis     net.Listener

	// 是否启用 https
	tlsConf *tls.Config

	// 拦截器
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor

	grpcOpts []grpc.ServerOption

	// other options for implementations of the interface
	// can be stored in a context
	ctx context.Context

	// 心跳检测
	healthServer *health.Server
	isHealth     bool
}

// defaultOptions .
func defaultOptions() *Options {
	return &Options{
		ctx:     context.Background(),
		network: "tcp",
		addr:    ":0",
	}
}

func Apply(opts ...Option) *Options {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithNetwork with server network.
func WithNetwork(network string) Option {
	return func(s *Options) {
		s.network = network
	}
}

// WithAddr .
func WithAddr(addr string) Option {
	return func(o *Options) {
		o.addr = addr
	}
}

// WithListener with server lis
func WithListener(lis net.Listener) Option {
	return func(s *Options) {
		s.lis = lis
	}
}

// WithTLSConfig with TLS config.
func WithTLSConfig(conf *tls.Config) Option {
	return func(o *Options) {
		o.tlsConf = conf
	}
}

// WithIsHealth Checks server.
func WithIsHealth() Option {
	return func(s *Options) {
		s.isHealth = true
	}
}

// WithUnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) Option {
	return func(s *Options) {
		s.unaryInts = in
	}
}

// WithStreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the server.
func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) Option {
	return func(s *Options) {
		s.streamInts = in
	}
}

// WithGrpcOptions with gRPC options.
func WithGrpcOptions(opts ...grpc.ServerOption) Option {
	return func(c *Options) {
		c.grpcOpts = opts
	}
}
