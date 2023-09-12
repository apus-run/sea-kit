package wp

import (
	"log"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	count int64
)

type TestWorker struct {
}

func (t *TestWorker) Do(data interface{}) {
	atomic.AddInt64(&count, 1)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
}

func TestOnceRun(t *testing.T) {
	a := assert.New(t)

	jobQueue := make(chan interface{}, 200)
	num := 10000000
	p := &WorkerPool{
		JobQueue:       jobQueue,
		MaxWorkers:     1000,
		InitWorkers:    100,
		MaxIdleWorkers: 100,
		RunI:           &TestWorker{},
	}

	go func() {
		for i := 0; i < num; i++ {
			jobQueue <- i
		}
		p.Stop()
	}()

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			<-ticker.C
			log.Printf("%+v", p.Stats())
		}
	}()

	p.Run()

	a.Equal(count, int64(num))
}

func TestStop(t *testing.T) {
	a := assert.New(t)

	jobQueue := make(chan interface{}, 200)
	p := &WorkerPool{
		JobQueue:   jobQueue,
		MaxWorkers: 10,
		RunI:       &TestWorker{},
	}

	go func() {
		defer func() {
			err := recover()
			a.EqualError(err.(error), "send on closed channel")
		}()
		for {
			jobQueue <- struct {
			}{}
		}
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		p.Stop()
	}()

	p.Run()
}
