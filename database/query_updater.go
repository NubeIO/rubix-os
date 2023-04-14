package database

import "gorm.io/gorm"

func updateTagsTransaction(db *gorm.DB, model, tags interface{}) error {
	return db.Model(model).Association("Tags").Replace(tags)
}

func (d *GormDatabase) updateTags(model, tags interface{}) error {
	return updateTagsTransaction(d.DB, model, tags)
}
