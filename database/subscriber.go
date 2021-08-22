package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

var st = model.NewSubscriberTypeEnum()
var sa = model.NewSubscriberApplicationEnum()

type Subscriber struct {
	*model.Subscriber

}



var subscriberPointsChildTable = "Point"
var subscriberJobsChildTable = "Job"


// GetSubscribers get all of them
func (d *GormDatabase) GetSubscribers() ([]model.Subscriber, error) {
	var subscribersModel []model.Subscriber
	query := d.DB.Find(&subscribersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return subscribersModel, nil
}

// CreateSubscriber make it
func (d *GormDatabase) CreateSubscriber(body *model.Subscriber)  error {
	body.UUID, _ = utils.MakeTopicUUID(model.CommonNaming.Subscriber)
	if body.SubscriberType == st.Point {
		//call points and make it exists
		query, err := d.GetPoint(body.PointUUID, false)
		if err != nil {
			return errorMsg("CreateSubscriber", "error on trying to add", nil)
		}
		if query != nil {
			body.SubscriberApplication =sa.Mapping
			n := d.DB.Create(body).Error
			return n
		}
	}
	return nil
}

// GetSubscriber get it
func (d *GormDatabase) GetSubscriber(uuid string) (*model.Subscriber, error) {
	var subscriberModel *model.Subscriber
	query := d.DB.Where("uuid = ? ", uuid).First(&subscriberModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriberModel, nil
}

// DeleteSubscriber deletes it
func (d *GormDatabase) DeleteSubscriber(uuid string) (bool, error) {
	var subscriberModel *model.Subscriber
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
	var subscriberModel *model.Subscriber
	query := d.DB.Where("uuid = ?", uuid).Find(&subscriberModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&subscriberModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return subscriberModel, nil

}
