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

func (a *FlowClient) GetProducerHistoriesPoints(uuid string, lastSyncId int) (*[]model.History, error) {
	req := a.client.R().
		SetResult(&[]model.History{}).SetQueryParam("id_gt", fmt.Sprintf("%v", lastSyncId))
	resp, err := req.
		Get(fmt.Sprintf("api/fnc/%s/api/histories/producers/points", uuid))
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("GetPointsProducerHistories: %s", err)
		} else {
			return nil, fmt.Errorf("GetPointsProducerHistories: %s", resp)
		}
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetPointsProducerHistories: %s", resp)
	}
	return resp.Result().(*[]model.History), nil
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
