package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


// GetCommandGroups returns all objects.
func (d *GormDatabase) GetCommandGroups() ([]*model.CommandGroup, error) {
	var commandGroup []*model.CommandGroup
	query := d.DB.Find(&commandGroup)
	if query.Error != nil {
		return nil, query.Error
	}
	return commandGroup, nil
}



// GetCommandGroup returns object.
func (d *GormDatabase) GetCommandGroup(uuid string) (*model.CommandGroup, error) {
	var commandGroup *model.CommandGroup
	query := d.DB.Where("uuid = ? ", uuid).First(&commandGroup); if query.Error != nil {
		return nil, query.Error
	}
	return commandGroup, nil
}



// CreateCommandGroup creates a object.
func (d *GormDatabase) CreateCommandGroup(body *model.CommandGroup) (*model.CommandGroup, error) {
	body.UUID = utils.MakeTopicUUID("")
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}



// UpdateCommandGroup update it.
func (d *GormDatabase) UpdateCommandGroup(uuid string, body *model.CommandGroup) (*model.CommandGroup, error) {
	var commandGroup *model.CommandGroup
	query := d.DB.Where("uuid = ?", uuid).Find(&commandGroup);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&commandGroup).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return commandGroup, nil

}


// DeleteCommandGroup delete an object.
func (d *GormDatabase) DeleteCommandGroup(uuid string) (bool, error) {
	var commandGroup *model.CommandGroup
	query := d.DB.Where("uuid = ? ", uuid).Delete(&commandGroup);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}