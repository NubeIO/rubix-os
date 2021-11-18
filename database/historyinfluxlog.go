package database

import (
	"github.com/NubeIO/flow-framework/model"
	"time"
)

// UpdateHistoryInfluxLogLastSyncId update a thing.
func (d *GormDatabase) UpdateHistoryInfluxLogLastSyncId(lastSyncId int) (*model.HistoryInfluxLog, error) {
	var historyInfluxLogModel *model.HistoryInfluxLog
	query := d.DB.First(&historyInfluxLogModel)
	historyInfluxLogModel.LastSyncID = lastSyncId
	historyInfluxLogModel.Timestamp = time.Now()
	if query.RowsAffected == 0 {
		d.DB.Create(&historyInfluxLogModel)
	}
	d.DB.Save(&historyInfluxLogModel)
	return historyInfluxLogModel, nil
}
