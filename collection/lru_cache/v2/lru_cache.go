package lru_cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// LRUCache is a concurrent fixed size cache that evicts elements in LRU order as well as by TTL.
type LRUCache struct {
	mux     sync.Mutex
	lru     *list.List
	cache   map[string]*list.Element
	maxSize int // maximum number of key-value pairs
	ttl     time.Duration
	TimeNow func() time.Time
	onEvict EvictCallback
}

// NewLRUCache creates a new LRU cache with default options.
func NewLRUCache(maxSize int) *LRUCache {
	return NewLRUWithOptions(maxSize, nil)
}

// NewLRUWithOptions creates a new LRU cache with the given options.
func NewLRUWithOptions(maxSize int, opts *Options) *LRUCache {
	if opts == nil {
		opts = &Options{}
	}
	if opts.TimeNow == nil {
		opts.TimeNow = time.Now
	}
	return &LRUCache{
		lru:     list.New(),
		cache:   make(map[string]*list.Element, opts.InitialCapacity),
		ttl:     opts.TTL,
		maxSize: maxSize,
		TimeNow: opts.TimeNow,
		onEvict: opts.OnEvict,
	}
}

// Get retrieves the value stored under the given key
func (c *LRUCache) Get(key string) any {
	c.mux.Lock()
	defer c.mux.Unlock()

	elt := c.cache[key]
	if elt == nil {
		return nil
	}

	cacheEntry := elt.Value.(*cacheEntry)
	if !cacheEntry.expiration.IsZero() && c.TimeNow().After(cacheEntry.expiration) {
		// Entry has expired
		if c.onEvict != nil {
			c.onEvict(cacheEntry.key, cacheEntry.value)
		}
		c.lru.Remove(elt)
		delete(c.cache, cacheEntry.key)
		return nil
	}

	c.lru.MoveToFront(elt)
	return cacheEntry.value
}

// Put puts a new value associated with a given key, returning the existing value (if present)
func (c *LRUCache) Put(key string, value any) any {
	c.mux.Lock()
	defer c.mux.Unlock()
	elt := c.cache[key]
	return c.putWithMutexHold(key, value, elt)
}

// CompareAndSwap puts a new value associated with a given key if existing value matches oldValue.
// It returns itemInCache as the element in cache after the function is executed and replaced as true if value is replaced, false otherwise.
func (c *LRUCache) CompareAndSwap(key string, oldValue, newValue any) (itemInCache any, replaced bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	elt := c.cache[key]
	// If entry not found, old value should be nil
	if elt == nil && oldValue != nil {
		return nil, false
	}

	if elt != nil {
		// Entry found, compare it with that you expect.
		entry := elt.Value.(*cacheEntry)
		if entry.value != oldValue {
			return entry.value, false
		}
	}
	c.putWithMutexHold(key, newValue, elt)
	return newValue, true
}

// putWithMutexHold populates the cache and returns the inserted value.
// Caller is expected to hold the c.mut mutex before calling.
func (c *LRUCache) putWithMutexHold(key string, value any, elt *list.Element) any {
	if elt != nil {
		entry := elt.Value.(*cacheEntry)
		existing := entry.value
		entry.value = value
		if c.ttl != 0 {
			entry.expiration = c.TimeNow().Add(c.ttl)
		}
		c.lru.MoveToFront(elt)
		return existing
	}

	entry := &cacheEntry{
		key:   key,
		value: value,
	}

	if c.ttl != 0 {
		entry.expiration = c.TimeNow().Add(c.ttl)
	}
	c.cache[key] = c.lru.PushFront(entry)
	for len(c.cache) > c.maxSize {
		oldest := c.lru.Remove(c.lru.Back()).(*cacheEntry)
		if c.onEvict != nil {
			c.onEvict(oldest.key, oldest.value)
		}
		delete(c.cache, oldest.key)
	}

	return nil
}

// Del deletes a key, value pair associated with a key
func (c *LRUCache) Del(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	elt := c.cache[key]
	if elt != nil {
		entry := c.lru.Remove(elt).(*cacheEntry)
		if c.onEvict != nil {
			c.onEvict(entry.key, entry.value)
		}
		delete(c.cache, key)
	}
}

// Size returns the number of entries currently in the lru, useful if cache is not full
func (c *LRUCache) Size() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return len(c.cache)
}

func (c *LRUCache) Len() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.Len()
}

// Keys returns keys of items in tiny_cache
func (c *LRUCache) Keys() []interface{} {
	var keyList []interface{}
	c.mux.Lock()
	for key := range c.cache {
		keyList = append(keyList, key)
	}
	c.mux.Unlock()
	return keyList
}

// EnlargeCapacity enlarges the capacity of tiny_cache
func (c *LRUCache) EnlargeCapacity(newCapacity int) error {
	// lock
	c.mux.Lock()
	defer c.mux.Unlock()
	// check newCapacity
	if newCapacity < c.maxSize {
		return fmt.Errorf("newCapacity[%d] must be larger than current[%d]",
			newCapacity, c.maxSize)
	}
	c.maxSize = newCapacity
	return nil
}

type cacheEntry struct {
	key        string
	expiration time.Time
	value      any
}
