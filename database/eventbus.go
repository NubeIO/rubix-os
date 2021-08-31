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
	query := d.DB.Where("uuid = ?", uuid).Find(&pointModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}

	//TODO make a better query
	thingProducer, err := d.GetProducerByThingUUID(pointModel.UUID);if errors.Is(err, gorm.ErrRecordNotFound) {

	}

	thingConsumer, err := d.GetWriterByThing(pointModel.UUID);if errors.Is(err, gorm.ErrRecordNotFound) {

	}

	if thingProducer  != nil {
		//fmt.Println("thingProducer", thingProducer.ConsumerUUID, thingProducer.ProducerUUID)
		//do check's like is enabled what is the cov, is the stream enabled, is the flow-network enabled
		producer, err := d.GetProducerByThingUUID(thingProducer.ProducerThingUUID)
		if err != nil {
			return nil, query.Error
		}
		thingConsumerUUID := ""
		if thingConsumer != nil {
			thingConsumerUUID = thingConsumer.ConsumerThingUUID
		}

		gateway, err := d.GetStreamGateway(producer.StreamUUID)
		if err != nil {
			return nil, query.Error
		}
		flowNetwork, err := d.GetFlowNetwork(gateway.StreamListUUID)
		if err != nil {
			return nil, query.Error
		}
		fmt.Println("POINT UUID", pointModel.UUID)
		fmt.Println("UPDATE FROM POINT", "FromThingUUID", producer.UUID)
		fmt.Println("GATEWAY", "producer", producer.UUID, "gateway", gateway.UUID)
		fmt.Println("FLOW-NETWORK","flow-network", flowNetwork.UUID, "flow-network-remote-uuid", flowNetwork.RemoteFlowUUID)
		fmt.Println("SEND DATA TO", thingConsumerUUID, "Description", pointModel.Description)
	}
	if thingConsumer  != nil {
		fmt.Println("thingConsumer", thingConsumer.ConsumerThingUUID, thingConsumer.ConsumerUUID)
		//if thingProducer.ProducerUUID != "" || thingConsumer.ProducerThingUUID != "" {
		//	busUpdate(pointModel.UUID, "updates", pointModel)
		//}
	}


	return pointModel, nil
}
