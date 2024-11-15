package parallel

import (
	"errors"
	"sync"
)

type Parallel struct {
	Callbacks []func() interface{}
	channel   chan int
	status    ParallelStatus
}

type ParallelStatus int

const (
	NORMAL = iota
	LISTENING
	STOPPED
	GRACEFUL_STOP
)

func NewParallel(concurrent int) *Parallel {
	return &Parallel{
		Callbacks: make([]func() interface{}, 0),
		channel:   make(chan int, concurrent),
		status:    NORMAL,
	}
}

func (p *Parallel) Add(callback func() interface{}) error {
	if p.IsStopped() {
		return errors.New("the process has stopped")
	} else {
		p.Callbacks = append(p.Callbacks, callback)
		return nil
	}
}

func (p *Parallel) IsStopped() bool {
	return p.status == STOPPED || p.status == GRACEFUL_STOP
}

func (p *Parallel) Wait() (results map[int]interface{}) {
	queues := p.Callbacks
	p.Clear()

	wg := sync.WaitGroup{}
	wg.Add(len(queues))
	resultMutex := sync.RWMutex{}

	results = map[int]interface{}{}
	for key, callback := range queues {
		p.channel <- 0
		go func(key int, callback func() interface{}) {
			// 捕捉异常
			defer func() {
				if err := recover(); err != nil {
					resultMutex.Lock()
					results[key] = err
					resultMutex.Unlock()
				}

				<-p.channel
				wg.Done()
			}()

			result := callback()

			resultMutex.Lock()
			results[key] = result
			resultMutex.Unlock()
		}(key, callback)
	}

	wg.Wait()

	return
}

func (p *Parallel) Run() map[int]interface{} {
	return p.Wait()
}

func (p *Parallel) Stop() {
	p.status = STOPPED
}

func (p *Parallel) GracefulStop() {
	p.status = GRACEFUL_STOP
}

func (p *Parallel) Listen() (err error) {
	p.status = LISTENING

	defer func() {
		if p.status == GRACEFUL_STOP {
			p.Wait()
		}
	}()

	for {
		if p.status == LISTENING {
			p.Wait()
		} else {
			break
		}
	}

	return err
}

func (p *Parallel) Clear() {
	p.Callbacks = make([]func() interface{}, 0)
}
