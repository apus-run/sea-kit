package etcd_lock

import (
	"time"

	"github.com/apus-run/sea-kit/dlock"
	"github.com/apus-run/sea-kit/retry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

type EtcdClient struct {
	client *clientv3.Client
}

// NewEtcdClient 创建 etcd 客户端
func NewEtcdClient(c *clientv3.Client) *EtcdClient {
	return &EtcdClient{client: c}
}

// NewLock 初始化锁实例 expiration 代表锁的过期时间，从调用 NewLock 开始计算
func (c *EtcdClient) NewLock(ctx context.Context, key string, expiration time.Duration) (dlock.Lock, error) {
	strategy, _ := retry.NewExponentialBackoffRetryStrategy(time.Millisecond*100, time.Second, 10)
	lt := time.Millisecond * 200

	lock := etcdLock{
		client:     c.client,
		key:        key,
		expiration: expiration,

		lockRetry:   strategy,
		lockTimeout: lt,
	}

	return &lock, nil
}
