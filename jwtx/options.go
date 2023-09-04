package jwtx

import (
	"github.com/pkg/errors"
	"time"
)

// Option is config option.
type Option func(*Options) error

type Options struct {
	secretKey string

	userID    uint64
	userAgent string

	expireAt time.Time
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		expireAt: time.Now(),
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

// WithSecretKey .
func WithSecretKey(secretKey string) Option {
	return func(o *Options) error {
		if secretKey == "" {
			return errors.New("secretKey can not be empty")
		}
		o.secretKey = secretKey
		return nil
	}
}

// WithUserId .
func WithUserId(userId uint64) Option {
	return func(o *Options) error {
		if userId == 0 {
			return errors.New("UserID can not be empty")
		}
		o.userID = userId
		return nil
	}
}

// WithExpireAt .
func WithExpireAt(expireAt time.Time) Option {
	return func(o *Options) error {
		if expireAt.IsZero() {
			return errors.New("expireAt can not be empty")
		}
		o.expireAt = expireAt
		return nil
	}
}

func WithUserAgent(ua string) Option {
	return func(o *Options) error {
		if ua == "" {
			return errors.New("UserAgent can not be empty")
		}
		o.userAgent = ua
		return nil
	}
}
