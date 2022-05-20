package database

import (
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type Job struct {
	*model.Job
}

func (d *GormDatabase) GetJobs() ([]*model.Job, error) {
	var jobsModel []*model.Job
	query := d.DB.Find(&jobsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobsModel, nil
}

func (d *GormDatabase) CreateJob(body *model.Job) (*model.Job, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Job)
	body.Name = nameIsNil(body.Name)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) GetJob(uuid string) (*model.Job, error) {
	var jobModel *model.Job
	query := d.DB.Where("uuid = ? ", uuid).First(&jobModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil
}

func (d *GormDatabase) GetJobsByPluginConfigId(pcId string) ([]*model.Job, error) {
	var jobsModel []*model.Job
	query := d.DB.Where("plugin_conf_id = ?", pcId).Find(&jobsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobsModel, nil
}

func (d *GormDatabase) GetJobByPluginConfId(pcId string) (*model.Job, error) {
	var jobModel *model.Job
	query := d.DB.Where("plugin_conf_id = ?", pcId).First(&jobModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil
}

// DeleteJob delete a job
func (d *GormDatabase) DeleteJob(uuid string) (bool, error) {
	var jobModel *model.Job
	query := d.DB.Where("uuid = ? ", uuid).Delete(&jobModel)
	return d.deleteResponseBuilder(query)
}

// UpdateJob  returns the device for the given id or nil.
func (d *GormDatabase) UpdateJob(uuid string, body *model.Job) (*model.Job, error) {
	var jobModel *model.Job
	query := d.DB.Where("uuid = ?", uuid).First(&jobModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&jobModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil

}
