package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm/clause"
)

// GetHistoryLogByHostUUID return history log for the given fncUuid or nil..
func (d *GormDatabase) GetHistoryLogByHostUUID(hostUUID string) (*model.HistoryLog, error) {
	var historyLogModel *model.HistoryLog
	d.DB.Where("host_uuid = ?", hostUUID).First(&historyLogModel)
	return historyLogModel, nil
}

// CreateHistoryLog creates a thing.
func (d *GormDatabase) CreateHistoryLog(body *model.HistoryLog) (*model.HistoryLog, error) {
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

// UpdateHistoryLog update/create a thing.
func (d *GormDatabase) UpdateHistoryLog(body *model.HistoryLog) (*model.HistoryLog, error) {
	var historyLogModel *model.HistoryLog
	query := d.DB.Where("host_uuid = ?", body.HostUUID).First(&historyLogModel)
	if historyLogModel.ID == 0 {
		if err := d.DB.Create(&body).Error; err != nil {
			return nil, err
		}
		return body, nil
	}
	query = d.DB.Model(&historyLogModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return historyLogModel, nil
}

// UpdateBulkHistoryLogs update/create a thing in a bulk.
func (d *GormDatabase) UpdateBulkHistoryLogs(body []*model.HistoryLog) (bool, error) {
	if err := d.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(body, 1000).Error; err != nil {
		return false, err
	}
	return true, nil
}
