package edgerest

import (
	"fmt"
)

type ServerPing struct {
	State string `json:"1_state"`
}

type UI struct {
	State string `json:"1_state"`
	IoNum string `json:"2_ioNum"`
	Gpio  string `json:"3_gpio"`
	Val   struct {
		UI1 struct {
			Val float64 `json:"val"`
		} `json:"UI1"`
		UI2 struct {
			Val float64 `json:"val"`
		} `json:"UI2"`
		UI3 struct {
			Val float64 `json:"val"`
		} `json:"UI3"`
		UI4 struct {
			Val float64 `json:"val"`
		} `json:"UI4"`
		UI5 struct {
			Val float64 `json:"val"`
		} `json:"UI5"`
		UI6 struct {
			Val float64 `json:"val"`
		} `json:"UI6"`
		UI7 struct {
			Val float64 `json:"val"`
		} `json:"UI7"`
	} `json:"4_val"`
	Msg      string `json:"5_msg"`
	MinRange struct {
		UI1 int `json:"UI1"`
		UI2 int `json:"UI2"`
		UI3 int `json:"UI3"`
		UI4 int `json:"UI4"`
		UI5 int `json:"UI5"`
		UI6 int `json:"UI6"`
		UI7 int `json:"UI7"`
	} `json:"6_min_range"`
	MaxRange struct {
		UI1 int `json:"UI1"`
		UI2 int `json:"UI2"`
		UI3 int `json:"UI3"`
		UI4 int `json:"UI4"`
		UI5 int `json:"UI5"`
		UI6 int `json:"UI6"`
		UI7 int `json:"UI7"`
	} `json:"7_max_range"`
}

// PingServer all points
func (a *RestClient) PingServer() (*ServerPing, error) {
	resp, err := a.client.R().
		SetResult(&ServerPing{}).
		Get("/")
	if err != nil {
		return nil, fmt.Errorf("error geting server %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ServerPing), nil
}

// GetUIs all ui points
func (a *RestClient) GetUIs() (*UI, error) {
	resp, err := a.client.R().
		SetResult([]UI{}).
		Get("/api/1.1/read/all/ui")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	fmt.Println(resp.String())
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*UI), nil
}
