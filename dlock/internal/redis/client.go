package rlock

import (
	"context"
	"time"

	"github.com/apus-run/sea-kit/dlock"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb redis.Cmdable
}

func NewClient(rdb redis.Cmdable) *Client {
	return &Client{rdb: rdb}
}

func (c *Client) NewLock(ctx context.Context, key string, expiration time.Duration) (dlock.Lock, error) {
	return NewLock(c.rdb, key, expiration), nil
}
