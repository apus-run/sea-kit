package redisx

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	*redis.Options
}

// UniqKey 用来唯一标识一个RedisConfig配置
func (config *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}

type RedisOption func(*RedisConfig) error

// DefaultOptions .
func DefaultOptions() *RedisConfig {
	return &RedisConfig{
		&redis.Options{
			Password: "",
			Addr:     "",
			DB:       0,
		},
	}
}

func Apply(opts ...RedisOption) *RedisConfig {
	options := DefaultOptions()
	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil
		}
	}
	return options
}

// WithRedisConfig 表示自行配置Gorm的配置信息
func WithRedisConfig(f func(options *RedisConfig)) RedisOption {
	return func(config *RedisConfig) error {
		f(config)
		return nil
	}
}
