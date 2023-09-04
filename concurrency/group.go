// 懒加载对象容器

package concurrency

import "sync"

// Group is a lazy load container.
type Group struct {
	new  func() interface{}
	data map[string]interface{}
	sync.RWMutex
}

// NewGroup news a group container.
func NewGroup(new func() interface{}) *Group {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	return &Group{
		new:  new,
		data: make(map[string]interface{}),
	}
}

// Get gets the object by the given key.
func (g *Group) Get(key string) interface{} {
	g.RLock()
	obj, ok := g.data[key]
	if ok {
		g.RUnlock()
		return obj
	}
	g.RUnlock()

	// double check
	g.Lock()
	defer g.Unlock()
	obj, ok = g.data[key]
	if ok {
		return obj
	}
	obj = g.new()
	g.data[key] = obj
	return obj
}

// Reset resets the new function and deletes all existing objects.
func (g *Group) Reset(new func() interface{}) {
	if new == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	g.Lock()
	g.new = new
	g.Unlock()
	g.Clear()
}

// Clear deletes all objects.
func (g *Group) Clear() {
	g.Lock()
	g.data = make(map[string]interface{})
	g.Unlock()
}
