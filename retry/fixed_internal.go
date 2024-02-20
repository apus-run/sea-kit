package retry

import (
	"fmt"
	"sync/atomic"
	"time"
)

// FixedIntervalRetryStrategy 等间隔重试
type FixedIntervalRetryStrategy struct {
	maxRetries int32         // 最大重试次数，如果是 0 或负数，表示无限重试
	interval   time.Duration // 重试间隔时间
	retries    int32         // 当前重试次数
}

func NewFixedIntervalRetryStrategy(interval time.Duration, maxRetries int32) (*FixedIntervalRetryStrategy, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("无效的间隔时间 %d, 预期值应大于 0", interval)
	}
	return &FixedIntervalRetryStrategy{
		maxRetries: maxRetries,
		interval:   interval,
	}, nil
}

func (s *FixedIntervalRetryStrategy) Next() (time.Duration, bool) {
	retries := atomic.AddInt32(&s.retries, 1)
	if s.maxRetries <= 0 || retries <= s.maxRetries {
		return s.interval, true
	}
	return 0, false
}
