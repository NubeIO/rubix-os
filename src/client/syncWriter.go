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

func (a *FlowClient) SyncWriterCOV(body *model.SyncWriterCOV) error {
	resp, err := a.client.R().
		SetBody(body).
		Post("/api/sync/writer_cov")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return fmt.Errorf("SyncWriterCOV: %s", err)
		} else {
			return fmt.Errorf("SyncWriterCOV: %s", resp)
		}
	}
	return nil
}
