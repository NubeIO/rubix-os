package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type ProducerList struct {
	*model.ProducerList
}

// GetProducerLists get all of them
func (d *GormDatabase) GetProducerLists() ([]*model.ProducerList, error) {
	var producersModel []*model.ProducerList

	query := d.DB.Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateProducerList make it
func (d *GormDatabase) CreateProducerList(body *model.ProducerList) (*model.ProducerList, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.SubscriptionList)
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetProducerList get it
func (d *GormDatabase) GetProducerList(uuid string) (*model.ProducerList, error) {
	var producerModel *model.ProducerList
	query := d.DB.Where("uuid = ? ", uuid).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}

// GetProducerListByThing get it by its
func (d *GormDatabase) GetProducerListByThing(fromThingUUID string) (*model.ProducerList, error) {
	var producerModel *model.ProducerList
	query := d.DB.Where("from_thing_uuid = ? ", fromThingUUID).First(&producerModel); if query.Error != nil {
		return nil, query.Error
	}
	return producerModel, nil
}


// DeleteProducerList deletes it
func (d *GormDatabase) DeleteProducerList(uuid string) (bool, error) {
	var producerModel *model.ProducerList
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
func (d *GormDatabase) UpdateProducerList(uuid string, body *model.ProducerList) (*model.ProducerList, error) {
	var producerModel *model.ProducerList
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
	var producerModel *model.ProducerList
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
