package redisx

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type (
	// IntCmd is an alias of redis.IntCmd.
	IntCmd = redis.IntCmd
	// FloatCmd is an alias of redis.FloatCmd.
	FloatCmd = redis.FloatCmd
	// StringCmd is an alias of redis.StringCmd.
	StringCmd = redis.StringCmd
	// Script is an alias of redis.Script.
	Script = redis.Script
)

type Wrapper interface {
	// GetClient 获取redis连接实例
	GetClient(option ...RedisOption) (*redis.Client, error)
}

type Helper struct {
	lock *sync.RWMutex

	clients map[string]*redis.Client // key为uniqKey, value为redis.Client (连接池）
}

// NewHelper new a logger helper.
func NewHelper() *Helper {
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}

	return &Helper{
		clients: clients,
		lock:    lock,
	}
}

// GetClient 获取Client实例
func (h *Helper) GetClient(option ...RedisOption) (*redis.Client, error) {
	config := Apply(option...)
	if config == nil {
		return nil, fmt.Errorf("invalid redis config")
	}

	// 如果最终的config没有设置dsn,就生成dsn
	key := config.UniqKey()

	h.lock.RLock()
	c, ok := h.clients[key]
	if ok {
		h.lock.RUnlock()
		return c, nil
	}
	h.lock.RUnlock()

	// 如果没有实例化,就实例化
	h.lock.Lock()
	defer h.lock.Unlock()

	client := redis.NewClient(config.Options)

	// 挂载到map中
	h.clients[key] = client

	return client, nil
}

// NewScript returns a new Script instance.
func NewScript(script string) *Script {
	return redis.NewScript(script)
}
