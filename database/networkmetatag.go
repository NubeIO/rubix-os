package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

func CreateNetworkMetaTagsTransaction(tx *gorm.DB, networkUUID string, body []*model.NetworkMetaTag) ([]*model.NetworkMetaTag, error) {
	var keys []string
	for _, b := range body {
		var count int64
		tx.Model(&model.NetworkMetaTag{}).Where("network_uuid = ?", networkUUID).Where("key = ?",
			b.Key).Count(&count)
		if count > 0 {
			keys = append(keys, b.Key)
		}
		b.NetworkUUID = networkUUID
	}
	notIn := strings.Join(keys, ",")
	if err := tx.Where("network_uuid = ?", networkUUID).Where("key not in (?)", notIn).
		Delete(&model.NetworkMetaTag{}).Error; err != nil {
		return nil, err
	}
	if len(body) > 0 {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(body).Error; err != nil {
			return nil, err
		}
	}
	if err := tx.Where("network_uuid = ?", networkUUID).Find(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) CreateNetworkMetaTags(networkUUID string, body []*model.NetworkMetaTag) ([]*model.NetworkMetaTag, error) {
	tx := d.DB.Begin()
	mt, err := CreateNetworkMetaTagsTransaction(tx, networkUUID, body)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return mt, nil
}

func (d *GormDatabase) GetNetworkMetaTags() ([]*model.NetworkMetaTag, error) {
	var networkMetaTagsModel []*model.NetworkMetaTag
	if err := d.DB.Find(&networkMetaTagsModel).Error; err != nil {
		return nil, err
	}
	return networkMetaTagsModel, nil
}
