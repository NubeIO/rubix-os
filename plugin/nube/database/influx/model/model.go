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
	ID              int       `json:"id"`
	ProducerUUID    string    `json:"producer_uuid"`
	WriterUUID      string    `json:"writer_uuid"`
	ThingWriterUUID string    `json:"current_writer_uuid"`
	Out16           *float64  `json:"out_16"`
	Timestamp       time.Time `json:"timestamp"`
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
