package ttime

import (
	"time"
)

const utc = false // used for when doing dev to appstore as local time

// RealTime is a concrete implementation of Time interface.
type RealTime struct{}

// New initializes and returns a new Time instance.
func New() *RealTime {
	return &RealTime{}
}

// Now returns a timestamp of the current datetime in UTC.
func Now() time.Time {
	if utc {
		return time.Now().UTC()
	}
	return time.Now()
}

// Now returns a timestamp of the current datetime in UTC.
func (rt *RealTime) Now() time.Time {
	if utc {
		return time.Now().UTC()
	}
	return time.Now()
}

// Timestamp returns a timestamp of the current datetime in UTC.
func (rt *RealTime) Timestamp() string {
	if utc {
		return time.Now().UTC().Format(time.RFC3339)
	}
	return time.Now().Format(time.RFC3339)
}

// Pretty returns a timestamp of the current datetime in UTC.
func (rt *RealTime) Pretty() string {
	if utc {
		return time.Now().UTC().Format(time.RFC850)
	}
	return time.Now().Format(time.RFC850)
}
