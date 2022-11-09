package database

import (
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm/clause"
	"strings"
)

func (d *GormDatabase) CreatePointMetaTags(pointUUID string, body []*model.PointMetaTag) ([]*model.PointMetaTag,
	error) {
	tx := d.DB.Begin()
	var pointUUIDs []string
	for _, b := range body {
		var count int64
		tx.Model(&model.PointMetaTag{}).Where("point_uuid = ?", pointUUID).Where("uuid = ?", b.UUID).
			Count(&count)
		if count == 0 {
			b.UUID = nuuid.MakeTopicUUID(model.CommonNaming.PointMetaTag)
		} else {
			pointUUIDs = append(pointUUIDs, b.UUID)
		}
		b.PointUUID = pointUUID
	}
	notIn := strings.Join(pointUUIDs, ",")
	if err := tx.Where("point_uuid = ?", pointUUID).Where("uuid not in (?)", notIn).
		Delete(&model.PointMetaTag{}).Error; err != nil {
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
