package csmodel

import "time"

type Device struct {
	DevEUI                              string    `json:"devEUI"`
	Name                                string    `json:"name"`
	ApplicationID                       string    `json:"applicationID"`
	Description                         string    `json:"description"`
	DeviceProfileID                     string    `json:"deviceProfileID"`
	DeviceProfileName                   string    `json:"deviceProfileName"`
	DeviceStatusBattery                 int       `json:"deviceStatusBattery"`
	DeviceStatusMargin                  int       `json:"deviceStatusMargin"`
	DeviceStatusExternalPowerSource     bool      `json:"deviceStatusExternalPowerSource"`
	DeviceStatusBatteryLevelUnavailable bool      `json:"deviceStatusBatteryLevelUnavailable"`
	DeviceStatusBatteryLevel            int       `json:"deviceStatusBatteryLevel"`
	LastSeenAt                          time.Time `json:"lastSeenAt"`
	LastSeenAtTime                      string    `json:"lastSeenAtTime"`
	LastSeenAtReadable                  string    `json:"lastSeenAtReadable"`
}

type Devices struct {
	Result []Device `json:"result"`
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
	FPort  int                    `json:"fPort"`
	Object map[string]interface{} `json:"object"`
}
