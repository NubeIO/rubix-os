package model

import "time"

//History for storing the all history
type History struct {
	ID        int       `json:"id" gorm:"AUTO_INCREMENT;primary_key;index"`
	UUID      string    `json:"uuid"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
