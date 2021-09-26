package influxmodel

import (
	"time"
)

type Temperature struct {
	City             string  `json:"city"`
	Country          string  `json:"country"`
	TemperatureScale string  `json:"temperature_scale"`
	TemperatureValue float64 `json:"temperature_value"`
}

type FluxTemperature struct {
	City             string    `json:"city"`
	Country          string    `json:"country"`
	TemperatureScale string    `json:"temperature_scale"`
	Time             time.Time `json:"_time"`
	Start            time.Time `json:"_start"`
	Stop             time.Time `json:"_stop"`
	Measurement      string    `json:"_measurement"`
	Field            string    `json:"_field"`
	Value            float64   `json:"_value"`
	Table            int       `json:"table"`
}
