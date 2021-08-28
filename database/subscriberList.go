package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type SubscriberList struct {
	*model.SubscriberList
}

// GetSubscriberLists get all of them
func (d *GormDatabase) GetSubscriberLists() ([]*model.SubscriberList, error) {
	var subscribersModel []*model.SubscriberList

	query := d.DB.Find(&subscribersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return subscribersModel, nil
}

// CreateSubscriberList make it
func (d *GormDatabase) CreateSubscriberList(body *model.SubscriberList) (*model.SubscriberList, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.SubscriptionList)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetSubscriberList get it
func (d *GormDatabase) GetSubscriberList(uuid string) (*model.SubscriberList, error) {
	var subscriberModel *model.SubscriberList
	query := d.DB.Where("uuid = ? ", uuid).First(&subscriberModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriberModel, nil
}

// GetSubscriberListByThing get it by its
func (d *GormDatabase) GetSubscriberListByThing(fromThingUUID string) (*model.SubscriberList, error) {
	var subscriberModel *model.SubscriberList
	query := d.DB.Where("from_thing_uuid = ? ", fromThingUUID).First(&subscriberModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriberModel, nil
}


// DeleteSubscriberList deletes it
func (d *GormDatabase) DeleteSubscriberList(uuid string) (bool, error) {
	var subscriberModel *model.SubscriberList
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

// UpdateSubscriberList  update it
func (d *GormDatabase) UpdateSubscriberList(uuid string, body *model.SubscriberList) (*model.SubscriberList, error) {
	var subscriberModel *model.SubscriberList
	query := d.DB.Where("uuid = ?", uuid).Find(&subscriberModel);if query.Error != nil {
		return nil, query.Error
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Subscription)
	query = d.DB.Model(&subscriberModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return subscriberModel, nil

}

// DropSubscriberList delete all.
func (d *GormDatabase) DropSubscriberList() (bool, error) {
	var subscriberModel *model.SubscriberList
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
