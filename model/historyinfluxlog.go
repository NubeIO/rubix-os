package model

import "time"

//HistoryInfluxLog for storing the history influx logs
type HistoryInfluxLog struct {
	ID         int       `json:"id" gorm:"AUTO_INCREMENT;primary_key;index"`
	LastSyncID int       `json:"last_sync_id"`
	Timestamp  time.Time `json:"timestamp"`
}
