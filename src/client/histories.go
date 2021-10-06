package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)

func (a *FlowClient) GetProducerHistory(uuid string) (*model.ProducerHistory, error) {
	resp, err := a.client.R().
		SetResult(&model.ProducerHistory{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/histories/producers/{uuid}/one")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("GetProducerHistory: %s", err)
		} else {
			return nil, fmt.Errorf("GetProducerHistory: %s", resp)
		}
	}
	return resp.Result().(*model.ProducerHistory), nil
}

func (a *FlowClient) AddProducerHistory(body model.ProducerHistory) (bool, error) {
	resp, err := a.client.R().
		SetBody(body).
		Post("/api/histories/producers")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return false, fmt.Errorf("AddProducerHistory: %s", err)
		} else {
			return false, fmt.Errorf("AddProducerHistory: %s", resp)
		}
	}
	return true, nil
}
