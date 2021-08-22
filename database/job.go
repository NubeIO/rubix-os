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
	query := d.DB.Preload(gatewaySubscriberChildTable).Find(&jobsModel)
	if query.Error != nil {
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
	query := d.DB.Where("uuid = ? ", uuid).First(&jobModel); if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil
}

// DeleteJob delete a job
func (d *GormDatabase) DeleteJob(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ? ", uuid).Delete(&jobModel);if query.Error != nil {
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
	query := d.DB.Where("uuid = ?", uuid).Find(&jobModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&jobModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil

}
//
//func (d *GormDatabase) CreateJobSubscriber(body *model.JobSubscriber, jobUUID string)  error {
//	query := d.DB.Where("uuid = ?", jobUUID).Find(&jobModel); if query.Error != nil {
//		return query.Error
//	}
//	body.UUID, _ = utils.MakeUUID()
//	body.JobUUID = jobUUID
//	fmt.Println(body.UUID,body.JobUUID)
//	n := d.DB.Create(body).Error
//	return n
//
//}
//
//
//func (d *GormDatabase) GetJobSubscribers() ([]model.JobSubscriber, error) {
//	query := d.DB.Find(&jobSubscribersModel)
//	if query.Error != nil {
//		return nil, query.Error
//	}
//	return jobSubscribersModel, nil
//}
//
//
//// DeleteJobSubscriber delete a job subscriber(
//func (d *GormDatabase) DeleteJobSubscriber(uuid string) (bool, error) {
//	query := d.DB.Where("uuid = ? ", uuid).Delete(&jobSubscriberModel);if query.Error != nil {
//		return false, query.Error
//	}
//	r := query.RowsAffected
//	if r == 0 {
//		return false, nil
//	} else {
//		return true, nil
//	}
//
//}
//
//
//// UpdateJobSubscriber  returns the device for the given id or nil.
//func (d *GormDatabase) UpdateJobSubscriber(uuid string, body *model.JobSubscriber) (*model.JobSubscriber, error) {
//	query := d.DB.Where("uuid = ?", uuid).Find(&jobSubscriberModel)
//	if query.Error != nil {
//		return nil, query.Error
//	}
//	query = d.DB.Model(&jobSubscriberModel).Updates(body)
//	if query.Error != nil {
//		return nil, query.Error
//	}
//	return jobSubscriberModel, nil
//}