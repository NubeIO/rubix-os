package database

import (
	"github.com/NubeDev/flow-framework/model"
	"time"
)

// CreateHistoryLog creates a thing.
func (d *GormDatabase) CreateHistoryLog(body *model.HistoryLog) (*model.HistoryLog, error) {
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

// UpdateHistoryLogLastSyncId update a thing.
func (d *GormDatabase) UpdateHistoryLogLastSyncId(lastSyncId int) (*model.HistoryLog, error) {
	var historyLogModel *model.HistoryLog
	query := d.DB.First(&historyLogModel)
	historyLogModel.LastSyncID = lastSyncId
	historyLogModel.Timestamp = time.Now()
	if query.RowsAffected == 0 {
		d.DB.Create(&historyLogModel)
	}
	d.DB.Save(&historyLogModel)
	return historyLogModel, nil
}
