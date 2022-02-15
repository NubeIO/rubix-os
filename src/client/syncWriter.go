package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
)

func (a *FlowClient) SyncWriter(body *model.SyncWriter) (*model.WriterClone, error) {
	resp, err := a.client.R().
		SetResult(&model.WriterClone{}).
		SetBody(body).
		Post("/api/sync/writer")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("SyncWriter: %s", err)
		} else {
			return nil, fmt.Errorf("SyncWriter: %s", resp)
		}
	}
	return resp.Result().(*model.WriterClone), nil
}

func (a *FlowClient) SyncCOV(body *model.SyncCOV) error {
	resp, err := a.client.R().
		SetBody(body).
		Post("/api/sync/cov")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return fmt.Errorf("SyncCOV: %s", err)
		} else {
			return fmt.Errorf("SyncCOV: %s", resp)
		}
	}
	return nil
}

func (a *FlowClient) SyncWriterAction(body *model.SyncWriterAction) error {
	resp, err := a.client.R().
		SetBody(body).
		Post("/api/sync/writer_action")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return fmt.Errorf("SyncWriterAction: %s", err)
		} else {
			return fmt.Errorf("SyncWriterAction: %s", resp)
		}
	}
	return nil
}
