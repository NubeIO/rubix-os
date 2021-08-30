package database

import (
	"fmt"
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

// pass in subscriptionList UUID as what starts everything
// get subscription uuid
// get producer uuid
// get flow network uuid
// get subscription uuid
// update writeValue to subscriptionList value
// update writeValue to producerList value
// producer to decide to aspect the new write value
// producer to send a COV event to the subscription if the value was updated.
// add histories if enabled



//point
// - presentValue

// subscription needs readonly of producer presentValue
// - presentValue

// subscriptionList needs to write to the producerList
// - writeValue

//producer
// - presentValue
// - SLWriteUUID //subscriptionList UUID

//producerList (subscriptionUUID)
// - writeValue

//producerHist
// - presentValue
// - SLWriteUUID //subscriptionList UUID

// SubscriptionAction get its value
func (d *GormDatabase) SubscriptionAction(uuid string, body *model.SubscriptionList, write bool) (*model.Producer, error) {
	var slm *model.SubscriptionList
	subscriptionList := d.DB.Where("uuid = ? ", uuid).First(&slm); if subscriptionList.Error != nil {
		return nil, subscriptionList.Error
	}
	if slm == nil {
		return nil, nil
	}
	var sm *model.Subscription
	subscription := d.DB.Where("uuid = ? ", slm.SubscriptionUUID).First(&sm); if subscription.Error != nil {
		return nil, subscription.Error
	}
	subType := sm.SubscriptionType
	subscriptionUUID := sm.UUID
	streamUUID := sm.StreamUUID
	producerUUID := sm.ProducerUUID
	writeV := body.WriteValue
 
	var s *model.Stream
	stream := d.DB.Where("uuid = ? ", streamUUID).First(&s); if subscription.Error != nil {
		return nil, stream.Error
	}
	streamListUUID := s.StreamListUUID
	var fn *model.FlowNetwork
	flow := d.DB.Where("stream_list_uuid = ? ", streamListUUID).First(&fn); if subscription.Error != nil {
		return nil, flow.Error
	}
	flowUUID := fn.UUID
	isRemote := fn.IsRemote
	fmt.Println("subType", subType, "subscriptionUUID", subscriptionUUID, "streamUUID", streamUUID, "producerUUID", producerUUID,"flowUUID", flowUUID, "isRemote", isRemote, writeV, write)
	if !isRemote { // local
		pm := new(model.Producer)
		query := d.DB.Where("uuid = ?", producerUUID).Find(&pm);if query.Error != nil {
			return nil, query.Error
		}
		if query == nil {
			return nil, nil
		}
		if write { //write new value to producer
			pm.PresentValue = body.WriteValue
			query = d.DB.Model(&pm).Updates(pm);if query.Error != nil {
				return nil, query.Error
			}
			ph := new(model.ProducerHistory)
			ph.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
			ph.ProducerUUID = producerUUID
			ph.PresentValue = pm.PresentValue
			query = d.DB.Model(&ph).Updates(ph);if query.Error != nil {
				return nil, query.Error
			}
			return pm, nil
		} else {
			d.DB.Where("uuid = ? ", uuid).First(&pm); if query.Error != nil {
				return nil, query.Error
			}
			return pm, nil
		}
	} else {
		pm := new(model.Producer)
		pm.UUID = producerUUID
		pm.PresentValue = body.WriteValue
		point, err := eventbus.EventREST(fn, pm, write)
		if err != nil {
			return nil, err
		}
		return point, err
	}
}

// pass in subscriptionList UUID as what starts everything
// get subscription uuid
// get producer uuid
// get flow network uuid
// get subscription uuid
// get update writeValue to subscription
// send event to point to update writeValue and update the point
// update the priory array
// update the producer based of what the point returns
// producer to send a COV event to the subscription.
// subscription to update the presentValue and the history if enabled

// SubscriptionActionPoint get its value or write
func (d *GormDatabase) SubscriptionActionPoint(slUUID string, pointBody *model.Point, write bool) (*model.Producer, error) {
	var slm *model.SubscriptionList
	sl := d.DB.Where("uuid = ? ", slUUID).First(&slm); if sl.Error != nil {
		return nil, sl.Error
	}
	if slm == nil {
		return nil, nil
	}

	var sm *model.Subscription
	subscription := d.DB.Where("uuid = ? ", slm.SubscriptionUUID).First(&sm); if subscription.Error != nil {
		return nil, subscription.Error
	}
	subType := sm.SubscriptionType
	subscriptionUUID := sm.UUID
	streamUUID := sm.StreamUUID
	producerUUID := sm.ProducerUUID
	writeV := pointBody.WriteValue
	pointUUID := sm.ProducerThingUUID

	var s *model.Stream
	stream := d.DB.Where("uuid = ? ", streamUUID).First(&s); if subscription.Error != nil {
		return nil, stream.Error
	}
	streamListUUID := s.StreamListUUID
	var fn *model.FlowNetwork
	flow := d.DB.Where("stream_list_uuid = ? ", streamListUUID).First(&fn); if subscription.Error != nil {
		return nil, flow.Error
	}
	flowUUID := fn.UUID
	isRemote := fn.IsRemote
	fmt.Println("subType", subType, "subscriptionUUID", subscriptionUUID, "streamUUID", streamUUID, "producerUUID", producerUUID,"flowUUID", flowUUID, "isRemote", isRemote, writeV, write)
	if !isRemote { // local
		pm := new(model.Producer)
		query := d.DB.Where("uuid = ?", producerUUID).Find(&pm);if query.Error != nil {
			return nil, query.Error
		}
		if query == nil {
			return nil, nil
		}
		if write { //write new value to producer
			//pm.WriteValue = body.WriteValue
			//pm.PresentValue = body.WriteValue
			//query = d.DB.Model(&pm).Updates(pm);if query.Error != nil {
			//	return nil, query.Error
			//}
			//ph := new(model.ProducerHistory)
			//ph.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
			//ph.ProducerUUID = producerUUID
			//ph.PresentValue = pm.PresentValue
			//query = d.DB.Model(&ph).Updates(ph);if query.Error != nil {
			//	return nil, query.Error
			//}
			return pm, nil
		} else {
			//d.DB.Where("uuid = ? ", uuid).First(&pm); if query.Error != nil {
			//	return nil, query.Error
			//}
			return pm, nil
		}
	} else {
		if write {
			point, err := eventbus.EventRESTPoint(pointUUID, fn, pointBody, write)
			if err != nil {
				return nil, err
			}
			pm := new(model.Producer)
			pm.UUID = producerUUID
			pm.PresentValue = point.WriteValue
			query := d.DB.Model(&pm).Updates(pm);if query.Error != nil {
				return nil, query.Error
			}
			ph := new(model.ProducerHistory)
			ph.UUID = utils.MakeTopicUUID(model.CommonNaming.ProducerHistory)
			ph.ProducerUUID = producerUUID
			ph.PresentValue = pm.PresentValue
			query = d.DB.Model(&ph).Updates(ph);if query.Error != nil {
				return nil, query.Error
			}
			return pm, err
		} else {
			point, err := eventbus.EventRESTPoint(pointUUID, fn, pointBody, write)
			if err != nil {
				return nil, err
			}
			pm := new(model.Producer)
			pm.UUID = producerUUID
			pm.PresentValue = point.WriteValue


			return pm, err
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
