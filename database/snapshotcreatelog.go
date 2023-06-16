package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetSnapshotCreateLogs(hostUUID string) ([]*model.SnapshotCreateLog, error) {
	var snapshotCreateLogModel []*model.SnapshotCreateLog
	if err := d.DB.Where("host_uuid = ?", hostUUID).Find(&snapshotCreateLogModel).Error; err != nil {
		return nil, err
	}
	return snapshotCreateLogModel, nil
}

func (d *GormDatabase) CreateSnapshotCreateLog(body *model.SnapshotCreateLog) (*model.SnapshotCreateLog, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.SnapshotCreateLog)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateSnapshotCreateLog(uuid string, body *model.SnapshotCreateLog) (*model.SnapshotCreateLog, error) {
	var snapshotCreateLog *model.SnapshotCreateLog
	if err := d.DB.Where("uuid = ?", uuid).Find(&snapshotCreateLog).Updates(body).Error; err != nil {
		return nil, err
	}
	return snapshotCreateLog, nil
}

func (d *GormDatabase) DeleteSnapshotCreateLog(uuid string) (*interfaces.Message, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.SnapshotCreateLog{})
	return d.deleteResponse(query)
}
