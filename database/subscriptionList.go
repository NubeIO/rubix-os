package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type SubscriptionList struct {
	*model.SubscriptionList
}

// GetSubscriptionLists get all of them
func (d *GormDatabase) GetSubscriptionLists() ([]*model.SubscriptionList, error) {
	var subscriptionsModel []*model.SubscriptionList
	query := d.DB.Find(&subscriptionsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionsModel, nil
}

// CreateSubscriptionList make it
func (d *GormDatabase) CreateSubscriptionList(body *model.SubscriptionList) (*model.SubscriptionList, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.SubscriptionList)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetSubscriptionList get it
func (d *GormDatabase) GetSubscriptionList(uuid string) (*model.SubscriptionList, error) {
	var subscriptionModel *model.SubscriptionList
	query := d.DB.Where("uuid = ? ", uuid).First(&subscriptionModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionModel, nil
}

// GetSubscriptionListByThing get it by its thing uuid
func (d *GormDatabase) GetSubscriptionListByThing(toThingUUID string) (*model.SubscriptionList, error) {
	var subscriptionModel *model.SubscriptionList
	query := d.DB.Where("to_thing_uuid = ? ", toThingUUID).First(&subscriptionModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionModel, nil
}


// DeleteSubscriptionList deletes it
func (d *GormDatabase) DeleteSubscriptionList(uuid string) (bool, error) {
	var subscriptionModel *model.SubscriptionList
	query := d.DB.Where("uuid = ? ", uuid).Delete(&subscriptionModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// UpdateSubscriptionList  update it
func (d *GormDatabase) UpdateSubscriptionList(uuid string, body *model.SubscriptionList) (*model.SubscriptionList, error) {
	var subscriptionModel *model.SubscriptionList
	query := d.DB.Where("uuid = ?", uuid).Find(&subscriptionModel);if query.Error != nil {
		return nil, query.Error
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Subscription)
	query = d.DB.Model(&subscriptionModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionModel, nil

}

// DropSubscriptionsList delete all.
func (d *GormDatabase) DropSubscriptionsList() (bool, error) {
	var subscriptionModel *model.SubscriptionList
	query := d.DB.Where("1 = 1").Delete(&subscriptionModel)
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
