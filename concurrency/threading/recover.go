package threading

import (
	"context"

	"log"
	"runtime/debug"
)

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() {})
func Recover(cleanups ...func()) {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v\n", r)
		for _, cleanup := range cleanups {
			cleanup()
		}
	}
}

// RecoverCtx is used with defer to do cleanup on panics.
func RecoverCtx(ctx context.Context, cleanups ...func()) {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v\n", r)
		log.Printf("Stack trace: %s\n", debug.Stack())
		for _, cleanup := range cleanups {
			cleanup()
		}
	}
}
