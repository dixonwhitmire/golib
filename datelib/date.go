// Package datelib formats time.Time values to ISO8601 formats.
package datelib

import "time"

// iso8601Date provides a format string for YYYY-MM-DD.
const iso8601Date = "2006-01-02"

// iso8601DateTime provides a format string for YYYY-MM-DDTHH:MI:SS including timezone.
const iso8601DateTime = "2006-01-02T15:04:05Z"

// ReferenceDate is the canonical Go reference date.
const ReferenceDate = "Mon Jan 2 15:04:05 MST 2006"

// FormatIso8601Date returns a string formatted as YYYY-MM-DD.
func FormatIso8601Date(t time.Time) string {
	return t.Format(iso8601Date)
}

// FormatIso8601DateTime returns a string formatted as YYYY-MM-DDTHH:MI:SS.
func FormatIso8601DateTime(t time.Time) string {
	return t.Format(iso8601DateTime)
}
