package recovery

import (
	"fmt"
	"runtime/debug"

	"github.com/elliot14A/fincher/pkg/logger"
)

// SafeGo runs fn in a new goroutine with panic protection (delegating to WrapPanic).
func SafeGo(name string, fn func(), onPanic func(r any, stack string)) {
	go func() {
		WrapPanic(name, fn, onPanic)
	}()
}

// WrapPanic executes fn synchronously with panic protection.
// If fn panics, the panic is recovered, logged with full stack trace, and onPanic (if provided) is invoked.
func WrapPanic(name string, fn func(), onPanic func(r any, stack string)) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			logger.Error("panic recovered",
				"operation", name,
				"panic", fmt.Sprint(r),
				"stack", stack,
			)
			if onPanic != nil {
				onPanic(r, stack)
			}
		}
	}()
	fn()
}
