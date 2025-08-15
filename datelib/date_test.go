package datelib

import (
	"testing"
	"time"
)

// createTestDate returns a time.Time configured for 8/15/25 at 2PM UTC.
func createTestDate(t *testing.T) time.Time {
	t.Helper()
	return time.Date(2025, time.August, 15, 14, 0, 0, 0, time.UTC)
}

func TestFormatIso8601Date(t *testing.T) {
	const want = "2025-08-15"
	got := FormatIso8601Date(createTestDate(t))

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatIso8601DateTime(t *testing.T) {
	const want = "2025-08-15T14:00:00Z"
	got := FormatIso8601DateTime(createTestDate(t))

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
