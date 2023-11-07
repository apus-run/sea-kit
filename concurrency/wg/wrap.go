package wg

import (
	"sync"
)

// wgWrapper wrap go WaitGroup
type wgWrapper struct {
	wg sync.WaitGroup
}

// New create wgWrapper instance
func New() *wgWrapper {
	w := &wgWrapper{}
	return w
}

// Wrap fn func in goroutine to run without recovery func
func (w *wgWrapper) Wrap(fn func()) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		fn()
	}()
}

// WrapWithRecovery fn func in goroutine to run with customized recovery func
func (w *wgWrapper) WrapWithRecovery(fn func(), recoveryFunc func(r interface{})) {
	w.wg.Add(1)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				recoveryFunc(e)
			}
		}()

		defer w.wg.Done()
		fn()
	}()
}

// Wait blocks until the WaitGroup counter is zero.
func (w *wgWrapper) Wait() {
	w.wg.Wait()
}
