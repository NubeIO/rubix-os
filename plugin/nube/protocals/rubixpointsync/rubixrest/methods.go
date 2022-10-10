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

func (a *RestClient) GetAllPoints() (*[]RubixNet, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetResult([]RubixNet{}).
		Get("/api/generic/networks?with_children=true&points=true"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]RubixNet), nil
}
