package threading

import (
	"github.com/apus-run/sea-kit/lang"
)

// A TaskRunner is used to control the concurrency of goroutines.
type TaskRunner struct {
	limitChan chan lang.PlaceholderType
}

// NewTaskRunner returns a TaskRunner.
func NewTaskRunner(concurrency int) *TaskRunner {
	return &TaskRunner{
		limitChan: make(chan lang.PlaceholderType, concurrency),
	}
}

// Schedule schedules a task to run under concurrency control.
func (rp *TaskRunner) Schedule(task func()) {
	rp.limitChan <- lang.Placeholder

	go func() {
		defer Recover(func() {
			<-rp.limitChan
		})

		task()
	}()
}
