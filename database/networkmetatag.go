package database

import (
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm/clause"
	"strings"
)

func (d *GormDatabase) CreateNetworkMetaTags(networkUUID string, body []*model.NetworkMetaTag) ([]*model.NetworkMetaTag,
	error) {
	tx := d.DB.Begin()
	var networkUUIDs []string
	for _, b := range body {
		var count int64
		tx.Model(&model.NetworkMetaTag{}).Where("network_uuid = ?", networkUUID).Where("uuid = ?",
			b.UUID).Count(&count)
		if count == 0 {
			b.UUID = nuuid.MakeTopicUUID(model.CommonNaming.NetworkMetaTag)
		} else {
			networkUUIDs = append(networkUUIDs, b.UUID)
		}
		b.NetworkUUID = networkUUID
	}
	notIn := strings.Join(networkUUIDs, ",")
	if err := tx.Where("network_uuid = ?", networkUUID).Where("uuid not in (?)", notIn).
		Delete(&model.NetworkMetaTag{}).Error; err != nil {
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
