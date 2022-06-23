package csmodel

type Device struct {
	DevEUI      string `json:"devEUI"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Devices struct {
	TotalCount string   `json:"totalCount"`
	Result     []Device `json:"result"`
}

type DeviceAll struct {
	Device Device `json:"device"`
}

type BaseUplink struct {
	DeviceName string `json:"deviceName"`
	DevEUI     string `json:"devEUI"`
	RxInfo     []struct {
		Rssi    int     `json:"rssi"`
		LoRaSNR float64 `json:"loRaSNR"`
	} `json:"rxInfo"`
	FPort int    `json:"fPort"`
	Data  string `json:"data"`
}
