package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm/clause"
	"strings"
)

func (d *GormDatabase) UpdateHostTags(hostUUID string, body []*model.HostTag) ([]*model.HostTag, error) {
	tx := d.DB.Begin()
	var tags []string
	for i, b := range body {
		body[i].HostUUID = hostUUID
		tags = append(tags, b.Tag)
	}
	notIn := strings.Join(tags, ",")
	if err := tx.Where("host_uuid = ?", hostUUID).Where("tag not in (?)", notIn).Delete(&model.HostTag{}).
		Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(body) > 0 {
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(body).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Where("host_uuid = ?", hostUUID).Find(&body).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return body, nil
}
