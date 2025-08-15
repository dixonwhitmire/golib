package ctxlib

import (
	"context"
	"syscall"
	"testing"
	"time"
)

func createTestingContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithCancel(context.Background())
}

// TestNewSignalContext validates that a context is cancelled if a specified signal is received.
// The test uses SIGUSR1 for cancellation to ensure we don't have unintended testing side effects with SIGTERM or SIGKILL.
func TestNewSignalContext(t *testing.T) {
	parentCtx, parentCancelFunc := createTestingContext(t)
	t.Cleanup(parentCancelFunc)

	testSignal := syscall.SIGUSR1
	ctx := NewSignalContext(parentCtx, testSignal)

	go func() {
		// ensure signal handler is ready
		time.Sleep(100 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), testSignal)
		if err != nil {
			t.Errorf("failed to send signal: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second):
		// Failure: The context did not cancel within the timeout
		t.Fatal("Test timed out, context was not canceled by the signal")
	}
}

// TestNewSignalContext_ParentCancellation validates that a context is cancelled if the parent context is cancelled.
// The test uses SIGUSR1 for cancellation to ensure we don't have unintended testing side effects with SIGTERM or SIGKILL.
func TestNewSignalContext_ParentCancellation(t *testing.T) {
	// Create a parent context that we can cancel
	parentCtx, parentCancel := context.WithCancel(context.Background())

	testSignal := syscall.SIGUSR1
	ctx := NewSignalContext(parentCtx, testSignal)

	// Cancel the parent context
	parentCancel()

	// Use a select statement to wait for the child context to be done or for a timeout
	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second):
		// Failure: The context did not cancel within the timeout
		t.Fatal("Test timed out, context was not canceled by the parent")
	}
}
