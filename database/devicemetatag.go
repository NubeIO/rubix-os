package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

func CreateDeviceMetaTagsTransaction(tx *gorm.DB, deviceUUID string, body []*model.DeviceMetaTag) ([]*model.DeviceMetaTag, error) {
	var keys []string
	for _, b := range body {
		var count int64
		tx.Model(&model.DeviceMetaTag{}).Where("device_uuid = ?", deviceUUID).Where("key = ?", b.Key).
			Count(&count)
		if count > 0 {
			keys = append(keys, b.Key)
		}
		b.DeviceUUID = deviceUUID
	}
	notIn := strings.Join(keys, ",")
	if err := tx.Where("device_uuid = ?", deviceUUID).Where("key not in (?)", notIn).
		Delete(&model.DeviceMetaTag{}).Error; err != nil {
		return nil, err
	}
	if len(body) > 0 {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(body).Error; err != nil {
			return nil, err
		}
	}
	if err := tx.Where("device_uuid = ?", deviceUUID).Find(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) CreateDeviceMetaTags(deviceUUID string, body []*model.DeviceMetaTag) ([]*model.DeviceMetaTag, error) {
	tx := d.DB.Begin()
	mt, err := CreateDeviceMetaTagsTransaction(tx, deviceUUID, body)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return mt, nil
}

func (d *GormDatabase) GetDeviceMetaTags() ([]*model.DeviceMetaTag, error) {
	var deviceMetaTagsModel []*model.DeviceMetaTag
	if err := d.DB.Find(&deviceMetaTagsModel).Error; err != nil {
		return nil, err
	}
	return deviceMetaTagsModel, nil
}

func (d *GormDatabase) GetDevicesMetaTagsForPostgresSync() ([]*model.DeviceMetaTag, error) {
	var deviceMetaTagsModel []*model.DeviceMetaTag
	query := d.DB.Table("device_meta_tags").
		Select("devices.source_uuid AS device_uuid, device_meta_tags.key, device_meta_tags.value").
		Joins("INNER JOIN devices ON devices.uuid = device_meta_tags.device_uuid").
		Where("IFNULL(devices.source_uuid,'') != ''").
		Scan(&deviceMetaTagsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return deviceMetaTagsModel, nil
}
