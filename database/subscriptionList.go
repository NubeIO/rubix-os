package database

import (
	"github.com/NubeDev/flow-framework/eventbus"
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
func (d *GormDatabase) GetSubscriptionListByThing(producerThingUUID string) (*model.SubscriptionList, error) {
	var subscriptionModel *model.SubscriptionList
	query := d.DB.Where("producer_thing_uuid = ? ", producerThingUUID).First(&subscriptionModel); if query.Error != nil {
		return nil, query.Error
	}
	return subscriptionModel, nil
}



// SubscriptionAction get its value
func (d *GormDatabase) SubscriptionAction(uuid string, body interface{}, askRefresh bool, askResponse bool, write bool, thingType string, flowNetworkUUID string) (interface{}, error) { //TODO add in more logic
	var subscriptionListModel *model.SubscriptionList
	subscriptionList := d.DB.Where("uuid = ? ", uuid).First(&subscriptionListModel); if subscriptionList.Error != nil {
		return nil, subscriptionList.Error
	}
	if subscriptionListModel == nil {
		return nil, nil
	}
	var subscriptionModel *model.Subscription
	subscription := d.DB.Where("uuid = ? ", subscriptionListModel.SubscriptionUUID).First(&subscriptionModel); if subscription.Error != nil {
		return nil, subscription.Error
	}
	subType := subscriptionModel.SubscriptionType

	if flowNetworkUUID != "" { //remote point
		var flowModel *model.FlowNetwork
		 d.DB.Where("uuid = ? ", flowNetworkUUID).First(&flowModel); if subscription.Error != nil {
			return nil, subscription.Error
		}
		if subType == model.CommonNaming.Point {
			ip := flowModel.FlowIP
			port := flowModel.FlowPort
			token := flowModel.FlowToken
			pntUUID := subscriptionModel.ProducerThingUUID

			if write { //write
				point, err := eventbus.EventREST(pntUUID, body, ip, port, token, write, thingType)
				if err != nil {
					return nil, err
				}
				return point, err
			} else { // read
				point, err := eventbus.EventREST(pntUUID, body, ip, port, token, write, thingType)
				if err != nil {
					return nil, err
				}
				return point, err
			}
		} else {
			return nil, nil
		}

	} else { // local point
		if subType == model.CommonNaming.Point {
			pnt := new(model.Point)
			d.DB.Where("uuid = ? ", subscriptionModel.ProducerThingUUID).First(&pnt); if subscription.Error != nil {
				return nil, subscription.Error
			}
			if write {
				var pointModel *model.Point
				query := d.DB.Where("uuid = ?", subscriptionModel.ProducerThingUUID).Find(&pointModel);if query.Error != nil {
					return nil, query.Error
				}
				query = d.DB.Model(&pointModel).Updates(body);if query.Error != nil {
					return nil, query.Error
				}
				return pointModel, nil
			}
			return pnt, nil
		} else {
			return nil, nil
		}
	}

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

func readValue(){

}

func writeValue(){

}