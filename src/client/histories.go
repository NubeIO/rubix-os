package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
)

func (a *FlowClient) GetProducerHistoriesPoints(uuid string, lastSyncId int) (*[]model.History, error) {
	req := a.client.R().
		SetResult(&[]model.History{}).SetQueryParam("id_gt", fmt.Sprintf("%v", lastSyncId))
	resp, err := req.
		Get(fmt.Sprintf("api/fnc/%s/api/histories/producers/points", uuid)) // TODO: url check, `api/fnc/%s` is not needed most probably
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
