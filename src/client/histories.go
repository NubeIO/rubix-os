package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

func (a *FlowClient) GetProducerHistoriesPointsForSync(id int, timeStamp time.Time) (*[]model.History, error) {
	req := a.client.R().
		SetResult(&[]model.History{}).SetQueryParam("id", fmt.Sprintf("%v", id)).
		SetQueryParam("timestamp", fmt.Sprintf("%v", timeStamp.Format(time.RFC3339Nano)))
	resp, err := req.Get("/api/histories/producers/points_for_sync")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("GetProducerHistoriesPointsForSync: %s", err)
		} else {
			return nil, fmt.Errorf("GetProducerHistoriesPointsForSync: %s", resp)
		}
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetProducerHistoriesPointsForSync: %s", resp)
	}
	return resp.Result().(*[]model.History), nil
}
