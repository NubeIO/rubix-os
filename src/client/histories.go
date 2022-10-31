package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

func (inst *FlowClient) GetProducerHistoriesPointsForSync(id int, timeStamp time.Time) (*[]model.History, error) {
	req := inst.client.R().
		SetResult(&[]model.History{}).SetQueryParam("id", fmt.Sprintf("%v", id)).
		SetQueryParam("timestamp", fmt.Sprintf("%v", timeStamp.Format(time.RFC3339Nano)))
	resp, err := nresty.FormatRestyResponse(req.Get("/api/histories/producers/points_for_sync"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]model.History), nil
}
