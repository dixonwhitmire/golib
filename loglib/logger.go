package loglib

import (
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

// SetDefaultLogger configures the default loglib for the application.
func SetDefaultLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
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
//	   defer LogElapsedTime("someFunction", startTime)
//	   return x + y
//	}
func LogElapsedTime(eventName string, startTime time.Time) {
	elapsedTime := time.Since(startTime).Seconds()
	slog.Info("elapsed time",
		LogEventSourceKey, eventName,
		LogElapsedTimeKey, fmt.Sprintf(LogTimeFormat, elapsedTime))
}
