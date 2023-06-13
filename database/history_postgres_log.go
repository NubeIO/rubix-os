package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// GetHistoryPostgresLogLastSyncHistoryId returns the last history postgres log last sync history id if it exists in history table also.
func (d *GormDatabase) GetHistoryPostgresLogLastSyncHistoryId() (int, error) {
	var logsModel *model.HistoryPostgresLog
	query := d.DB
	query.First(&logsModel)
	if logsModel != nil {
		var historiesModel *model.History
		query.Where("id = ?", logsModel.ID).
			Where("point_uuid = ?", logsModel.PointUUID).
			Where("host_uuid = ?", logsModel.HostUUID).
			Where("datetime(timestamp) = datetime(?)", logsModel.Timestamp).
			First(&historiesModel)
		if query.Error != nil {
			return 0, query.Error
		}
		return historiesModel.HistoryID, nil
	}
	return 0, nil
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
