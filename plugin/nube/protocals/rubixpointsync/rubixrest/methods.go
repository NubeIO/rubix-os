package rubixrest

import (
	"github.com/NubeIO/flow-framework/nresty"
)

type RubixNet struct {
	UUID          string      `json:"uuid"`
	CreatedOn     string      `json:"created_on"`
	UpdatedOn     string      `json:"updated_on"`
	Name          string      `json:"name"`
	Enable        bool        `json:"enable"`
	Fault         bool        `json:"fault"`
	HistoryEnable bool        `json:"history_enable"`
	Devices       []*RubixDev `json:"devices"`
}

type RubixDev struct {
	UUID          string      `json:"uuid"`
	CreatedOn     string      `json:"created_on"`
	UpdatedOn     string      `json:"updated_on"`
	NetworkUUID   string      `json:"network_uuid"`
	Name          string      `json:"name"`
	Enable        bool        `json:"enable"`
	Fault         bool        `json:"fault"`
	HistoryEnable bool        `json:"history_enable"`
	Points        []*RubixPnt `json:"points"`
}

type RubixPnt struct {
	UUID          string          `json:"uuid"`
	CreatedOn     string          `json:"created_on"`
	UpdatedOn     string          `json:"updated_on"`
	DeviceUUID    string          `json:"device_uuid"`
	Name          string          `json:"name"`
	Enable        bool            `json:"enable"`
	HistoryEnable bool            `json:"history_enable"`
	PointStore    RubixPointStore `json:"point_store"`
}

type RubixPointStore struct {
	Value         float64 `json:"value"`
	ValueOriginal float64 `json:"value_original"`
	Fault         bool    `json:"fault"`
}

type RubixPointWrite struct {
	PriorityArray PriorityArrayWrite `json:"priority_array_write"`
	Enable        bool               `json:"enable"`
	Writable      bool               `json:"writable"`
	// DisableMQTT					bool    					`json:"disable_mqtt"`
}

type PriorityArrayWrite struct {
	P1  *float64 `json:"_1"`
	P2  *float64 `json:"_2"`
	P3  *float64 `json:"_3"`
	P4  *float64 `json:"_4"`
	P5  *float64 `json:"_5"`
	P6  *float64 `json:"_6"`
	P7  *float64 `json:"_7"`
	P8  *float64 `json:"_8"`
	P9  *float64 `json:"_9"`
	P10 *float64 `json:"_10"`
	P11 *float64 `json:"_11"`
	P12 *float64 `json:"_12"`
	P13 *float64 `json:"_13"`
	P14 *float64 `json:"_14"`
	P15 *float64 `json:"_15"`
	P16 *float64 `json:"_16"`
}

func (a *RestClient) GetAllPoints() (*[]RubixNet, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetResult([]RubixNet{}).
		Get("/api/generic/networks?with_children=true&points=true"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]RubixNet), nil
}

func (a *RestClient) CreateNewRubixPoint(pointName, deviceUUID string) (*RubixPnt, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(RubixPnt{
			Name:          pointName,
			Enable:        true,
			DeviceUUID:    deviceUUID,
			HistoryEnable: true,
		}).
		SetResult(RubixPnt{}).
		Post("/api/generic/points"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*RubixPnt), nil
}

func (a *RestClient) WriteRubixPointByPathNames(networkName, deviceName, pointName string, writeValue *float64) (*RubixPnt, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetBody(RubixPointWrite{
			PriorityArray: PriorityArrayWrite{
				P3: writeValue,
			},
		}).
		SetResult(RubixPnt{}).
		Patch(`/api/generic/points/name/` + networkName + `/` + deviceName + `/` + pointName))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*RubixPnt), nil
}

func (a *RestClient) WriteRubixPointByUUID(pointUUID string, writeValue *float64) (*RubixPnt, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetBody(RubixPointWrite{
			PriorityArray: PriorityArrayWrite{
				P3: writeValue,
			},
		}).
		SetResult(RubixPnt{}).
		Patch(`/api/generic/points/uuid/` + pointUUID))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*RubixPnt), nil
}

func (a *RestClient) CreateNewRubixDevice(deviceName, networkUUID string) (*RubixDev, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(RubixDev{
			Name:          deviceName,
			Enable:        true,
			NetworkUUID:   networkUUID,
			HistoryEnable: true,
		}).
		SetResult(RubixDev{}).
		Post("/api/generic/devices"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*RubixDev), nil
}

func (a *RestClient) CreateNewRubixNetwork(netName string) (*RubixNet, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(RubixNet{
			Name:          netName,
			Enable:        true,
			HistoryEnable: true,
		}).
		SetResult(RubixNet{}).
		Post("/api/generic/networks"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*RubixNet), nil
}
