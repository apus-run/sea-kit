package retry

import (
	"context"
	"fmt"
	"time"
)

// Retry 会在以下条件满足的情况下返回：
// 1. 重试达到了最大次数，而后返回重试耗尽的错误
// 2. ctx 被取消或者超时
// 3. bizFunc 没有返回 error
// 而只要 bizFunc 返回 error，就会尝试重试
func Retry(ctx context.Context,
	s Strategy,
	bizFunc func() error) error {
	var ticker *time.Ticker
	defer func() {
		if ticker != nil {
			ticker.Stop()
		}
	}()
	for {
		err := bizFunc()
		// 直接退出
		if err == nil {
			return nil
		}
		duration, ok := s.Next()
		if !ok {
			return fmt.Errorf("ekit: 重试耗尽, 原因: %w", err)
		}
		if ticker == nil {
			ticker = time.NewTicker(duration)
		} else {
			ticker.Reset(duration)
		}
		select {
		case <-ctx.Done():
			// 超时或者被取消了，直接返回
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
