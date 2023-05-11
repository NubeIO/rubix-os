package csrest

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/NubeIO/flow-framework/nresty"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var limit = "200"

func (inst *ChirpClient) Proxy(ctx *gin.Context) {
	ctx.Request.Header.Del("Authorization")
	ctx.Request.Header.Set("Grpc-Metadata-Authorization", inst.ClientToken)
	inst.proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

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

// GetDevices all
func (inst *ChirpClient) GetDevices() (*Devices, error) {
	var allDevices Devices
	for _, application := range csApplications.Result {
		q := fmt.Sprintf("/devices?limit=%s&applicationID=%s", limit, application.ID)
		resp, err := nresty.FormatRestyResponse(inst.client.R().
			SetResult(Devices{}).
			Get(q))
		err = checkResponse(resp, err)
		if err != nil {
			log.Error("lorawan: rest GetDevices error: ", err)
			return nil, err
		}
		if resp.Result() == nil {
			log.Error("lorawan: rest GetDevices result nil", err)
		}
		currDevices := resp.Result().(*Devices)
		allDevices.Result = append(allDevices.Result, currDevices.Result...)
	}
	return &allDevices, nil
}

// GetDevice single
func (inst *ChirpClient) GetDevice(devEui string) (*DeviceSingle, error) {
	q := fmt.Sprintf("/devices/%s", devEui)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(DeviceSingle{}).
		Get(q))
	err = checkResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*DeviceSingle), nil
}

// AddDevice add cs device
func (inst *ChirpClient) AddDevice(body *DeviceSingle) error {
	q := "/devices"
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Post(q))
	return err
}

// UpdateDevice update cs device
func (inst *ChirpClient) UpdateDevice(body *DeviceSingle) error {
	q := fmt.Sprintf("/devices/%s", body.Device.DevEUI)
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Put(q))
	return err
}

// DeleteDevice delete
func (inst *ChirpClient) DeleteDevice(devEui string) error {
	q := fmt.Sprintf("/devices/%s", devEui)
	_, err := nresty.FormatRestyResponse(inst.client.R().
		Delete(q))
	return err
}

// DeviceOTAAKeyGet get cs otaa device key
func (inst *ChirpClient) DeviceOTAAKeyGet(devEui string) (string, error) {
	q := fmt.Sprintf("/devices/%s/keys", devEui)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(DeviceKey{}).
		Get(q))
	return resp.Result().(*DeviceKey).Keys.NwkKey, err
}

// DeviceOTAAKeyAdd set cs otaa device key
func (inst *ChirpClient) DeviceOTAAKeyAdd(devEui string, key string) error {
	keys := DeviceKey{
		Keys: DeviceKeys{
			AppKey: key,
			NwkKey: key,
		},
	}
	q := fmt.Sprintf("/devices/%s/keys", devEui)
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(keys).
		Post(q))
	return err
}

// DeviceOTAAKeyPost update cs otaa device key
func (inst *ChirpClient) DeviceOTAAKeyUpdate(devEui string, key string) error {
	keys := DeviceKey{
		Keys: DeviceKeys{
			AppKey: key,
			NwkKey: key,
		},
	}
	q := fmt.Sprintf("/devices/%s/keys", devEui)
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(keys).
		Put(q))
	return err
}

// ActivateDevice activate a device
func (inst *ChirpClient) ActivateDevice(devEui string, body *DeviceActivation) error {
	q := fmt.Sprintf("/devices/%s/activate", devEui)
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Put(q))
	return err
}
