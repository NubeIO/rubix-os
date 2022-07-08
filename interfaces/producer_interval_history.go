package interfaces

import "time"

type ProducerIntervalHistory struct {
	UUID               string     `json:"uuid"`
	ProducerThingClass string     `json:"producer_thing_class"`
	HistoryInterval    *int       `json:"history_interval,omitempty"`
	Timestamp          *time.Time `json:"timestamp,omitempty"`
	PresentValue       *float64   `json:"present_value,omitempty"`
}
