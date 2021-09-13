package rest

import (
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)

func WriteClone(uuid string, flowBody *model.FlowNetwork, body *model.WriterClone, write bool) (*model.WriterClone, error) {
	if !flowBody.IsMQTT {
		ip := flowBody.FlowIP
		port := flowBody.FlowPort
		token := flowBody.FlowToken
		c := client.NewSessionWithToken(token, ip, port)
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
	return nil, nil
}

func WriteProducer(uuid string, flowBody *model.FlowNetwork, body *model.Producer, write bool) (*model.Producer, error) {
	if !flowBody.IsMQTT {
		ip := flowBody.FlowIP
		port := flowBody.FlowPort
		token := flowBody.FlowToken
		c := client.NewSessionWithToken(token, ip, port)
		if write {
			res, err := c.ClientEditProducer(uuid, *body)
			if err != nil {
				return nil, err
			}
			return res, err
		} else {
			res, err := c.ClientGetProducer(uuid)
			if err != nil {
				return nil, err
			}
			return res, err
		}
	}
	return nil, nil
}

func ProducerRead(flowBody *model.FlowNetwork, producerUUID string) (*model.Producer, error) {
	ip := flowBody.FlowIP
	port := flowBody.FlowPort
	token := flowBody.FlowToken
	c := client.NewSessionWithToken(token, ip, port)
	point, err := c.ClientGetProducer(producerUUID)
	if err != nil {
		return nil, err
	}
	return point, err
}

func ProducerHistory(flowBody *model.FlowNetwork, producerUUID string) (*model.ProducerHistory, error) {
	ip := flowBody.FlowIP
	port := flowBody.FlowPort
	token := flowBody.FlowToken
	c := client.NewSessionWithToken(token, ip, port)
	point, err := c.ClientGetHistory(producerUUID)
	if err != nil {
		return nil, err
	}
	return point, err
}
