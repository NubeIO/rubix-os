package rest

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/client"
)

func WriteClone(uuid string, flowBody *model.FlowNetwork, body *model.WriterClone, write bool) (*model.WriterClone, error) {
	c := client.NewFlowClientCli(flowBody.FlowIP, flowBody.FlowPort, flowBody.FlowToken, flowBody.IsMasterSlave, flowBody.GlobalUUID, model.IsFNCreator(flowBody))
	if write {
		res, err := c.EditWriterClone(uuid, *body, write)
		if err != nil {
			return nil, err
		}
		return res, err
	} else {
		res, err := c.GetWriterClone(uuid)
		if err != nil {
			return nil, err
		}
		return res, err
	}
}

func WriteProducer(uuid string, flowBody *model.FlowNetwork, body *model.Producer, write bool) (*model.Producer, error) {
	c := client.NewFlowClientCli(flowBody.FlowIP, flowBody.FlowPort, flowBody.FlowToken, flowBody.IsMasterSlave, flowBody.GlobalUUID, model.IsFNCreator(flowBody))
	if write {
		res, err := c.EditProducer(uuid, *body)
		if err != nil {
			return nil, err
		}
		return res, err
	} else {
		res, err := c.GetProducer(uuid)
		if err != nil {
			return nil, err
		}
		return res, err
	}
}

func ProducerRead(flowBody *model.FlowNetwork, producerUUID string) (*model.Producer, error) {
	c := client.NewFlowClientCli(flowBody.FlowIP, flowBody.FlowPort, flowBody.FlowToken, flowBody.IsMasterSlave, flowBody.GlobalUUID, model.IsFNCreator(flowBody))
	point, err := c.GetProducer(producerUUID)
	if err != nil {
		return nil, err
	}
	return point, err
}

func ProducerHistory(flowBody *model.FlowNetwork, producerUUID string) (*model.ProducerHistory, error) {
	c := client.NewFlowClientCli(flowBody.FlowIP, flowBody.FlowPort, flowBody.FlowToken, flowBody.IsMasterSlave, flowBody.GlobalUUID, model.IsFNCreator(flowBody))
	point, err := c.GetProducerHistory(producerUUID)
	if err != nil {
		return nil, err
	}
	return point, err
}
