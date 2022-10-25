package chirpstackrest

import (
	"github.com/NubeIO/flow-framework/nresty"
)

type ChirpstackCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChirpstackToken struct {
	Token string `json:"jwt"`
}

type ChirpstackDev struct {
	ApplicationID     string      `json:"applicationID"`
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

type ChirpstackDevProfile struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	OrganizationID    string `json:"organizationID"`
	NetworkServerID   string `json:"networkServerID"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	NetworkServerName string `json:"skipFCntCheck"`
}

func (a *RestClient) GetChirpstackToken() (*ChirpstackToken, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetBody(ChirpstackCredentials{
			Email:    "admin",
			Password: "Helensburgh2508",
		}).
		SetResult(ChirpstackToken{}).
		Post("/api/internal/login"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ChirpstackToken), nil
}

func (a *RestClient) GetChirpstackDeviceProfileUUID(token string) ([]*ChirpstackDevProfile, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Accept", "application/json").
		SetResult([]*ChirpstackDevProfile{}).
		Post("/api/device-profiles?limit=50"))
	if err != nil {
		return nil, err
	}
	return resp.Result().([]*ChirpstackDevProfile), nil
}

func (a *RestClient) AddChirpstackDevice(chirpstackAppNum, modbusAddress int, deviceName, lorawanDeviceEUI, chirpstackDeviceProfileUUID string) (*interface{}, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetBody(ChirpstackDev{
			ApplicationID:     string(chirpstackAppNum),
			Description:       "Modbus Address: " + string(modbusAddress),
			DeviceEUI:         lorawanDeviceEUI,
			DeviceProfileID:   chirpstackDeviceProfileUUID,
			IsDisabled:        false,
			Name:              deviceName,
			ReferenceAltitude: 0,
			SkipFCntCheck:     true,
		}).
		SetResult(ChirpstackDev{}).
		Post("/api/devices"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*interface{}), nil
}
