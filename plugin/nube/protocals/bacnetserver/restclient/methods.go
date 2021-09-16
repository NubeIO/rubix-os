package pkgrest

import (
	"fmt"
	pkgmodel "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
	"strconv"
)

// GetPoints all points
func (a *RestClient) GetPoints() (*[]pkgmodel.BacnetPoint, error) {
	resp, err := a.client.R().
		SetResult([]pkgmodel.BacnetPoint{}).
		Get("/api/bacnet/points")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*[]pkgmodel.BacnetPoint), nil
}

// AddPoint an object
func (a *RestClient) AddPoint(body pkgmodel.BacnetPoint) (*pkgmodel.BacnetPoint, error) {
	fmt.Println("ADD POINT ON IN BACNET REST CALL", body)
	resp, err := a.client.R().
		SetResult(&pkgmodel.BacnetPoint{}).
		SetBody(body).
		Post("/api/bacnet/points")
	if err != nil {
		return nil, fmt.Errorf("failed to add %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*pkgmodel.BacnetPoint), nil
}

// EditPoint an object
func (a *RestClient) EditPoint(body pkgmodel.BacnetPoint, obj string, addr int) (*pkgmodel.BacnetPoint, error) {
	resp, err := a.client.R().
		SetResult(&pkgmodel.BacnetPoint{}).
		SetBody(body).
		SetPathParams(map[string]string{"obj": obj, "addr": strconv.Itoa(addr)}).
		Patch("/api/bacnet/points/obj/{obj}/{addr}")
	if err != nil {
		return nil, fmt.Errorf("failed to update %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*pkgmodel.BacnetPoint), nil
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
