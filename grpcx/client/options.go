package client

import (
	"crypto/tls"
	"time"

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
	timeout    time.Duration
	unaryInts  []grpc.UnaryClientInterceptor
	streamInts []grpc.StreamClientInterceptor

	dialOpts []grpc.DialOption
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		timeout: 2000 * time.Millisecond,
		secure:  false,
	}
}

func Apply(opts ...Option) *Options {
	options := DefaultOptions()
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

// WithTimeout with client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

// WithDialOptions with gRPC client connection options.
func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(c *Options) {
		c.dialOpts = opts
	}
}
