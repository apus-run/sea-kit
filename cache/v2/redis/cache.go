package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	cache "github.com/apus-run/sea-kit/cache/v2"
)

var _ cache.Cache = (*Cache)(nil)

type Cache struct {
	client redis.Cmdable
}

func NewCache(client redis.Cmdable) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	return c.client.Set(ctx, key, val, expiration).Err()
}

func (c *Cache) SetNX(ctx context.Context, key string, val any, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, val, expiration).Result()
}

func (c *Cache) Delete(ctx context.Context, key ...string) (int64, error) {
	return c.client.Del(ctx, key...).Result()
}

func (c *Cache) Get(ctx context.Context, key string) (val cache.Value) {
	val.Value, val.Error = c.client.Get(ctx, key).Result()
	if val.Error != nil && errors.Is(val.Error, redis.Nil) {
		val.Error = cache.ErrKeyNotExist
	}
	return
}

func (c *Cache) GetSet(ctx context.Context, key string, val string) (result cache.Value) {
	result.Value, result.Error = c.client.GetSet(ctx, key, val).Result()
	if result.Error != nil && errors.Is(result.Error, redis.Nil) {
		result.Error = cache.ErrKeyNotExist
	}
	return
}

func (c *Cache) LPush(ctx context.Context, key string, val ...any) (int64, error) {
	return c.client.LPush(ctx, key, val...).Result()
}

func (c *Cache) LPop(ctx context.Context, key string) (result cache.Value) {
	result.Value, result.Error = c.client.LPop(ctx, key).Result()
	if result.Error != nil && errors.Is(result.Error, redis.Nil) {
		result.Error = cache.ErrKeyNotExist
	}
	return
}

func (c *Cache) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return c.client.SAdd(ctx, key, members...).Result()
}

func (c *Cache) SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return c.client.SRem(ctx, key, members...).Result()
}

func (c *Cache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, key, value).Result()
}

func (c *Cache) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.DecrBy(ctx, key, value).Result()
}

func (c *Cache) IncrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	return c.client.IncrByFloat(ctx, key, value).Result()
}
