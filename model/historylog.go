package model

import "time"

//HistoryLog for storing the history logs
type HistoryLog struct {
	ID         int       `json:"id" gorm:"AUTO_INCREMENT;primary_key;index"`
	LastSyncID int       `json:"last_sync_id"`
	Timestamp  time.Time `json:"timestamp"`
}
