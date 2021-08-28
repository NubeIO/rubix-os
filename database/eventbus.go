package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"gorm.io/gorm"
)

// DBBusEvent and event on the bus.
func (d *GormDatabase) DBBusEvent(uuid string, body *model.Point) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Find(&pointModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&pointModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}

	//TODO make a better query
	thingSubscriber, err := d.GetSubscriberListByThing(pointModel.UUID);if errors.Is(err, gorm.ErrRecordNotFound) {

	}

	thingSubscription, err := d.GetSubscriptionListByThing(pointModel.UUID);if errors.Is(err, gorm.ErrRecordNotFound) {

	}

	if thingSubscriber  != nil {
		fmt.Println("thingSubscriber", thingSubscriber.FromThingUUID, thingSubscriber.SubscriberUUID)
		//do check's like is enabled what is the cov, is the stream enabled, is the flow-network enabled
		subscriber, err := d.GetSubscriber(thingSubscriber.SubscriberUUID)
		if err != nil {
			return nil, query.Error
		}
		thingSubscriptionUUID := ""
		if thingSubscription != nil {
			thingSubscriptionUUID = thingSubscription.ToThingUUID
		}

		gateway, err := d.GetStreamGateway(subscriber.StreamUUID)
		if err != nil {
			return nil, query.Error
		}
		flowNetwork, err := d.GetFlowNetwork(gateway.FlowNetworkUUID)
		if err != nil {
			return nil, query.Error
		}
		fmt.Println("POINT UUID", pointModel.UUID)
		fmt.Println("UPDATE FROM POINT", "FromThingUUID", subscriber.FromThingUUID)
		fmt.Println("GATEWAY", "subscriber", subscriber.UUID, "gateway", gateway.UUID)
		fmt.Println("FLOW-NETWORK","flow-network", flowNetwork.UUID, "flow-network-remote-uuid", flowNetwork.RemoteUUID)
		fmt.Println("SEND DATA TO", thingSubscriptionUUID, "Description", pointModel.Description)
	}
	if thingSubscription  != nil {
		fmt.Println("thingSubscription", thingSubscription.ToThingUUID, thingSubscription.SubscriptionUUID)
		if thingSubscriber.FromThingUUID != "" || thingSubscription.ToThingUUID != "" {
			busUpdate(pointModel.UUID, "updates", pointModel)
		}
	}


	return pointModel, nil
}
