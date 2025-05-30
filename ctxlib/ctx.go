package ctxlib

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// SignalContext returns a context which responds to a SIGTERM signal.
func SignalContext() context.Context {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		slog.Info("received shutdown signal . . . stopping application")
		stop()
	}()
	return ctx
}
