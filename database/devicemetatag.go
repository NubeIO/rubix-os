package database

import (
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm/clause"
	"strings"
)

func (d *GormDatabase) CreateDeviceMetaTags(deviceUUID string, body []*model.DeviceMetaTag) ([]*model.DeviceMetaTag,
	error) {
	tx := d.DB.Begin()
	var deviceUUIDs []string
	for _, b := range body {
		var count int64
		tx.Model(&model.DeviceMetaTag{}).Where("device_uuid = ?", deviceUUID).Where("uuid = ?", b.UUID).
			Count(&count)
		if count == 0 {
			b.UUID = nuuid.MakeTopicUUID(model.CommonNaming.DeviceMetaTag)
		} else {
			deviceUUIDs = append(deviceUUIDs, b.UUID)
		}
		b.DeviceUUID = deviceUUID
	}
	notIn := strings.Join(deviceUUIDs, ",")
	if err := tx.Where("device_uuid = ?", deviceUUID).Where("uuid not in (?)", notIn).
		Delete(&model.DeviceMetaTag{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(body) > 0 {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(body).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return body, nil
}
