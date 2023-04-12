package csrest

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"

	"github.com/NubeIO/flow-framework/nresty"
	"github.com/go-resty/resty/v2"
)

var limit = "200"

const orgID = "1"

func (inst *ChirpClient) SetDeviceLimit(newLimit int) {
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
// i.e. "401 Unauthorised" or "connection refused"
func IsCSConnectionError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "401") ||
		strings.Contains(err.Error(), "authentication failed") ||
		strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no route to host") ||
		strings.Contains(err.Error(), "501 Not Implemented"))
}

// GetOrganizations get all
func (inst *ChirpClient) GetOrganizations() (*Organizations, error) {
	q := fmt.Sprintf("/organizations?limit=%s", limit)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(Organizations{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*Organizations), nil
}

// GetGateways get all
func (inst *ChirpClient) GetGateways() (*Gateways, error) {
	q := fmt.Sprintf("/gateways?limit=%s", limit)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(Gateways{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*Gateways), nil
}

// GetApplications get all
func (inst *ChirpClient) GetApplications() (*Applications, error) {
	q := fmt.Sprintf("/applications?limit=%s", limit)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(Applications{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*Applications), nil
}

// GetDeviceProfiles get all
func (inst *ChirpClient) GetDeviceProfiles() (*DeviceProfiles, error) {
	q := fmt.Sprintf("/device-profiles?limit=%s", limit)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(DeviceProfiles{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*DeviceProfiles), nil
}

// GetServiceProfiles get all
func (inst *ChirpClient) GetServiceProfiles() (*ServiceProfiles, error) {
	q := fmt.Sprintf("/service-profiles?limit=%s&organizationID=%s", limit, orgID)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(ServiceProfiles{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ServiceProfiles), nil
}

// GetGatewayProfiles get all
func (inst *ChirpClient) GetGatewayProfiles() (*GatewayProfiles, error) {
	q := fmt.Sprintf("/gateway-profiles?limit=%s&organizationID=%s", limit, orgID)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(GatewayProfiles{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*GatewayProfiles), nil
}

// GetDevices all
func (inst *ChirpClient) GetDevices() (*csmodel.Devices, error) {
	var allDevices csmodel.Devices
	for _, application := range csApplications.Result {
		q := fmt.Sprintf("/devices?limit=%s&applicationID=%s", limit, application.ID)
		resp, err := nresty.FormatRestyResponse(inst.client.R().
			SetResult(csmodel.Devices{}).
			Get(q))
		err = checkResponse(resp, err)
		if err != nil {
			log.Error("lorawan: rest GetDevices error: ", err)
			return nil, err
		}
		if resp.Result() == nil {
			log.Error("lorawan: rest GetDevices result nil", err)
		}
		currDevices := resp.Result().(*csmodel.Devices)
		allDevices.Result = append(allDevices.Result, currDevices.Result...)
	}
	return &allDevices, nil
}

// GetDevice single
func (inst *ChirpClient) GetDevice(devEui string) (*csmodel.Device, error) {
	q := fmt.Sprintf("/devices/%s", devEui)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(csmodel.DeviceAll{}).
		Get(q))
	err = checkResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return &resp.Result().(*csmodel.DeviceAll).Device, nil
}

// AddDevice add all
func (inst *ChirpClient) AddDevice(body *Device) (*Device, error) {
	q := fmt.Sprintf("/devices")
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(Device{}).
		SetBody(body).
		Post(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*Device), nil
}

// EditDevice edit object
func (inst *ChirpClient) EditDevice(devEui string, body *Device) (*Device, error) {
	q := fmt.Sprintf("/devices/%s", devEui)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(Device{}).
		SetBody(body).
		Put(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*Device), nil
}

// DeleteDevice delete
func (inst *ChirpClient) DeleteDevice(devEui string) (bool, error) {
	q := fmt.Sprintf("/devices/%s", devEui)
	_, err := nresty.FormatRestyResponse(inst.client.R().
		Delete(q))
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeviceOTAKeysUpdate active a device
func (inst *ChirpClient) DeviceOTAKeysUpdate(devEui string, body *DeviceKey) (*DeviceKey, error) {
	q := fmt.Sprintf("/devices/%s/keys", devEui)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(DeviceKey{}).
		SetBody(body).
		Put(q))
	if err != nil {
		return nil, err
	}
	r := resp.Result().(*DeviceKey)
	return r, nil
}

// DeviceOTAKeys active a device
func (inst *ChirpClient) DeviceOTAKeys(devEui string, body *DeviceKey) (*DeviceKey, error) {
	q := fmt.Sprintf("/devices/%s/keys", devEui)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(DeviceKey{}).
		SetBody(body).
		Post(q))
	if err != nil {
		return nil, err
	}
	r := resp.Result().(*DeviceKey)
	return r, nil
}

// ActivateDevice activate a device
func (inst *ChirpClient) ActivateDevice(devEui string, body *DeviceActivation) (*DeviceActivation, error) {
	q := fmt.Sprintf("/devices/%s/activate", devEui)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(DeviceActivation{}).
		SetBody(body).
		Put(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*DeviceActivation), nil
}
