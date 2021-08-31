package rest

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)

func WriteClone(uuid string, flowBody *model.FlowNetwork, producerBody *model.WriterClone, write bool) (*model.WriterClone, error) {
	if !flowBody.IsMQTT {
		ip := flowBody.FlowIP
		port := flowBody.FlowPort
		token := flowBody.FlowToken
		c := client.NewSessionWithToken(token, ip, port)
		if write {
			point, err := c.ClientEditWriterClone(uuid, *producerBody);if err != nil {
				return nil, err
			}
			return point, err
		} else {
			point, err := c.ClientGetWriterClone(uuid);if err != nil {
				return nil, err
			}
			return point, err
		}
	}
	return nil, nil
}


func ProducerRead(flowBody  *model.FlowNetwork, producerUUID string) (*model.Producer, error) {
	ip := flowBody.FlowIP
	port := flowBody.FlowPort
	token := flowBody.FlowToken
	c := client.NewSessionWithToken(token, ip, port)
	point, err := c.ClientGetProducer(producerUUID);if err != nil {
		return nil, err
	}
	return point, err
}


func EventRESTPoint(pointUUID string, flowBody *model.FlowNetwork, pointBody *model.Point, write bool) (*model.Point, error) {
	if !flowBody.IsMQTT {
		ip := flowBody.FlowIP
		port := flowBody.FlowPort
		token := flowBody.FlowToken
		pointUUID := pointUUID
		fmt.Println(pointUUID, 99999999999)
		c := client.NewSessionWithToken(token, ip, port)
		if write {
			point, err := c.ClientEditPoint(pointUUID, *pointBody);if err != nil {
				return nil, err
			}
			return point, err
		} else {
			point, err := c.ClientGetPoint(pointUUID);if err != nil {
				return nil, err
			}
			return point, err
		}
	}
	return nil, nil
}
