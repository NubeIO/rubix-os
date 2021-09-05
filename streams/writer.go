package streams

import (
	"encoding/json"
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


//ValidateTypes check the type of the producer or consumer, as in type=point
func ValidateTypes(t string, body *model.WriterBody) ([]byte, string, error) {
	if t == model.CommonNaming.Point {
		var bk model.WriterBody
		if body.Action == model.CommonNaming.Write {
			if body.Priority == bk.Priority {
				return nil, body.Action, errors.New("error: invalid json on writerBody")
			}
			b, err := json.Marshal(body.Priority)
			if err != nil {
				return nil, body.Action, errors.New("error: failed to marshal json on writeBody")
			}
			return b, body.Action, err
		} else {
			if body.Action == model.CommonNaming.Read {
				return nil, body.Action, nil
			} else {
				return nil, body.Action, errors.New("error: invalid action, try read or write")
			}
		}
	}
	return nil, body.Action, errors.New("error: invalid data type on writerBody, ie type could be a point")

}
