package mediator_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apus-run/sea-kit/mediator"
)

type testEvent struct {
	called int32
}

func (e *testEvent) Kind() mediator.EventKind {
	return "test-event"
}

type testHandler struct{}

func (h testHandler) Listening() []mediator.EventKind {
	return []mediator.EventKind{"test-event"}
}

func (h testHandler) Handle(_ context.Context, ev mediator.Event) {
	te, ok := ev.(*testEvent)
	if !ok {
		panic("unexpected event type")
	}
	atomic.AddInt32(&te.called, 1)
}

func TestEvent(t *testing.T) {
	mediator, _ := mediator.NewInMemMediator(mediator.Options{Concurrent: 3})
	mediator.Subscribe(testHandler{})

	ev := &testEvent{}
	mediator.Dispatch(ev)
	<-time.After(100 * time.Millisecond)

	if atomic.LoadInt32(&ev.called) != 1 {
		t.FailNow()
	}
}

func TestEventCollection(t *testing.T) {
	eb, _ := mediator.NewInMemMediator(mediator.Options{Concurrent: 3})
	eb.Subscribe(testHandler{})

	collection := mediator.NewEventCollection()
	ev := &testEvent{}
	collection.Add(ev)
	collection.Raise(eb)
	<-time.After(100 * time.Millisecond)

	if atomic.LoadInt32(&ev.called) != 1 {
		t.FailNow()
	}

	collection.Raise(eb)
	<-time.After(100 * time.Millisecond)

	if atomic.LoadInt32(&ev.called) != 1 {
		t.FailNow()
	}
}
