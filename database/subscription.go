package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type Subscriptions struct {
	*model.Subscription
}

// GetSubscriptions get all of them
func (d *GormDatabase) GetSubscriptions() ([]*model.Subscription, error) {
	var subscriptionsModel []*model.Subscription
	query := d.DB.Preload("SubscriptionList").Find(&subscriptionsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionsModel, nil
}

// CreateSubscription make it
func (d *GormDatabase) CreateSubscription(body *model.Subscription) (*model.Subscription, error) {
	_, err := d.GetStreamGateway(body.StreamUUID);if err != nil {
		return nil, errorMsg("GetStreamGateway", "error on trying to get validate the stream UUID", nil)
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Subscription)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetSubscription get it
func (d *GormDatabase) GetSubscription(uuid string) (*model.Subscription, error) {
	var subscriptionModel *model.Subscription
	query := d.DB.Where("uuid = ? ", uuid).First(&subscriptionModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionModel, nil
}

// DeleteSubscription deletes it
func (d *GormDatabase) DeleteSubscription(uuid string) (bool, error) {
	var subscriptionModel *model.Subscription
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

// UpdateSubscription  update it
func (d *GormDatabase) UpdateSubscription(uuid string, body *model.Subscription) (*model.Subscription, error) {
	var subscriptionModel *model.Subscription
	query := d.DB.Where("uuid = ?", uuid).Find(&subscriptionModel);if query.Error != nil {
		return nil, query.Error
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Subscription)
	query = d.DB.Model(&subscriptionModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionModel, nil

}

// DropSubscriptions delete all.
func (d *GormDatabase) DropSubscriptions() (bool, error) {
	var subscriptionModel *model.Subscription
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
