package eventbus

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)

func EventREST(flowBody  *model.FlowNetwork, producerBody *model.Producer, write bool) (*model.Producer, error) {
	if !flowBody.IsMQTT {
		ip := flowBody.FlowIP
		port := flowBody.FlowPort
		token := flowBody.FlowToken
		producerUUID := producerBody.UUID
		c := client.NewSessionWithToken(token, ip, port)
		if write {
			point, err := c.ClientEditProducer(producerUUID, *producerBody);if err != nil {
				return nil, err
			}
			return point, err
		} else {
			point, err := c.ClientGetProducer(producerUUID);if err != nil {
				return nil, err
			}
			return point, err
		}
	}
	return nil, nil
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
