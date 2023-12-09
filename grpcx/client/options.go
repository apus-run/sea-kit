package client

import (
	"context"
	"crypto/tls"

	"google.golang.org/grpc"
)

// Option is config option.
type Option func(*Options)

type Options struct {
	addr string

	// 是否启用 https
	secure  bool
	tlsConf *tls.Config

	// 拦截器
	unaryInts  []grpc.UnaryClientInterceptor
	streamInts []grpc.StreamClientInterceptor

	dialOpts []grpc.DialOption

	// other options for implementations of the interface
	// can be stored in a context
	ctx context.Context
}

// defaultOptions .
func defaultOptions() *Options {
	return &Options{
		ctx:    context.Background(),
		secure: false,
	}
}

func Apply(opts ...Option) *Options {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithSecure .
func WithSecure(secure bool) Option {
	return func(o *Options) {
		o.secure = secure
	}
}

// WithAddr .
func WithAddr(addr string) Option {
	return func(o *Options) {
		o.addr = addr
	}
}

// WithTLSConfig with TLS config.
func WithTLSConfig(conf *tls.Config) Option {
	return func(o *Options) {
		o.tlsConf = conf
	}
}

// WithUnaryInterceptor returns a DialOption that specifies the interceptor for unary RPCs.
func WithUnaryInterceptor(in ...grpc.UnaryClientInterceptor) Option {
	return func(o *Options) {
		o.unaryInts = in
	}
}

// WithStreamInterceptor returns a DialOption that specifies the interceptor for streaming RPCs.
func WithStreamInterceptor(in ...grpc.StreamClientInterceptor) Option {
	return func(o *Options) {
		o.streamInts = in
	}
}

// WithDialOptions with gRPC client connection options.
func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(c *Options) {
		c.dialOpts = opts
	}
}
