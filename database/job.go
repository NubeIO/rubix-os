package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type Job struct {
	*model.Job
}

var jobsModel []model.Job
var jobModel *model.Job


func (d *GormDatabase) GetJobs() ([]model.Job, error) {
	query := d.DB.Find(&jobsModel);if query.Error != nil {
		return nil, query.Error
	}
	return jobsModel, nil
}

func (d *GormDatabase) CreateJob(body *model.Job)  error {
	body.UUID, _ = utils.MakeUUID()
	n := d.DB.Create(body).Error
	return n
}

func (d *GormDatabase) GetJob(uuid string) (*model.Job, error) {
	query := d.DB.Where("id = ? ", uuid).First(&jobModel); if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil
}

// DeleteJob delete a job
func (d *GormDatabase) DeleteJob(uuid string) (bool, error) {
	query := d.DB.Where("id = ? ", uuid).Delete(&jobModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// UpdateJob  returns the device for the given id or nil.
func (d *GormDatabase) UpdateJob(uuid string, body *model.Job) (*model.Job, error) {
	query := d.DB.Where("id = ?", uuid).Find(&jobModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&jobModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil

}
