package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

type Subscriber struct {
	*model.Subscriber

}

// GetSubscribers get all of them
func (d *GormDatabase) GetSubscribers() ([]*model.Subscriber, error) {
	var subscribersModel []*model.Subscriber
	query := d.DB.Find(&subscribersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return subscribersModel, nil
}


// CreateSubscriber make it
func (d *GormDatabase) CreateSubscriber(body *model.Subscriber) (*model.Subscriber, error) {
	//call points and make it exists
	_, err := d.GetStreamGateway(body.StreamUUID);if err != nil {
		return nil, errorMsg("GetStreamGateway", "error on trying to get validate the gateway UUID", nil)
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Subscriber)
	err = d.DB.Create(&body).Error; if err != nil {
		return nil, errorMsg("CreateSubscriber", "error on trying to add a new Subscriber", nil)
	}
	return body, nil
}



// GetSubscriber get it
func (d *GormDatabase) GetSubscriber(uuid string) (*model.Subscriber, error) {
	var subscriberModel *model.Subscriber
	query := d.DB.Where("uuid = ? ", uuid).First(&subscriberModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriberModel, nil
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

// DropSubscribers delete all.
func (d *GormDatabase) DropSubscribers() (bool, error) {
	var subscriberModel *model.Subscriber
	query := d.DB.Where("1 = 1").Delete(&subscriberModel)
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