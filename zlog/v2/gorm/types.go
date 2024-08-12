package gorm

import (
	"time"

	"github.com/apus-run/sea-kit/zlog/v2"
)

type Option func(*Config)

func DefaultDBConfig() *Config {
	return &Config{
		SlowThreshold:             time.Second * 2,
		IgnoreRecordNotFoundError: true,
		LogInfo:                   true,
		Options:                   zlog.DefaultOptions(),
	}
}

func Apply(opts ...Option) *Config {
	config := DefaultDBConfig()
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func WithSlowThreshold(threshold time.Duration) Option {
	return func(config *Config) {
		config.SlowThreshold = threshold
	}
}
func WithIgnoreRecordNotFoundError() Option {
	return func(config *Config) {
		config.IgnoreRecordNotFoundError = true
	}
}

func WithLogInfo(logInfo bool) Option {
	return func(config *Config) {
		config.LogInfo = logInfo
	}
}

// WithConfig 设置所有配置
func WithConfig(fn func(config *Config)) Option {
	return func(config *Config) {
		fn(config)
	}
}

// WithZlogOptions 设置 zlog.Options
func WithZlogOptions(options *zlog.Options) Option {
	return func(conf *Config) {
		conf.Options = options
	}
}

// Config is logger config
type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogInfo                   bool

	// 集成 zlog 配置
	*zlog.Options
}
