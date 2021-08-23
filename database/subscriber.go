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
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Subscriber)
	if body.SubscriberType == model.CommonNaming.Point  {
		//call points and make it exists
		gateway, err := d.GetGateway(body.GatewayUUID)
		if err != nil {
			return nil, errorMsg("GetGateway", "error on trying to get", nil)
		}
		query, err := d.GetPoint(body.ToUUID, false)
		if err != nil {
			return nil, errorMsg("CreateSubscriber", "error on trying to add", nil)
		}
		if query != nil {
			body.SubscriberApplication = model.CommonNaming.Mapping
			n := d.DB.Create(&body).Error

			fmt.Println(gateway.UUID, "gateway")
			// if its local
			 if !gateway.IsRemote {
				 var sm model.Subscription
				 sm.GatewayUUID = body.GatewayUUID
				 sm.ToUUID = body.ToUUID
				 sm.Name = "internal point mapping"
				 sm.Description = "internal point mapping"
				 sm.Enable = body.Enable
				 sm.SubscriberType = body.SubscriberType
				 sm.SubscriberApplication = body.SubscriberApplication
				 subscription, err := d.CreateSubscription(&sm)
				 if err != nil {
					return nil, errorMsg("CreateSubscription", "error on trying to add", nil)
				 }
				 u := utils.MakeTopicUUID("")
				 err = 	d.DB.Create(&model.PointSubscriberLedger{UUID: u, GatewayUUID: body.GatewayUUID, SubscriberUUID: body.UUID, SubscriptionUUID: subscription.UUID, PointUUID: body.ToUUID}).Error
				 if err != nil {
					return nil, errorMsg("CreateSubscription PointLedger", "error on trying to add", nil)
				 } else {
					 u := utils.MakeTopicUUID("")
					d.DB.Create(&model.PointSubscriptionLedger{UUID: u, GatewayUUID: body.GatewayUUID, SubscriberUUID: body.UUID, PointUUID: body.FromUUID})
				 }
			 }
			if err != nil {
				return nil, err
			}
			return body, n
		}
	} else if body.SubscriberType == model.CommonNaming.Network {

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
