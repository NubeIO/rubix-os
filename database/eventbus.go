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
	thingProducer, err := d.GetProducerListByThing(pointModel.UUID);if errors.Is(err, gorm.ErrRecordNotFound) {

	}

	thingSubscription, err := d.GetSubscriptionListByThing(pointModel.UUID);if errors.Is(err, gorm.ErrRecordNotFound) {

	}

	if thingProducer  != nil {
		fmt.Println("thingProducer", thingProducer.SubscriptionUUID, thingProducer.ProducerUUID)
		//do check's like is enabled what is the cov, is the stream enabled, is the flow-network enabled
		producer, err := d.GetProducer(thingProducer.ProducerUUID)
		if err != nil {
			return nil, query.Error
		}
		thingSubscriptionUUID := ""
		if thingSubscription != nil {
			thingSubscriptionUUID = thingSubscription.ProducerThingUUID
		}

		gateway, err := d.GetStreamGateway(producer.StreamUUID)
		if err != nil {
			return nil, query.Error
		}
		flowNetwork, err := d.GetFlowNetwork(gateway.FlowNetworkUUID)
		if err != nil {
			return nil, query.Error
		}
		fmt.Println("POINT UUID", pointModel.UUID)
		fmt.Println("UPDATE FROM POINT", "FromThingUUID", producer.UUID)
		fmt.Println("GATEWAY", "producer", producer.UUID, "gateway", gateway.UUID)
		fmt.Println("FLOW-NETWORK","flow-network", flowNetwork.UUID, "flow-network-remote-uuid", flowNetwork.RemoteUUID)
		fmt.Println("SEND DATA TO", thingSubscriptionUUID, "Description", pointModel.Description)
	}
	if thingSubscription  != nil {
		fmt.Println("thingSubscription", thingSubscription.ProducerThingUUID, thingSubscription.SubscriptionUUID)
		if thingProducer.ProducerUUID != "" || thingSubscription.ProducerThingUUID != "" {
			busUpdate(pointModel.UUID, "updates", pointModel)
		}
	}


	return pointModel, nil
}
