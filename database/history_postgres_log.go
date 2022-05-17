package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// GetLastSyncHistoryPostgresLog returns the last sync history postgres log if it exists in history table also.
func (d *GormDatabase) GetLastSyncHistoryPostgresLog() (*model.HistoryPostgresLog, error) {
	var logsModel *model.HistoryPostgresLog
	query := d.DB
	query.First(&logsModel)
	if logsModel != nil {
		var historiesModel *model.History
		query.Where("id = ?", logsModel.ID).
			Where("uuid = ?", logsModel.UUID).
			Where("value = ?", logsModel.Value).
			Where("datetime(timestamp) = datetime(?)", logsModel.Timestamp).
			First(&historiesModel)
		if historiesModel.ID != 0 {
			return logsModel, nil
		}
	}
	return nil, nil
}

// UpdateHistoryPostgresLog update or create if not found a thing.
func (d *GormDatabase) UpdateHistoryPostgresLog(body *model.HistoryPostgresLog) (*model.HistoryPostgresLog, error) {
	var historyPostgresLogModel *model.HistoryPostgresLog
	query := d.DB.First(&historyPostgresLogModel)
	if historyPostgresLogModel.ID == 0 {
		if err := d.DB.Create(&body).Error; err != nil {
			return nil, err
		}
		return body, nil
	}
	query = d.DB.Model(&historyPostgresLogModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyPostgresLogModel, nil
}
