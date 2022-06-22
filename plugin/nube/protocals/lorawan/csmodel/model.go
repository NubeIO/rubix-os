package csmodel

type Devices struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		DevEUI string `json:"devEUI"`
		Name   string `json:"name"`
	} `json:"result"`
}

type Device struct {
	Device struct {
		DevEUI string `json:"devEUI"`
		Name   string `json:"name"`
	} `json:"device"`
}
