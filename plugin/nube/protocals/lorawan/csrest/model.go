package csrest

import "time"

type DevicesResult struct {
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

type DeviceSingle struct {
	Device              *DeviceBody `json:"device"`
	LastSeenAt          string      `json:"lastSeenAt"`
	DeviceStatusBattery int         `json:"deviceStatusBattery"`
	DeviceStatusMargin  int         `json:"deviceStatusMargin"`
	Location            interface{} `json:"location"`
}

// Devices GET
type Devices struct {
	TotalCount string           `json:"totalCount"`
	Result     []*DevicesResult `json:"result"`
}

type DeviceBody struct {
	DevEUI            string `json:"devEUI"`
	Name              string `json:"name"`
	ApplicationID     string `json:"applicationID"`
	Description       string `json:"description"`
	DeviceProfileID   string `json:"deviceProfileID"`
	SkipFCntCheck     bool   `json:"skipFCntCheck"`
	ReferenceAltitude int    `json:"referenceAltitude"`
	Variables         struct {
	} `json:"variables"`
	Tags struct {
	} `json:"tags"`
	IsDisabled bool `json:"isDisabled"`
}

type DeviceKeys struct {
	DevEUI    string `json:"devEUI"`
	NwkKey    string `json:"nwkKey"`
	AppKey    string `json:"appKey"`
	GenAppKey string `json:"genAppKey"`
}

type DeviceKey struct {
	Keys DeviceKeys `json:"deviceKeys"`
}

type DeviceActivation struct {
	DeviceActivation struct {
		DevAddr     string `json:"devAddr"`
		NwkSEncKey  string `json:"nwkSEncKey"`
		AppSKey     string `json:"appSKey"`
		DevEUI      string `json:"devEUI"`
		FNwkSIntKey string `json:"fNwkSIntKey"`
		SNwkSIntKey string `json:"sNwkSIntKey"`
	} `json:"deviceActivation"`
}

// DeviceActive used to POST active
type DeviceActive struct {
	DActive struct {
		DevEUI      string `json:"devEUI"`
		DevAddr     string `json:"devAddr"`
		AppSKey     string `json:"appSKey"`
		NwkSEncKey  string `json:"nwkSEncKey"`
		SNwkSIntKey string `json:"sNwkSIntKey"`
		FNwkSIntKey string `json:"fNwkSIntKey"`
		FCntUp      int    `json:"fCntUp"`
		NFCntDown   int    `json:"nFCntDown"`
		AFCntDown   int    `json:"aFCntDown"`
	} `json:"deviceActivation"`
}

type DeviceProfileResult struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	OrganizationID    string    `json:"organizationID"`
	NetworkServerID   string    `json:"networkServerID"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	NetworkServerName string    `json:"networkServerName"`
}

type DeviceProfiles struct {
	TotalCount string                 `json:"totalCount"`
	Result     []*DeviceProfileResult `json:"result"`
}

type ApplicationResult struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	OrganizationID     string `json:"organizationID"`
	ServiceProfileID   string `json:"serviceProfileID"`
	ServiceProfileName string `json:"serviceProfileName"`
}

type Applications struct {
	TotalCount string               `json:"totalCount"`
	Result     []*ApplicationResult `json:"result"`
}

type MQTTUplink struct {
	DeviceName string `json:"deviceName"`
	DevEUI     string `json:"devEUI"`
	RxInfo     []struct {
		Rssi    int     `json:"rssi"`
		LoRaSNR float64 `json:"loRaSNR"`
	} `json:"rxInfo"`
	FPort  int                    `json:"fPort"`
	Object map[string]interface{} `json:"object"`
}

type MQTTError struct {
	ApplicationId   string `json:"applicationID"`
	ApplicationName string `json:"applicationName"`
	DeviceName      string `json:"deviceName"`
	DevEUI          string `json:"devEUI"`
	Type            string `json:"type"`
	Error           string `json:"error"`
}
