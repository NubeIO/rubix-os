package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
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
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Job)
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

// DeleteJob delete a job
func (d *GormDatabase) DeleteJob(uuid string) (bool, error) {
	var jobModel *model.Job
	query := d.DB.Where("uuid = ? ", uuid).Delete(&jobModel)
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

// UpdateJob  returns the device for the given id or nil.
func (d *GormDatabase) UpdateJob(uuid string, body *model.Job) (*model.Job, error) {
	var jobModel *model.Job
	query := d.DB.Where("uuid = ?", uuid).Find(&jobModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&jobModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil

}

// DropJobs delete all.
func (d *GormDatabase) DropJobs() (bool, error) {
	var jobModel *model.Job
	query := d.DB.Where("1 = 1").Delete(&jobModel)
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

//
//func (d *GormDatabase) CreateJobProducer(body *model.JobProducer, jobUUID string)  error {
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
//func (d *GormDatabase) GetJobProducers() ([]model.JobProducer, error) {
//	query := d.DB.Find(&jobProducersModel)
//	if query.Error != nil {
//		return nil, query.Error
//	}
//	return jobProducersModel, nil
//}
//
//
//// DeleteJobProducer delete a job producer(
//func (d *GormDatabase) DeleteJobProducer(uuid string) (bool, error) {
//	query := d.DB.Where("uuid = ? ", uuid).Delete(&jobProducerModel);if query.Error != nil {
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
//// UpdateJobProducer  returns the device for the given id or nil.
//func (d *GormDatabase) UpdateJobProducer(uuid string, body *model.JobProducer) (*model.JobProducer, error) {
//	query := d.DB.Where("uuid = ?", uuid).Find(&jobProducerModel)
//	if query.Error != nil {
//		return nil, query.Error
//	}
//	query = d.DB.Model(&jobProducerModel).Updates(body)
//	if query.Error != nil {
//		return nil, query.Error
//	}
//	return jobProducerModel, nil
//}
