package redislock

import "time"

type lockOptions struct {
	key        string
	expiration time.Duration
	retry      RetryStrategy
	timeout    time.Duration
}

type LockOption func(*lockOptions)

// DefaultOptions .
func DefaultOptions() *lockOptions {
	return &lockOptions{
		expiration: time.Second * 30,
		// 每隔 100ms 重试一次，每次重试的超时时间是 1s
		retry: &FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
		},
		timeout: time.Second,
	}
}

func Apply(opts ...LockOption) *lockOptions {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithKey(key string) LockOption {
	return func(o *lockOptions) {
		o.key = key
	}
}

func WithRetry(retry RetryStrategy) LockOption {
	return func(o *lockOptions) {
		o.retry = retry
	}
}

func WithExpiration(expiration time.Duration) LockOption {
	return func(o *lockOptions) {
		o.expiration = expiration
	}
}

func WithTimeout(timeout time.Duration) LockOption {
	return func(o *lockOptions) {
		o.timeout = timeout
	}
}
