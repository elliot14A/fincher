package recovery_test

import (
	"sync"
	"testing"

	"github.com/elliot14A/fincher/pkg/recovery"
)

func TestSafeGo_RecoversPanic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	var recoveredVal any
	var capturedStack string

	recovery.SafeGo("test-routine", func() {
		panic("simulated routine panic")
	}, func(r any, stack string) {
		recoveredVal = r
		capturedStack = stack
		wg.Done()
	})

	wg.Wait()

	if recoveredVal != "simulated routine panic" {
		t.Errorf("expected panic value 'simulated routine panic', got: %v", recoveredVal)
	}
	if capturedStack == "" {
		t.Error("expected non-empty captured stack trace")
	}
}

func TestWrapPanic_RecoversPanic(t *testing.T) {
	var recoveredVal any

	recovery.WrapPanic("test-wrap", func() {
		panic("sync panic")
	}, func(r any, stack string) {
		recoveredVal = r
	})

	if recoveredVal != "sync panic" {
		t.Errorf("expected panic value 'sync panic', got: %v", recoveredVal)
	}
}
