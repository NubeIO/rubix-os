package streams

import (
	"errors"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/rest"
)



func ProducerFeedback(producerUUID string, flowBody *model.FlowNetwork) (*model.ProducerHistory, error)  {
	if producerUUID == "" {
		return nil, errors.New("error: producer uuid is none")
	}
	producerFeedback, err := ProducerHistory(flowBody, producerUUID);if err != nil {
		return nil, errors.New("error: on get feedback from producer history")
	}
	return producerFeedback, err
}

func WriteClone(uuid string, flowBody *model.FlowNetwork, body *model.WriterClone, write bool) (*model.WriterClone, error) {
	if !flowBody.IsMQTT {
		call, err := rest.WriteClone(uuid, flowBody, body, write)
		if err != nil {
			return nil, err
		}
		return call, nil
	}
	return nil, nil
}


func ProducerHistory(flowBody  *model.FlowNetwork, producerUUID string) (*model.ProducerHistory, error) {
	producerFeedback, err := rest.ProducerHistory(flowBody, producerUUID);if err != nil {
		return nil, err
	}
	return producerFeedback, err
}
