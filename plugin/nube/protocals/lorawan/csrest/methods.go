package csrest

import (
	"fmt"
	"strconv"

	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
)

// TODO: add to config
var limit = "200"

func (a *RestClient) SetDeviceLimit(newLimit int) {
	limit = strconv.Itoa(newLimit)
}

// GetDevices all
func (a *RestClient) GetDevices() (*csmodel.Devices, error) {
	q := fmt.Sprintf("/api/devices?limit=%s", limit)
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetResult(csmodel.Devices{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*csmodel.Devices), nil
}

// GetDevice single
func (a *RestClient) GetDevice(devEui string) (*csmodel.Device, error) {
	q := fmt.Sprintf("/api/devices/%s", devEui)
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetResult(csmodel.Device{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*csmodel.Device), nil
}
