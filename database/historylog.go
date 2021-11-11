package database

import (
	"github.com/NubeDev/flow-framework/model"
)

// GetHistoryLogByFlowNetworkCloneUUID return history log for the given fncUuid or nil..
func (d *GormDatabase) GetHistoryLogByFlowNetworkCloneUUID(fncUuid string) (*model.HistoryLog, error) {
	var historyLogModel *model.HistoryLog
	d.DB.Where("flow_network_clone_uuid = ?", fncUuid).First(&historyLogModel)
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
	query := d.DB.Where("flow_network_clone_uuid = ?", body.FlowNetworkCloneUUID).First(&historyLogModel)
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
