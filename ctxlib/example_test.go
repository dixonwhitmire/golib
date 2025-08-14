package ctxlib

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"
)

func ExampleNewSignalContext() {
	// Example of using a custom signal
	// You can replace syscall.SIGUSR1 with any signal you'd like to test with.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a new context that listens for SIGUSR1.
	// In a real-world scenario, this would be syscall.SIGINT and syscall.SIGTERM.
	sigCtx := NewSignalContext(ctx, syscall.SIGUSR1)

	// Simulate receiving the signal after 1 second.
	// This goroutine mimics an external signal being sent.
	go func() {
		// A real application would not do this. It is only for demonstration purposes.
		time.Sleep(1 * time.Second)
		fmt.Println("Sending signal...")
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGUSR1)
	}()

	// Wait for the context to be cancelled.
	// This will happen when the SIGUSR1 signal is received.
	<-sigCtx.Done()

	// Print the cancellation reason.
	// The output here is what `go doc` expects to see.
	fmt.Printf("Context cancelled with error: %v\n", sigCtx.Err())

	// Un-comment this to see the output.
	// Output:
	// Sending signal...
	// Context cancelled with error: context canceled
}

func ExampleNewSignalContext_default() {
	// The default behavior listens for os.Interrupt and syscall.SIGTERM.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a new context that listens for default signals.
	sigCtx := NewSignalContext(ctx)

	// In a real-world scenario, you would press Ctrl+C to trigger cancellation.
	// Here, we simulate it for the example.
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Sending interrupt signal...")
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(os.Interrupt)
	}()

	<-sigCtx.Done()

	fmt.Printf("Context cancelled with error: %v\n", sigCtx.Err())

	// Output:
	// Sending interrupt signal...
	// Context cancelled with error: context canceled
}
