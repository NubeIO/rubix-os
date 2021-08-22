package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type Subscriber struct {
	*model.Subscriber
}

var subscribersModel []model.Subscriber
var subscriberModel *model.Subscriber
//var jobSubscribersModel []model.JobSubscriber
//var jobSubscriberModel *model.JobSubscriber


// GetSubscribers get all of them
func (d *GormDatabase) GetSubscribers() ([]model.Subscriber, error) {
	query := d.DB.Preload(subscriberChildTable).Find(&jobsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return subscribersModel, nil
}

// CreateSubscriber make it
func (d *GormDatabase) CreateSubscriber(body *model.Subscriber)  error {
	body.UUID, _ = utils.MakeUUID()
	n := d.DB.Create(body).Error
	return n
}

// GetSubscriber get it
func (d *GormDatabase) GetSubscriber(uuid string) (*model.Subscriber, error) {
	query := d.DB.Where("uuid = ? ", uuid).First(&subscriberModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriberModel, nil
}

// DeleteSubscriber deletes it
func (d *GormDatabase) DeleteSubscriber(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ? ", uuid).Delete(&subscriberModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// UpdateSubscriber  update it
func (d *GormDatabase) UpdateSubscriber(uuid string, body *model.Subscriber) (*model.Subscriber, error) {
	query := d.DB.Where("uuid = ?", uuid).Find(&subscriberModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&subscriberModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return subscriberModel, nil

}

//func (d *GormDatabase) CreateJobSubscriber(body *model.Subscriber, jobUUID string)  error {
//	query := d.DB.Where("uuid = ?", jobUUID).Find(&subscriberModel); if query.Error != nil {
//		return query.Error
//	}
//	body.UUID, _ = utils.MakeUUID()
//	body.UUID = jobUUID
//	fmt.Println(body.UUID,body.UUID)
//	n := d.DB.Create(body).Error
//	return n
//
//}

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