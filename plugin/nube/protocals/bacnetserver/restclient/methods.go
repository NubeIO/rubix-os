package pkgrest

import (
	"fmt"
	plgmodel "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
)

// AddPoint an object
func (a *RestClient) AddPoint(body plgmodel.BacnetPoint) (*plgmodel.BacnetPoint, error) {
	resp, err := a.client.R().
		SetResult(&plgmodel.BacnetPoint{}).
		SetBody(body).
		Post("/api/bacnet/points")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*plgmodel.BacnetPoint), nil
}

// EditPoint an object
func (a *RestClient) EditPoint(body plgmodel.BacnetPoint) (*plgmodel.BacnetPoint, error) {
	addr := body.Address
	obj := body.ObjectType
	u := fmt.Sprintf("/api/bacnet/points/obj/%s/%d", obj, addr)
	resp, err := a.client.R().
		SetResult(&plgmodel.BacnetPoint{}).
		SetBody(body).
		Patch(u)
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*plgmodel.BacnetPoint), nil
}
