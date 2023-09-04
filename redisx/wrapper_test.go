package redisx

import (
	"context"
	"testing"
	"time"
)

type Client struct {
	*Helper
}

func TestRedis_GetClient(t *testing.T) {
	ctx := context.Background()
	h := &Client{NewHelper()}
	client, err := h.GetClient(WithRedisConfig(func(options *RedisConfig) {
		options.Addr = "localhost:16379"
		options.DB = 0
		options.Username = "root"
	}))
	if err != nil {
		t.Fatal(err)
	}

	// 检测数据库是否可以连接
	cmd := client.Ping(ctx)
	if cmd.Err() != nil {
		t.Fatal(cmd.Err())
	}

	err = client.Set(ctx, "foo", "bar", 1*time.Hour).Err()

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("foo", val)

	err = client.Del(ctx, "foo").Err()
	if err != nil {
		t.Fatal(err)
	}

}
