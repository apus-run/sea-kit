package redis

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ v1.Cache = (*Cache)(nil)

// Cache 代表Redis缓存
type Cache struct {
	client redis.Cmdable
	lock   sync.RWMutex
}

// NewCache 初始化redis服务
func NewCache(client redis.Cmdable) *Cache {
	// 返回RedisCache实例
	obj := &Cache{
		client: client,
		lock:   sync.RWMutex{},
	}
	return obj
}

// Get 获取某个key对应的值
func (r *Cache) Get(ctx context.Context, key string) (string, error) {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redis.Nil) {
		return "", v1.ErrKeyNotFound
	}

	return cmd.Result()
}

// GetObj 获取某个key对应的对象, 对象必须实现 https://pkg.go.dev/encoding#BinaryUnMarshaler
func (r *Cache) GetObj(ctx context.Context, key string, model interface{}) error {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redis.Nil) {
		return v1.ErrKeyNotFound
	}

	err := cmd.Scan(model)
	if err != nil {
		return err
	}
	return nil
}

// GetMany 获取某些key对应的值
func (r *Cache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := r.client.Pipeline()
	vals := make(map[string]string)
	cmds := make([]*redis.StringCmd, 0, len(keys))

	for _, key := range keys {
		cmds = append(cmds, pipeline.Get(ctx, key))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	errs := make([]string, 0, len(keys))
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		key := cmd.Args()[1].(string)
		vals[key] = val
	}
	return vals, nil
}

// Set 设置某个key和值到缓存，带超时时间
func (r *Cache) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

func (r *Cache) SetNX(ctx context.Context, key string, val string, timeout time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, val, timeout).Result()
}

// SetObj 设置某个key和对象到缓存, 对象必须实现 https://pkg.go.dev/encoding#BinaryMarshaler
func (r *Cache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

// SetMany 设置多个key和值到缓存
func (r *Cache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipline := r.client.Pipeline()
	cmds := make([]*redis.StatusCmd, 0, len(data))
	for k, v := range data {
		cmds = append(cmds, pipline.Set(ctx, k, v, timeout))
	}
	_, err := pipline.Exec(ctx)
	return err
}

// SetForever 设置某个key和值到缓存，不带超时时间
func (r *Cache) SetForever(ctx context.Context, key string, val string) error {
	return r.client.Set(ctx, key, val, v1.NoneDuration).Err()
}

// SetForeverObj 设置某个key和对象到缓存，不带超时时间，对象必须实现 https://pkg.go.dev/encoding#BinaryMarshaler
func (r *Cache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return r.client.Set(ctx, key, val, v1.NoneDuration).Err()
}

// SetTTL 设置某个key的超时时间
func (r *Cache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return r.client.Expire(ctx, key, timeout).Err()
}

// GetTTL 获取某个key的超时时间
func (r *Cache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *Cache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return r.client.IncrBy(ctx, key, step).Result()
}

func (r *Cache) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, 1).Result()
}

func (r *Cache) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, -1).Result()
}

func (r *Cache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Cache) DelMany(ctx context.Context, keys []string) error {
	pipline := r.client.Pipeline()
	cmds := make([]*redis.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipline.Del(ctx, key))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (r *Cache) Contains(key string) bool {
	i, _ := r.client.Exists(context.Background(), key).Result()
	return i > 0
}
