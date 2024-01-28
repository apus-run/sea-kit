package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Uri            string
	DatabaseName   string
	CollectionName string

	*options.ClientOptions
}

type Option func(*Config)

// DefaultOptions .
func DefaultOptions() *Config {
	connectTimeout := 30 * time.Second
	maxConnIdleTime := 3 * time.Minute
	minPoolSize := uint64(20)
	maxPoolSize := uint64(300)

	return &Config{
		ClientOptions: &options.ClientOptions{
			ConnectTimeout:  &connectTimeout,
			MaxConnIdleTime: &maxConnIdleTime,
			MinPoolSize:     &minPoolSize,
			MaxPoolSize:     &maxPoolSize,
		},
	}
}

func Apply(opts ...Option) *Config {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithURI(uri string) Option {
	return func(config *Config) {
		config.Uri = uri
	}
}

func WithDatabaseName(dbname string) Option {
	return func(config *Config) {
		config.DatabaseName = dbname
	}
}

func WithCollectionName(collname string) Option {
	return func(config *Config) {
		config.CollectionName = collname
	}
}

// WithMongoConfig 设置所有配置
func WithMongoConfig(fn func(options *Config)) Option {
	return func(config *Config) {
		fn(config)
	}
}

// WithClientOptions 表示自行配置ClientOptions的配置信息
func WithClientOptions(options *options.ClientOptions) Option {
	return func(config *Config) {
		config.ClientOptions = options
	}
}
