package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

//var st = model.NewSubscriberTypeEnum()
//var sa = model.NewSubscriberApplicationEnum()

type Subscriber struct {
	*model.Subscriber

}



//var subscriberPointsChildTable = "Point"
//var subscriberJobsChildTable = "Job"


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
	if body.SubscriberType == model.CommonNaming.Point {
		//call points and make it exists
		_, err := d.GetGateway(body.GatewayUUID);if err != nil {
			return nil, errorMsg("GetGateway", "error on trying to get validate the gateway UUID", nil)
		}
		_, err = d.GetPoint(body.ToUUID, false);if err != nil {
			return nil, errorMsg("CreateSubscriber", "error on trying to get validate the point UUID", nil)
		}
		err = d.DB.Create(&body).Error; if err != nil {
			return nil, errorMsg("CreateSubscriber", "error on trying to add a new Subscriber", nil)
		}
		fmt.Println(9999999)
		 if !body.IsRemote {
			 fmt.Println(8888, body.IsRemote, 888888)
			 u := utils.MakeTopicUUID("")
			 d.DB.Create(&model.PointSubscriptionLedger{UUID: u, GatewayUUID: body.GatewayUUID, SubscriberUUID: body.UUID, PointUUID: body.FromUUID})
		 } else if body.IsRemote {
			 var sm model.Subscription
			 fmt.Println(5555555, body.IsRemote, 5555555)
			sm.UUID = utils.MakeTopicUUID(model.CommonNaming.Subscriber)
			sm.GatewayUUID = body.GatewayUUID
			sm.ToUUID = body.ToUUID
			sm.Name = "internal point mapping"
			sm.Description = "internal point mapping"
			sm.Enable = body.Enable
			sm.IsRemote = false
			sm.SubscriberType = body.SubscriberType
			sm.SubscriberApplication = body.SubscriberApplication
			subscription, err := d.CreateSubscription(&sm); if err != nil {
				 return nil, errorMsg("CreateSubscription", "error on trying to add", nil)
			}
			if err != nil {
				 return nil, errorMsg("CreateSubscription PointLedger", "error on trying to add", nil)
			}
			 fmt.Println(5555555, body.IsRemote, 5555555)
			u := utils.MakeTopicUUID("")
			d.DB.Create(&model.PointSubscriberLedger{UUID: u, GatewayUUID: body.GatewayUUID, SubscriberUUID: body.UUID, SubscriptionUUID: subscription.UUID, PointUUID: body.ToUUID})
		 }
		return body, nil

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
