package loglib

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const (
	LogErrorKey       = "err"
	LogEventSourceKey = "event_source"
	LogElapsedTimeKey = "elapsed_time"
	LogPathKey        = "path"
	LogTimeFormat     = "%.6f"
)

// SetDefaultLogger configures the default loglib for the consuming application.
// The default logger logs INFO events and higher.
func SetDefaultLogger() {
	ConfigureLogger(&slog.HandlerOptions{Level: slog.LevelInfo})
}

// ConfigureLogger configures the default logger given the provided handlerOptions.
func ConfigureLogger(handlerOptions *slog.HandlerOptions) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	slog.SetDefault(logger)
}

// LogElapsedTime calculates the elapsed time for an event using a start time.
// The event being measured could be a function execution, or a more granular processing event.
// LogElapsedTime is meant to be integrated as a deferred function.
//
// Example:
//
//	func someFunction(x, y int) int {
//	   startTime := time.Now()
//	   defer LogElapsedTime(slog.LevelInfo, "someFunction", startTime)
//	   return x + y
//	}
func LogElapsedTime(level slog.Level, eventName string, startTime time.Time) {
	elapsedTime := time.Since(startTime).Seconds()
	slog.Log(
		context.Background(),
		level,
		"elapsed time",
		LogEventSourceKey, eventName,
		LogElapsedTimeKey, fmt.Sprintf(LogTimeFormat, elapsedTime))
}
