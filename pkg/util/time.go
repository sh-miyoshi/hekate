package util

import "time"

// TimeTruncate returns truncated time
// output is only year, month, and day (not includes hour, minutes, seconds)
func TimeTruncate(in time.Time) time.Time {
	return time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, in.Location())
}
