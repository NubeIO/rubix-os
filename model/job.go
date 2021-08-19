package model

import (
	"time"
)


type Job struct {
	UUID          				string    `json:"uuid" sql:"uuid"`
	Frequency   				string    `json:"frequency,omitempty" sql:"frequency"`
	StartDate   				time.Time `json:"start_date,omitempty" sql:"start_date"`
	EndDate     				time.Time `json:"end_date,omitempty" sql:"end_date"`
	IsActive    				bool      `json:"is_active" sql:"is_active"`
	DestroyAfterCompleted   	bool      `json:"destroy_after_completed" sql:"destroy_after_completed"`
}