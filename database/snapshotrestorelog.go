package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetSnapshotRestoreLogs(hostUUID string) ([]*model.SnapshotRestoreLog, error) {
	var snapshotRestoreLogsModel []*model.SnapshotRestoreLog
	if err := d.DB.Where("host_uuid = ?", hostUUID).Find(&snapshotRestoreLogsModel).Error; err != nil {
		return nil, err
	}
	return snapshotRestoreLogsModel, nil
}

func (d *GormDatabase) CreateSnapshotRestoreLog(body *model.SnapshotRestoreLog) (*model.SnapshotRestoreLog, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.SnapshotRestoreLog)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateSnapshotRestoreLog(uuid string, body *model.SnapshotRestoreLog) (*model.SnapshotRestoreLog, error) {
	var snapshotRestoreLogModel *model.SnapshotRestoreLog
	if err := d.DB.Where("uuid = ?", uuid).Find(&snapshotRestoreLogModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return snapshotRestoreLogModel, nil
}

func (d *GormDatabase) DeleteSnapshotRestoreLog(uuid string) (bool, error) {
	var snapshotRestoreLogModel *model.SnapshotRestoreLog
	query := d.DB.Where("uuid = ? ", uuid).Delete(&snapshotRestoreLogModel)
	return d.deleteResponseBuilder(query)
}
