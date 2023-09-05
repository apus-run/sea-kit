package memory

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/apus-run/sea-kit/cache"
)

var _ cache.Cache = (*Cache)(nil)

type MemoryData struct {
	val        interface{}
	createTime time.Time
	ttl        time.Duration
}

type Cache struct {
	data map[string]*MemoryData
	lock sync.RWMutex
}

func NewCache() *Cache {
	obj := &Cache{
		data: map[string]*MemoryData{},
		lock: sync.RWMutex{},
	}
	return obj
}

func (m *Cache) Get(ctx context.Context, key string) (string, error) {
	var val string
	if err := m.GetObj(ctx, key, &val); err != nil {
		return "", err
	}
	return val, nil
}

func (m *Cache) GetObj(ctx context.Context, key string, obj interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if md, ok := m.data[key]; ok {
		if md.ttl != cache.NoneDuration {
			if time.Now().Sub(md.createTime) > md.ttl {
				delete(m.data, key)
				return cache.ErrKeyNotFound
			}
		}

		bt, _ := json.Marshal(md.val)
		err := json.Unmarshal(bt, obj)
		if err != nil {
			return err
		}
		return nil
	}

	return cache.ErrKeyNotFound
}

func (m *Cache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	errs := make([]string, 0, len(keys))
	rets := make(map[string]string)
	for _, key := range keys {
		val, err := m.Get(ctx, key)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		rets[key] = val
	}
	if len(errs) == 0 {
		return rets, nil
	}
	return rets, errors.New(strings.Join(errs, "||"))
}

func (m *Cache) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return m.SetObj(ctx, key, val, timeout)
}

func (m *Cache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	md := &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        timeout,
	}
	m.data[key] = md
	return nil
}

func (m *Cache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	errs := []string{}
	for k, v := range data {
		err := m.Set(ctx, k, v, timeout)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "||"))
	}
	return nil
}

func (m *Cache) SetForever(ctx context.Context, key string, val string) error {
	return m.Set(ctx, key, val, cache.NoneDuration)
}

func (m *Cache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return m.SetObj(ctx, key, val, cache.NoneDuration)
}

func (m *Cache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if md, ok := m.data[key]; ok {
		md.ttl = timeout
		return nil
	}
	return cache.ErrKeyNotFound
}

func (m *Cache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if md, ok := m.data[key]; ok {
		return md.ttl, nil
	}
	return cache.NoneDuration, cache.ErrKeyNotFound
}

func (m *Cache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	var val int64
	err := m.GetObj(ctx, key, &val)
	val = val + step
	if err == nil {
		m.data[key].val = val
		return val, nil
	}

	if !errors.Is(err, cache.ErrKeyNotFound) {
		return 0, err
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	// key not found
	m.data[key] = &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        cache.NoneDuration,
	}

	return val, nil
}

func (m *Cache) Increment(ctx context.Context, key string) (int64, error) {
	return m.Calc(ctx, key, 1)
}

func (m *Cache) Decrement(ctx context.Context, key string) (int64, error) {
	return m.Calc(ctx, key, -1)
}

func (m *Cache) Del(ctx context.Context, key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.data, key)
	return nil
}

func (m *Cache) DelMany(ctx context.Context, keys []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}

func (m *Cache) Contains(key string) bool {
	_, err := m.Get(context.Background(), key)
	if err != nil {
		return false
	}
	return err == nil
}
