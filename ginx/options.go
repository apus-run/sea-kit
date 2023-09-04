package ginx

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

// Option is config option.
type Option func(*Options) error

type Options struct {
	mode         string // dev or prod
	host         string
	port         string
	maxPingCount int

	// Before and After funcs
	beforeStart []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		mode:         "dev",
		host:         "localhost",
		port:         "8080",
		maxPingCount: 5,
	}
}

func Apply(opts ...Option) *Options {
	options := DefaultOptions()
	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil
		}
	}
	return options
}

// WithMode .
func WithMode(mode string) Option {
	return func(o *Options) error {
		if strings.ToLower(mode) != "dev" && strings.ToLower(mode) != "prod" {
			return errors.New("mode must be dev or prod")
		}
		o.mode = mode
		return nil
	}
}

// WithAddr .
func WithAddr(host string) Option {
	return func(o *Options) error {
		if host == "" {
			return errors.New("host can not be empty")
		}
		o.host = host
		return nil
	}
}

// WithPort .
func WithPort(port string) Option {
	return func(o *Options) error {
		if port == "" {
			return errors.New("port can not be empty")
		}
		o.port = port
		return nil
	}
}

// WithMaxPingCount .
func WithMaxPingCount(maxPingCount int) Option {
	return func(o *Options) error {
		if maxPingCount <= 0 {
			return errors.New("maxPingCount must be greater than 0")
		}
		o.maxPingCount = maxPingCount
		return nil
	}
}

// Before and Afters

// BeforeStart run funcs before app starts
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *Options) error {
		if fn == nil {
			return errors.New("beforeStart func can not be nil")
		}
		o.beforeStart = append(o.beforeStart, fn)
		return nil
	}
}

// AfterStart run funcs after app starts
func AfterStart(fn func(context.Context) error) Option {
	return func(o *Options) error {
		if fn == nil {
			return errors.New("afterStart func can not be nil")
		}
		o.afterStart = append(o.afterStart, fn)
		return nil
	}
}

// AfterStop run funcs after app stops
func AfterStop(fn func(context.Context) error) Option {
	return func(o *Options) error {
		if fn == nil {
			return errors.New("afterStop func can not be nil")
		}
		o.afterStop = append(o.afterStop, fn)
		return nil
	}
}
