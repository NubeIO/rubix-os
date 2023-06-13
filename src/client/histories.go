package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
	"time"
)

func (inst *FlowClient) GetPointHistoriesForSync(id int, timeStamp time.Time) (*[]model.PointHistory, error) {
	req := inst.client.R().
		SetResult(&[]model.PointHistory{}).SetQueryParam("id", fmt.Sprintf("%v", id)).
		SetQueryParam("timestamp", fmt.Sprintf("%v", timeStamp.Format(time.RFC3339Nano)))
	resp, err := nresty.FormatRestyResponse(req.Get("/api/histories/points/sync"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]model.PointHistory), nil
}
