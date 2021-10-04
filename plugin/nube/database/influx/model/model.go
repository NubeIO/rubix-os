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

type HistPayload struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
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

type InfluxTest struct {
	Org       string
	Bucket    string
	ServerURL string
	AuthToken string
}
