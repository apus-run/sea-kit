package lru

import (
	"container/list"
	"fmt"
	"sync"
)

type LRUCache struct {
	lock     sync.Mutex
	capacity int                           // maximum number of key-value pairs
	cache    map[interface{}]*list.Element // safemap for cached key-value pairs
	lru      *list.List                    // LRU list
}

type Pair struct {
	key   interface{} // tiny_cache key
	value interface{} // tiny_cache value
}

// NewLRUCache returns a new, empty LRUCache
func NewLRUCache(capacity int) *LRUCache {
	c := new(LRUCache)
	c.capacity = capacity
	c.cache = make(map[interface{}]*list.Element)
	c.lru = list.New()
	return c
}

// Get get cached value from LRU tiny_cache
// The second return value indicates whether key is found or not, true if found, false if not
func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.lru.MoveToFront(elem) // move node to head of lru list
		return elem.Value.(*Pair).value, true
	}
	return nil, false
}

// Add adds a key-value pair to LRU tiny_cache, true if eviction occurs, false if not
func (c *LRUCache) Add(key interface{}, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	// update item if found in tiny_cache
	if elem, ok := c.cache[key]; ok {
		c.lru.MoveToFront(elem) // update lru list
		elem.Value.(*Pair).value = value
		return false
	}
	// add item if not found
	elem := c.lru.PushFront(&Pair{key, value})
	c.cache[key] = elem
	// evict item if needed
	if c.lru.Len() > c.capacity {
		c.evict()
		return true
	}
	return false
}

// evict a key-value pair from LRU cache
func (c *LRUCache) evict() {
	elem := c.lru.Back()
	if elem == nil {
		return
	}
	// remove item at the end of lru list
	c.lru.Remove(elem)
	delete(c.cache, elem.Value.(*Pair).key)
}

// Del deletes cached value from tiny_cache
func (c *LRUCache) Del(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.lru.Remove(elem)
		delete(c.cache, key)
	}
}

// Len returns number of items in tiny_cache
func (c *LRUCache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Len()
}

// Keys returns keys of items in tiny_cache
func (c *LRUCache) Keys() []interface{} {
	var keyList []interface{}
	c.lock.Lock()
	for key := range c.cache {
		keyList = append(keyList, key)
	}
	c.lock.Unlock()
	return keyList
}

// EnlargeCapacity enlarges the capacity of tiny_cache
func (c *LRUCache) EnlargeCapacity(newCapacity int) error {
	// lock
	c.lock.Lock()
	defer c.lock.Unlock()
	// check newCapacity
	if newCapacity < c.capacity {
		return fmt.Errorf("newCapacity[%d] must be larger than current[%d]",
			newCapacity, c.capacity)
	}
	c.capacity = newCapacity
	return nil
}
