package ratelimit_redis

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:embed lua/slide_window.lua
var luaScript string

// RedisSlidingWindowLimiter Redis 上的滑动窗口算法限流器实现
type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable

	// 窗口大小
	interval time.Duration

	// 阈值
	// interval 内允许 rate 个请求
	// 1s 内允许 3000 个请求
	rate int
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable,
	interval time.Duration, rate int) Limiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return false, fmt.Errorf("generate uuid failed: %w", err)
	}
	return r.cmd.Eval(ctx, luaScript, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli(), uid.String()).Bool()
}
