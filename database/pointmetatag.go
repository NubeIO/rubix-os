package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm/clause"
	"strings"
)

func (d *GormDatabase) CreatePointMetaTags(pointUUID string, body []*model.PointMetaTag) ([]*model.PointMetaTag,
	error) {
	tx := d.DB.Begin()
	var keys []string
	for _, b := range body {
		var count int64
		tx.Model(&model.PointMetaTag{}).Where("point_uuid = ?", pointUUID).Where("key = ?", b.Key).
			Count(&count)
		if count > 0 {
			keys = append(keys, b.Key)
		}
		b.PointUUID = pointUUID
	}
	notIn := strings.Join(keys, ",")
	if err := tx.Where("point_uuid = ?", pointUUID).Where("key not in (?)", notIn).
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
	if err := tx.Where("point_uuid = ?", pointUUID).Find(&body).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return body, nil
}
