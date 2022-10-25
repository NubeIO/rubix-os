package chirpstackrest

import (
	"errors"
	"github.com/NubeIO/flow-framework/nresty"
	"strconv"
)

type ChirpstackCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChirpstackToken struct {
	Token string `json:"jwt"`
}

type ChirpstackDevWrapper struct {
	Device ChirpstackDev `json:"device"`
}

type ChirpstackDev struct {
	ApplicationID     int         `json:"applicationID"`
	Description       string      `json:"description"`
	DeviceEUI         string      `json:"devEUI"`
	DeviceProfileID   string      `json:"deviceProfileID"`
	IsDisabled        bool        `json:"isDisabled"`
	Name              string      `json:"name"`
	ReferenceAltitude float64     `json:"referenceAltitude"`
	SkipFCntCheck     bool        `json:"skipFCntCheck"`
	Tags              interface{} `json:"tags"`
	Variables         interface{} `json:"variables"`
}

type ChirpstackDevProfileResponse struct {
	TotalCount string                 `json:"totalCount"`
	Result     []ChirpstackDevProfile `json:"result"`
}

type ChirpstackDevProfile struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	OrganizationID    string `json:"organizationID"`
	NetworkServerID   string `json:"networkServerID"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	NetworkServerName string `json:"skipFCntCheck"`
}

type ChirpstackDevActivateWrapper struct {
	DeviceKeys DevActivateDeviceKeys `json:"deviceKeys"`
}

type DevActivateDeviceKeys struct {
	ApplicationKey string `json:"appKey"`
	NetworkKey     string `json:"nwkKey"`
	DeviceEUI      string `json:"devEUI"`
}

func (a *RestClient) GetChirpstackToken(user, pass string) (*ChirpstackToken, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetBody(ChirpstackCredentials{
			Email:    user,
			Password: pass,
		}).
		SetResult(ChirpstackToken{}).
		Post("/api/internal/login"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ChirpstackToken), nil
}

func (a *RestClient) GetChirpstackDeviceProfileUUID(token string) ([]ChirpstackDevProfile, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Accept", "application/json").
		SetResult(ChirpstackDevProfileResponse{}).
		Get("/api/device-profiles?limit=50"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ChirpstackDevProfileResponse).Result, nil
}

func (a *RestClient) AddChirpstackDevice(chirpstackAppNum, modbusAddress int, deviceName, lorawanDeviceEUI, chirpstackDeviceProfileUUID, token string) error {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("content-type", "application/json").
		SetBody(ChirpstackDevWrapper{
			Device: ChirpstackDev{
				ApplicationID:     chirpstackAppNum,
				Description:       "Modbus Address: " + strconv.Itoa(modbusAddress),
				DeviceEUI:         lorawanDeviceEUI,
				DeviceProfileID:   chirpstackDeviceProfileUUID,
				IsDisabled:        false,
				Name:              deviceName,
				ReferenceAltitude: 0,
				SkipFCntCheck:     true,
			},
		}).
		SetResult(map[string]interface{}{}).
		Post("/api/devices"))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New("error response code: " + strconv.Itoa(resp.StatusCode()))
	}
	return nil
}

func (a *RestClient) ActivateChirpstackDevice(applicationKey, lorawanDeviceEUI, token, lorawanNetworkKey string) error {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("content-type", "application/json").
		SetBody(ChirpstackDevActivateWrapper{
			DeviceKeys: DevActivateDeviceKeys{
				ApplicationKey: applicationKey,
				NetworkKey:     lorawanNetworkKey,
				DeviceEUI:      lorawanDeviceEUI,
			},
		}).
		SetResult(map[string]interface{}{}).
		Post("/api/devices/" + lorawanDeviceEUI + "/keys"))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return errors.New("error activating chirpstack device (it might already be activated)")
	}
	return nil
}
