package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/simplelru"

	"github.com/apus-run/sea-kit/list"
	"github.com/apus-run/sea-kit/set"
	"github.com/apus-run/sea-kit/simple_cache"
	"github.com/apus-run/sea-kit/simple_cache/internal/errs"
)

var (
	_ simple_cache.Cache = (*Cache)(nil)
)

type Cache struct {
	lock   sync.RWMutex
	client simplelru.LRUCache[string, any]
}

func NewCache(client simplelru.LRUCache[string, any]) *Cache {
	return &Cache{
		lock:   sync.RWMutex{},
		client: client,
	}
}

// Set expiration 无效 由lru 统一控制过期时间
func (c *Cache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.client.Add(key, val)
	return nil
}

// SetNX expiration 无效 由lru 统一控制过期时间
func (c *Cache) SetNX(ctx context.Context, key string, val any, expiration time.Duration) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.client.Contains(key) {
		return false, nil
	}

	c.client.Add(key, val)

	return true, nil
}

func (c *Cache) Get(ctx context.Context, key string) (val simple_cache.Value) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var ok bool
	val.Value, ok = c.client.Get(key)
	if !ok {
		val.Error = errs.ErrKeyNotExist
	}

	return
}

func (c *Cache) GetSet(ctx context.Context, key string, val string) (result simple_cache.Value) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var ok bool
	result.Value, ok = c.client.Get(key)
	if !ok {
		result.Error = errs.ErrKeyNotExist
	}

	c.client.Add(key, val)

	return
}

func (c *Cache) Delete(ctx context.Context, key ...string) (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	n := int64(0)
	for _, k := range key {
		if ctx.Err() != nil {
			return n, ctx.Err()
		}
		_, ok := c.client.Get(k)
		if !ok {
			continue
		}
		if c.client.Remove(k) {
			n++
		} else {
			return n, fmt.Errorf("%w: key = %s", errs.ErrDeleteKeyFailed, k)
		}
	}
	return n, nil
}

// anySliceToValueSlice 公共转换
func (c *Cache) anySliceToValueSlice(data ...any) []simple_cache.Value {
	newVal := make([]simple_cache.Value, len(data), cap(data))
	for key, value := range data {
		anyVal := simple_cache.Value{}
		anyVal.Value = value
		newVal[key] = anyVal
	}
	return newVal
}

func (c *Cache) LPush(ctx context.Context, key string, val ...any) (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var (
		ok     bool
		result = simple_cache.Value{}
	)
	result.Value, ok = c.client.Get(key)
	if !ok {
		l := &list.ConcurrentList[simple_cache.Value]{
			List: list.NewLinkedListOf[simple_cache.Value](c.anySliceToValueSlice(val...)),
		}
		c.client.Add(key, l)
		return int64(l.Len()), nil
	}

	data, ok := result.Value.(list.List[simple_cache.Value])
	if !ok {
		return 0, errors.New("当前key不是list类型")
	}

	err := data.Append(c.anySliceToValueSlice(val)...)
	if err != nil {
		return 0, err
	}

	c.client.Add(key, data)
	return int64(data.Len()), nil
}

func (c *Cache) LPop(ctx context.Context, key string) (val simple_cache.Value) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var (
		ok bool
	)
	val.Value, ok = c.client.Get(key)
	if !ok {
		val.Error = errs.ErrKeyNotExist
		return
	}

	data, ok := val.Value.(list.List[simple_cache.Value])
	if !ok {
		val.Error = errors.New("当前key不是list类型")
		return
	}

	value, err := data.Delete(0)
	if err != nil {
		val.Error = err
		return
	}

	val = value
	return
}

func (c *Cache) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var (
		ok     bool
		result = simple_cache.Value{}
	)
	result.Value, ok = c.client.Get(key)
	if !ok {
		result.Value = set.NewMapSet[any](8)
	}

	s, ok := result.Value.(set.Set[any])
	if !ok {
		return 0, errors.New("当前key已存在不是set类型")
	}

	for _, value := range members {
		s.Add(value)
	}
	c.client.Add(key, s)

	return int64(len(s.Keys())), nil
}

func (c *Cache) SRem(ctx context.Context, key string, members ...any) (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	result, ok := c.client.Get(key)
	if !ok {
		return 0, errs.ErrKeyNotExist
	}

	s, ok := result.(set.Set[any])
	if !ok {
		return 0, errors.New("当前key已存在不是set类型")
	}

	var rems int64
	for _, member := range members {
		if s.Exist(member) {
			s.Delete(member)
			rems++
		}
	}
	return rems, nil
}

func (c *Cache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var (
		ok     bool
		result = simple_cache.Value{}
	)
	result.Value, ok = c.client.Get(key)
	if !ok {
		c.client.Add(key, value)
		return value, nil
	}

	incr, err := result.Int64()
	if err != nil {
		return 0, errors.New("当前key不是int64类型")
	}

	newVal := incr + value
	c.client.Add(key, newVal)

	return newVal, nil
}

func (c *Cache) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var (
		ok     bool
		result = simple_cache.Value{}
	)
	result.Value, ok = c.client.Get(key)
	if !ok {
		c.client.Add(key, -value)
		return -value, nil
	}

	decr, err := result.Int64()
	if err != nil {
		return 0, errors.New("当前key不是int64类型")
	}

	newVal := decr - value
	c.client.Add(key, newVal)

	return newVal, nil
}

func (c *Cache) IncrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var (
		ok     bool
		result = simple_cache.Value{}
	)
	result.Value, ok = c.client.Get(key)
	if !ok {
		c.client.Add(key, value)
		return value, nil
	}

	val, err := result.Float64()
	if err != nil {
		return 0, errors.New("当前key不是float64类型")
	}

	newVal := val + value
	c.client.Add(key, newVal)

	return newVal, nil
}
