package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type Subscriptions struct {
	*model.Subscriptions
}

// GetSubscriptions get all of them
func (d *GormDatabase) GetSubscriptions() ([]model.Subscriptions, error) {
	var subscriptionsModel []model.Subscriptions
	query := d.DB.Preload(subscriberPointsChildTable).Find(&subscriptionsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionsModel, nil
}

// CreateSubscription make it
func (d *GormDatabase) CreateSubscription(body *model.Subscriptions)  error {
	body.UUID, _ = utils.MakeTopicUUID(model.CommonNaming.Subscription)
	n := d.DB.Create(body).Error
	return n
}

// GetSubscription get it
func (d *GormDatabase) GetSubscription(uuid string) (*model.Subscriptions, error) {
	var subscriptionModel *model.Subscriptions
	query := d.DB.Where("uuid = ? ", uuid).First(&subscriptionModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionModel, nil
}

// DeleteSubscription deletes it
func (d *GormDatabase) DeleteSubscription(uuid string) (bool, error) {
	var subscriptionModel *model.Subscriptions
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
func (d *GormDatabase) UpdateSubscription(uuid string, body *model.Subscriptions) (*model.Subscriptions, error) {
	var subscriptionModel *model.Subscriptions
	query := d.DB.Where("uuid = ?", uuid).Find(&subscriptionModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&subscriptionModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionModel, nil

}
