package threading

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/apus-run/sea-kit/lang"
	"github.com/stretchr/testify/assert"
)

func TestRoutineId(t *testing.T) {
	assert.True(t, RoutineId() > 0)
}

func TestRunSafe(t *testing.T) {
	log.SetOutput(io.Discard)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan lang.PlaceholderType)
	go RunSafe(func() {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestRunSafeCtx(t *testing.T) {
	ctx := context.Background()
	ch := make(chan lang.PlaceholderType)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	go RunSafeCtx(ctx, func() {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestGoSafeCtx(t *testing.T) {
	ctx := context.Background()
	ch := make(chan lang.PlaceholderType)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	GoSafeCtx(ctx, func() {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}
