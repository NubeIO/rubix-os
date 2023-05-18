package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm"
)

func updateTagsTransaction(db *gorm.DB, model interface{}, tags []*model.Tag) error {
	return db.Model(model).Association("Tags").Replace(tags)
}

func (d *GormDatabase) updateTags(model interface{}, tags []*model.Tag) error {
	return updateTagsTransaction(d.DB, model, tags)
}

func (d *GormDatabase) updateGroups(model interface{}, groups []*model.Group) error {
	return d.DB.Model(model).Association("Groups").Replace(groups)
}
