package wg

import (
	"sync"

	"github.com/apus-run/sea-kit/concurrency/panics"
)

// wgWrapper wrap go WaitGroup
type wgWrapper struct {
	wg sync.WaitGroup
	pc panics.Catcher
}

// New create wgWrapper instance
func New() *wgWrapper {
	return &wgWrapper{}
}

// Wrap fn func in goroutine to run without recovery func
func (w *wgWrapper) Wrap(fn func()) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.pc.Try(fn)
	}()
}

// Go fn func in goroutine to run without recovery func
func (w *wgWrapper) Go(fn func()) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.pc.Try(fn)
	}()
}

// WaitAndRecover will block until all goroutines spawned with Go exit and
// will return a *panics.Recovered if one of the child goroutines panics.
func (w *wgWrapper) WaitAndRecover() *panics.Recovered {
	w.wg.Wait()

	// Return a recovered panic if we caught one from a child goroutine.
	return w.pc.Recovered()
}

// Wait blocks until the WaitGroup counter is zero.
func (w *wgWrapper) Wait() {
	w.wg.Wait()

	// Propagate a panic if we caught one from a child goroutine.
	w.pc.Repanic()
}
