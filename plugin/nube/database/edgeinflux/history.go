package main

import "time"

// History for storing the raw history data
type History struct {
	UUID      string
	Value     float64
	Timestamp time.Time
	Tags      map[string]string
}
