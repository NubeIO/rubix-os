package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

// GetSchedules returns all things.
func (d *GormDatabase) GetSchedules() ([]*model.Schedule, error) {
	var schModel []*model.Schedule
	query := d.DB.Find(&schModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return schModel, nil
}

// GetSchedule returns the thing for the given id or nil.
func (d *GormDatabase) GetSchedule(uuid string) (*model.Schedule, error) {
	var schModel *model.Schedule
	query := d.DB.Where("uuid = ? ", uuid).First(&schModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return schModel, nil
}

// GetScheduleByField returns the sch for the given field ie name or nil.
func (d *GormDatabase) GetScheduleByField(field string, value string) (*model.Schedule, error) {
	var schModel *model.Schedule
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).First(&schModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return schModel, nil
}

// CreateSchedule creates a thing.
func (d *GormDatabase) CreateSchedule(body *model.Schedule) (*model.Schedule, error) {
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Schedule)
	body.Name = nameIsNil(body.Name)
	body.ThingClass = model.ThingClass.Schedule
	body.ThingType = model.ThingClass.Schedule
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

// UpdateSchedule  update it
func (d *GormDatabase) UpdateSchedule(uuid string, body *model.Schedule) (*model.Schedule, error) {
	var schModel *model.Schedule
	query := d.DB.Where("uuid = ?", uuid).Find(&schModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&schModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return schModel, nil

}

// DeleteSchedule delete a thing.
func (d *GormDatabase) DeleteSchedule(uuid string) (bool, error) {
	var schModel *model.Schedule
	query := d.DB.Where("uuid = ? ", uuid).Delete(&schModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// DropSchedules delete all things.
func (d *GormDatabase) DropSchedules() (bool, error) {
	var schModel *model.Schedule
	query := d.DB.Where("1 = 1").Delete(&schModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
