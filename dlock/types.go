package dlock

import (
	"context"
	"time"
)

type Client interface {
	// NewLock 创建一个新的锁
	NewLock(ctx context.Context, key string, options *LockOptions) (Lock, error)
}

type Lock interface {
	// Lock 锁定
	Lock(ctx context.Context) error

	// Unlock 解锁
	Unlock(ctx context.Context) error
	
	// Refresh 续约，延长过期时间
	Refresh(ctx context.Context) error
}

type LockOptions struct {
	// 锁的有效日期
	Expiration time.Duration

	// 锁的重试时间间隔
	RetryDelay time.Duration

	// 锁的重试次数
	RetryCount int
}
