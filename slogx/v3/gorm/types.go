package gorm

import (
	"time"

	log "github.com/apus-run/van/log/slog"
)

// Config is logger config
type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogInfo                   bool

	*log.Options
}

type Option func(*Config)

func DefaultDBConfig() *Config {
	return &Config{
		SlowThreshold:             time.Second * 2,
		IgnoreRecordNotFoundError: true,
		LogInfo:                   true,
		Options:                   log.DefaultOptions(),
	}
}

func Apply(opts ...Option) *Config {
	config := DefaultDBConfig()
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// WithSlowThreshold set slow threshold
func WithSlowThreshold(threshold time.Duration) Option {
	return func(c *Config) {
		c.SlowThreshold = threshold
	}
}

// WithIgnoreRecordNotFoundError ignore record not found error
func WithIgnoreRecordNotFoundError() Option {
	return func(c *Config) {
		c.IgnoreRecordNotFoundError = true
	}
}

// WithLogInfo set log info
func WithLogInfo(logInfo bool) Option {
	return func(c *Config) {
		c.LogInfo = logInfo
	}
}

// WithConfig set all config
func WithConfig(fn func(config *Config)) Option {
	return func(config *Config) {
		fn(config)
	}
}

// WithSlogOptions 设置 slog.Options
func WithSlogOptions(options *log.Options) Option {
	return func(conf *Config) {
		conf.Options = options
	}
}
