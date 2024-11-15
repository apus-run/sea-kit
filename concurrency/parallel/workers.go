package parallel

import (
	"errors"
	"log"
)

type Workers struct {
	processes int
	pipe      chan bool
	stopped   bool
}

func NewWorkers(num int) *Workers {
	return &Workers{
		processes: num,
		pipe:      make(chan bool, num),
	}
}

func (workers *Workers) Handle(job func()) error {
	if workers.stopped {
		return errors.New("the process has stopped")
	}
	workers.pipe <- true

	go func() {
		defer func() {
			<-workers.pipe
			if panicValue := recover(); panicValue != nil {
				log.Default().Println(panicValue)
			}
		}()
		job()
	}()

	return nil
}

func (workers *Workers) Stop() {
	workers.stopped = true
}
