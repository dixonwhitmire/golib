package loglib

import (
	"testing"
	"time"
)

// TestSetDefault validates [SetDefaultLogger]
func TestSetDefault(t *testing.T) {
	SetDefaultLogger()
}

func TestLogElapsedTime(t *testing.T) {
	startTime := time.Now()
	defer LogElapsedTime("TestLogElapsedTime", startTime)
	func() {
		for i := 0; i < 10_000; i++ {
			i *= i
		}
	}()
}
