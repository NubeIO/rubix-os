package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
)

func (d *GormDatabase) producerBroadcast(producer model.ProducerBody) error {
	t := fmt.Sprintf("%s", eventbus.ProducerEvent)
	stream, err := d.GetStream(producer.StreamUUID, false)
	if err != nil {
		return err
	}
	flowUUID := ""
	for _, net := range stream.FlowNetworks {
		flowUUID = net.UUID
	}
	//TODO
	// check if flow is enabled
	// check if stream is enabled
	// check if producer is enabled
	// then broadcast
	producer.FlowNetworkUUID = flowUUID
	d.Bus.RegisterTopic(t)
	err = d.Bus.Emit(eventbus.CTX(), t, producer)
	if err != nil {
		return err
	}
	return nil
}

//compare a COV event
func compare(p1, p2 *model.Point) bool {
	return p1.PresentValue == p2.PresentValue
}
