package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
)

func (d *GormDatabase) GetSnapshotLog() ([]*model.SnapshotLog, error) {
	var snapshotLogsModel []*model.SnapshotLog
	if err := d.DB.Find(&snapshotLogsModel).Error; err != nil {
		return nil, err
	}
	return snapshotLogsModel, nil
}

func (d *GormDatabase) CreateSnapshotLog(body *model.SnapshotLog) (*model.SnapshotLog, error) {
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateSnapshotLog(file string, body *model.SnapshotLog) (*model.SnapshotLog, error) {
	var snapshotLogsModel []*model.SnapshotLog
	if err := d.DB.Where("file = ?", file).Find(&snapshotLogsModel).Error; err != nil {
		return nil, err
	}
	if len(snapshotLogsModel) == 0 {
		body.File = file
		return d.CreateSnapshotLog(body)
	}
	var snapshotLogModel *model.SnapshotLog
	if err := d.DB.Where("file = ?", file).Find(&snapshotLogModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return snapshotLogModel, nil
}

func (d *GormDatabase) DeleteSnapshotLog(file string) (*interfaces.Message, error) {
	var snapshotLogModel *model.SnapshotLog
	query := d.DB.Where("file = ? ", file).Delete(&snapshotLogModel)
	return d.deleteResponse(query)
}

// DeleteSnapshotLogs avoids discrepancies between snapshot raw files and logs
// discrepancies happens when snapshot gets deleted manually
func (d *GormDatabase) DeleteSnapshotLogs(files []string) (*interfaces.Message, error) {
	var snapshotLogModel *model.SnapshotLog
	query := d.DB.Where("file NOT IN ?", files).Delete(&snapshotLogModel)
	return d.deleteResponse(query)
}
