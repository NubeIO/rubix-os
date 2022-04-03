package plgrest

import (
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnetmodel"
	model "github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwmodel"
	baseModel "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"strconv"
)

// GetPoints all points
func (a *RestClient) GetPoints() (*[]bacnetmodel.BacnetPoint, error) {
	resp, err := a.client.R().
		SetResult([]bacnetmodel.BacnetPoint{}).
		Get("/api/bacnet/points")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*[]bacnetmodel.BacnetPoint), nil
}

// AddPoint an object
func (a *RestClient) AddPoint(body bacnetmodel.BacnetPoint) (*bacnetmodel.BacnetPoint, error) {
	fmt.Println("ADD POINT ON IN BACNET REST CALL", body)
	resp, err := a.client.R().
		SetResult(&bacnetmodel.BacnetPoint{}).
		SetBody(body).
		Post("/api/bacnet/points")
	if err != nil {
		return nil, fmt.Errorf("failed to add %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*bacnetmodel.BacnetPoint), nil
}

type PriorityWriter struct {
	PriorityArrayWrite baseModel.Priority `json:"priority_array_write,omitempty"`
}

// EditPoint an object
func (a *RestClient) EditPoint(body bacnetmodel.BacnetPoint, obj string, addr int) (*bacnetmodel.BacnetPoint, error) {
	priorityWriter := new(PriorityWriter)
	priorityWriter.PriorityArrayWrite = *body.Priority
	resp, err := a.client.R().
		SetResult(&bacnetmodel.BacnetPoint{}).
		SetBody(priorityWriter).
		SetPathParams(map[string]string{"obj": obj, "addr": strconv.Itoa(addr)}).
		Patch("/api/bacnet/points/obj/{obj}/{addr}")
	if err != nil {
		return nil, fmt.Errorf("failed to update %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*bacnetmodel.BacnetPoint), nil
}

// DeletePoint an object
func (a *RestClient) DeletePoint(obj string, addr int) (bool, error) {
	resp, err := a.client.R().
		SetPathParams(map[string]string{"obj": obj, "addr": strconv.Itoa(addr)}).
		Delete("/api/bacnet/points/obj/{obj}/{addr}")
	if err != nil {
		return false, fmt.Errorf("failed to delete %s", err)
	}
	if resp.Error() != nil {
		return false, getAPIError(resp)
	}
	return true, nil
}

// PingServer all points
func (a *RestClient) PingServer() (*bacnetmodel.ServerPing, error) {
	resp, err := a.client.R().
		SetResult(&bacnetmodel.ServerPing{}).
		Get("/api/system/ping")
	if err != nil {
		return nil, fmt.Errorf("error geting server %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*bacnetmodel.ServerPing), nil
}

// GetServer all points
func (a *RestClient) GetServer() (*model.Server, error) {
	resp, err := a.client.R().
		SetResult(&model.Server{}).
		Get("/api/bacnet/server")
	if err != nil {
		return nil, fmt.Errorf("error geting server %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.Server), nil
}

// EditServer an object
func (a *RestClient) EditServer(body model.Server) (*model.Server, error) {
	resp, err := a.client.R().
		SetResult(&model.Server{}).
		SetBody(body).
		Patch("/api/bacnet/server")
	if err != nil {
		return nil, fmt.Errorf("failed to update %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.Server), nil
}
