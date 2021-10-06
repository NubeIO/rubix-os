package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)

func (a *FlowClient) SyncStream(body *model.SyncStream) (*model.StreamClone, error) {
	resp, err := a.client.R().
		SetResult(&model.StreamClone{}).
		SetBody(body).
		Post("/api/sync/stream")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("SyncStream: %s", err)
		} else {
			return nil, fmt.Errorf("SyncStream: %s", resp)
		}
	}
	return resp.Result().(*model.StreamClone), nil
}
