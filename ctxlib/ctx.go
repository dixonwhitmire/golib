package ctxlib

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// NewSignalContext returns a new context that is cancelled when one of the specified signals is received.
// If no signals are provided, it defaults to listening for os.Interrupt and syscall.SIGTERM.
func NewSignalContext(parent context.Context, signals ...os.Signal) context.Context {
	if len(signals) == 0 {
		signals = []os.Signal{os.Interrupt, syscall.SIGTERM}
	}

	ctx, stop := signal.NotifyContext(parent, signals...)

	slog.Info("signal registration complete", "signals", signals)

	go func() {
		<-ctx.Done()
		slog.Info("received shutdown signal. stopping application.")
		stop()
	}()

	return ctx
}
