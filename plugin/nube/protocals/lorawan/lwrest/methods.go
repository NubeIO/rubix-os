package lwrest

import (
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwmodel"
	"github.com/NubeIO/flow-framework/src/client"
)

const limit = "50"
const orgID = "1"

// GetOrganizations get all
func (a *RestClient) GetOrganizations() (*lwmodel.Organizations, error) {
	q := fmt.Sprintf("/api/organizations?limit=%s", limit)
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.Organizations{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.Organizations), nil
}

// GetGateways get all
func (a *RestClient) GetGateways() (*lwmodel.Gateways, error) {
	q := fmt.Sprintf("/api/gateways?limit=%s&organizationID=%s", limit, orgID)
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.Gateways{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.Gateways), nil
}

// GetApplications get all
func (a *RestClient) GetApplications() (*lwmodel.Applications, error) {
	q := fmt.Sprintf("/api/applications?limit=%s&organizationID=%s", limit, orgID)
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.Applications{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.Applications), nil
}

// GetDevices get all
func (a *RestClient) GetDevices() (*lwmodel.Devices, error) {
	q := fmt.Sprintf("/api/devices?limit=%s&organizationID=%s", limit, orgID)
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.Devices{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.Devices), nil
}

// GetDeviceProfiles get all
func (a *RestClient) GetDeviceProfiles() (*lwmodel.DeviceProfiles, error) {
	q := fmt.Sprintf("/api/device-profiles?limit=%s&organizationID=%s", limit, orgID)
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.DeviceProfiles{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.DeviceProfiles), nil
}

// GetServiceProfiles get all
func (a *RestClient) GetServiceProfiles() (*lwmodel.ServiceProfiles, error) {
	q := fmt.Sprintf("/api/service-profiles?limit=%s&organizationID=%s", limit, orgID)
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.ServiceProfiles{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.ServiceProfiles), nil
}

// GetGatewayProfiles get all
func (a *RestClient) GetGatewayProfiles() (*lwmodel.GatewayProfiles, error) {
	q := fmt.Sprintf("/api/gateway-profiles?limit=%s&organizationID=%s", limit, orgID)
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.GatewayProfiles{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.GatewayProfiles), nil
}

// AddDevice add all
func (a *RestClient) AddDevice(body lwmodel.Devices) (*lwmodel.Devices, error) {
	q := fmt.Sprintf("/api/devices")
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.Devices{}).
		SetBody(body).
		Post(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.Devices), nil
}

// GetDevice get an object
func (a *RestClient) GetDevice(devEui string) (*lwmodel.GetDevice, error) {
	q := fmt.Sprintf("/api/devices/%s", devEui)
	resp, err := client.CheckError(a.client.R().
		SetResult(lwmodel.GetDevice{}).
		Get(q))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*lwmodel.GetDevice), nil
}

// EditDevice edit object
func (a *RestClient) EditDevice(devEui string, body lwmodel.Device) (bool, error) {
	q := fmt.Sprintf("/api/devices/%s", devEui)
	_, err := client.CheckError(a.client.R().
		SetResult(lwmodel.Device{}).
		SetBody(body).
		Put(q))
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteDevice delete
func (a *RestClient) DeleteDevice(devEui string) (bool, error) {
	q := fmt.Sprintf("/api/devices/%s", devEui)
	_, err := client.CheckError(a.client.R().
		Delete(q))
	if err != nil {
		return false, err
	}
	return true, nil
}
