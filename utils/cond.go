package utils

import (
	"sync"
	"time"
)

// Cond is similar to sync.Cond, but you can wait with a timeout.
type Cond struct {
	L      sync.Locker
	signal chan struct{}
}

// Signal wakes up a waiter.
// It's required for the caller to hold L.
func (c *Cond) Signal() {
	select {
	case c.signal <- struct{}{}:
	default:
	}
}

// Broadcast wake up all the waiters.
// It's required for the caller to hold L.
func (c *Cond) Broadcast() {
	close(c.signal)
	c.signal = make(chan struct{})
}

var timerPool = sync.Pool{
	New: func() interface{} {
		return time.NewTimer(time.Second)
	},
}

// WaitWithTimeout wait for a signal or a period of timeout eclipsed.
// returns true in case of timeout else false
func (c *Cond) WaitWithTimeout(d time.Duration) bool {
	ch := c.signal
	c.L.Unlock()
	t := timerPool.Get().(*time.Timer)
	t.Reset(d)
	defer func() {
		t.Stop()
		timerPool.Put(t)
		c.L.Lock()
	}()
	select {
	case <-ch:
		return false
	case <-t.C:
		return true
	}
}

// NewCond creates a Cond.
func NewCond(lock sync.Locker) *Cond {
	return &Cond{lock, make(chan struct{})}
}
