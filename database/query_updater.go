package database

func (d *GormDatabase) updateTags(model, tags interface{}) error {
	return d.DB.Model(model).Association("Tags").Replace(tags)
}
