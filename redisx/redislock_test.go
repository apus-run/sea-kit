package redisx

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/apus-run/sea-kit/stringx"
)

func runOnRedis(t *testing.T, fn func(client *redis.Client)) {
	s := miniredis.RunT(t)

	h := NewHelper()
	client, err := h.GetClient(WithRedisConfig(func(options *RedisConfig) {
		options.Addr = s.Addr()
		options.DB = 0
		options.Username = "root"
	}))

	if err != nil {
		t.Fatal(err)
	}

	fn(client)
}

func TestRedisLock(t *testing.T) {
	testFn := func(ctx context.Context) func(client *redis.Client) {
		return func(client *redis.Client) {
			key := stringx.Rand()
			firstLock := NewRedisLock(client, key)
			firstLock.SetExpire(5)
			firstAcquire, err := firstLock.Acquire()
			assert.Nil(t, err)
			assert.True(t, firstAcquire)

			secondLock := NewRedisLock(client, key)
			secondLock.SetExpire(5)
			againAcquire, err := secondLock.Acquire()
			assert.Nil(t, err)
			assert.False(t, againAcquire)

			release, err := firstLock.Release()
			assert.Nil(t, err)
			assert.True(t, release)

			endAcquire, err := secondLock.Acquire()
			assert.Nil(t, err)
			assert.True(t, endAcquire)
		}
	}

	t.Run("normal", func(t *testing.T) {
		runOnRedis(t, testFn(nil))
	})

	t.Run("withContext", func(t *testing.T) {
		runOnRedis(t, testFn(context.Background()))
	})
}

func TestRedisLock_Expired(t *testing.T) {
	runOnRedis(t, func(client *redis.Client) {
		key := stringx.Rand()
		redisLock := NewRedisLock(client, key)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := redisLock.AcquireCtx(ctx)
		assert.NotNil(t, err)
	})

	runOnRedis(t, func(client *redis.Client) {
		key := stringx.Rand()
		redisLock := NewRedisLock(client, key)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := redisLock.ReleaseCtx(ctx)
		assert.NotNil(t, err)
	})
}
