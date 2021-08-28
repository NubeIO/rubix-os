package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type ProducerList struct {
	*model.SubscriberList
}

// GetProducerLists get all of them
func (d *GormDatabase) GetProducerLists() ([]*model.SubscriberList, error) {
	var producersModel []*model.SubscriberList

	query := d.DB.Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateProducerList make it
func (d *GormDatabase) CreateProducerList(body *model.SubscriberList) (*model.SubscriberList, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.SubscriptionList)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetProducerList get it
func (d *GormDatabase) GetProducerList(uuid string) (*model.SubscriberList, error) {
	var producerModel *model.SubscriberList
	query := d.DB.Where("uuid = ? ", uuid).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// GetProducerListBySubUUID get it by its
func (d *GormDatabase) GetProducerListBySubUUID(subscriptionUUID string) (*model.SubscriberList, error) {
	var producerModel *model.SubscriberList
	query := d.DB.Where("subscription_uuid = ? ", subscriptionUUID).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// DeleteProducerList deletes it
func (d *GormDatabase) DeleteProducerList(uuid string) (bool, error) {
	var producerModel *model.SubscriberList
	query := d.DB.Where("uuid = ? ", uuid).Delete(&producerModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}

// UpdateProducerList  update it
func (d *GormDatabase) UpdateProducerList(uuid string, body *model.SubscriberList) (*model.SubscriberList, error) {
	var producerModel *model.SubscriberList
	query := d.DB.Where("uuid = ?", uuid).Find(&producerModel);if query.Error != nil {
		return nil, query.Error
	}
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Subscription)
	query = d.DB.Model(&producerModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil

}

// DropProducerList delete all.
func (d *GormDatabase) DropProducerList() (bool, error) {
	var producerModel *model.SubscriberList
	query := d.DB.Where("1 = 1").Delete(&producerModel)
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
