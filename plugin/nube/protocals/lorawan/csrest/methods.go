package csrest

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/go-resty/resty/v2"
)

// TODO: handle reconnect

var limit = "200"

func (a *RestClient) SetDeviceLimit(newLimit int) {
	limit = strconv.Itoa(newLimit)
}

func checkResponse(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return errors.New(resp.Status())
	}
	return err
}

// IsCSConnectionError returns true if error is related to connection.
//  i.e. "401 Unauthorised" or "connection refused"
func IsCSConnectionError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "401") ||
		strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no route to host") ||
		strings.Contains(err.Error(), "501 Not Implemented"))
}

// GetDevices all
func (a *RestClient) GetDevices() (*csmodel.Devices, error) {
	q := fmt.Sprintf("/api/devices?limit=%s", limit)
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetResult(csmodel.Devices{}).
		Get(q))
	err = checkResponse(resp, err)
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
	err = checkResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*csmodel.Device), nil
}
