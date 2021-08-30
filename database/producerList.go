package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type ProducerList struct {
	*model.ProducerSubscriptionList
}

// GetProducerLists get all of them
func (d *GormDatabase) GetProducerLists() ([]*model.ProducerSubscriptionList, error) {
	var producersModel []*model.ProducerSubscriptionList

	query := d.DB.Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateProducerList make it
func (d *GormDatabase) CreateProducerList(body *model.ProducerSubscriptionList) (*model.ProducerSubscriptionList, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.SubscriptionList)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetProducerList get it
func (d *GormDatabase) GetProducerList(uuid string) (*model.ProducerSubscriptionList, error) {
	var producerModel *model.ProducerSubscriptionList
	query := d.DB.Where("uuid = ? ", uuid).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// GetProducerListBySubUUID get it by its
func (d *GormDatabase) GetProducerListBySubUUID(subscriptionUUID string) (*model.ProducerSubscriptionList, error) {
	var producerModel *model.ProducerSubscriptionList
	query := d.DB.Where("subscription_uuid = ? ", subscriptionUUID).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// DeleteProducerList deletes it
func (d *GormDatabase) DeleteProducerList(uuid string) (bool, error) {
	var producerModel *model.ProducerSubscriptionList
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
func (d *GormDatabase) UpdateProducerList(uuid string, body *model.ProducerSubscriptionList) (*model.ProducerSubscriptionList, error) {
	var producerModel *model.ProducerSubscriptionList
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
	var producerModel *model.ProducerSubscriptionList
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

/*
// get
 */

//// AddHistory  add a history record
//func (d *GormDatabase) AddHistory(uuid string, body *model.ProducerSubscriptionList) (*model.ProducerHistory, error) {
//	var producerModel *model.ProducerSubscriptionList
//	var producerHist *model.ProducerHistory
//	query := d.DB.Where("uuid = ?", uuid).Find(&producerModel);if query.Error != nil {
//		return nil, query.Error
//	}
//	body.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
//	query = d.DB.Model(&producerModel).Updates(body);if query.Error != nil {
//		return nil, query.Error
//	}
//	return producerModel, nil
//
//}