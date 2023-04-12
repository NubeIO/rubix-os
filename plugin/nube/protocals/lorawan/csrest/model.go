package csrest

import "time"

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

type OrganizationsResult struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	DisplayName     string    `json:"displayName"`
	CanHaveGateways bool      `json:"canHaveGateways"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Organizations struct {
	TotalCount string                 `json:"totalCount"`
	Result     []*OrganizationsResult `json:"result"`
}

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

// Devices GET
type Devices struct {
	TotalCount string           `json:"totalCount"`
	Result     []*DevicesResult `json:"result"`
}

type T struct {
	Device struct {
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
	} `json:"device"`
}

type DeviceBody struct {
	DevEUI          string `json:"devEUI"`
	Name            string `json:"name"`
	ApplicationID   string `json:"applicationID"`
	Description     string `json:"description"`
	DeviceProfileID string `json:"deviceProfileID"`
	IsDisabled      bool   `json:"isDisabled"`
}

type Device struct {
	Device              *DeviceBody `json:"device"`
	LastSeenAt          interface{} `json:"lastSeenAt"`
	DeviceStatusBattery int         `json:"deviceStatusBattery"`
	DeviceStatusMargin  int         `json:"deviceStatusMargin"`
	Location            interface{} `json:"location"`
}

// type DeviceAdd struct {
// 	Device *DeviceBody `json:"device"`
// }

type DeviceAdd struct {
	Device struct {
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
	} `json:"device"`
}

type DeviceKey struct {
	DeviceKeys struct {
		DevEUI    string `json:"devEUI"`
		NwkKey    string `json:"nwkKey"`
		AppKey    string `json:"appKey"`
		GenAppKey string `json:"genAppKey"`
	} `json:"deviceKeys"`
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

type DeviceProfilesResult struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	OrganizationID    string    `json:"organizationID"`
	NetworkServerID   string    `json:"networkServerID"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	NetworkServerName string    `json:"networkServerName"`
}

type DeviceProfiles struct {
	TotalCount string                  `json:"totalCount"`
	Result     []*DeviceProfilesResult `json:"result"`
}

type ApplicationsResult struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	OrganizationID     string `json:"organizationID"`
	ServiceProfileID   string `json:"serviceProfileID"`
	ServiceProfileName string `json:"serviceProfileName"`
}

type Applications struct {
	TotalCount string                `json:"totalCount"`
	Result     []*ApplicationsResult `json:"result"`
}

type ServiceProfilesResult struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	OrganizationID    string    `json:"organizationID"`
	NetworkServerID   string    `json:"networkServerID"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	NetworkServerName string    `json:"networkServerName"`
}

type ServiceProfiles struct {
	TotalCount string                   `json:"totalCount"`
	Result     []*ServiceProfilesResult `json:"result"`
}

type GatewayProfilesResult struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	NetworkServerID   string    `json:"networkServerID"`
	NetworkServerName string    `json:"networkServerName"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type GatewayProfiles struct {
	TotalCount string                   `json:"totalCount"`
	Result     []*GatewayProfilesResult `json:"result"`
}

type UsersResult struct {
	Id         string      `json:"id"`
	Email      string      `json:"email"`
	SessionTTL int         `json:"sessionTTL"`
	IsAdmin    interface{} `json:"isAdmin"`
	IsActive   interface{} `json:"isActive"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

type Users struct {
	TotalCount string         `json:"totalCount"`
	Result     []*UsersResult `json:"result"`
}

type GatewaysResult struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	FirstSeenAt       time.Time `json:"firstSeenAt"`
	LastSeenAt        time.Time `json:"lastSeenAt"`
	FirstSeenAtString string    `json:"firstSeenAtString"`
	LastSeenAtString  string    `json:"lastSeenAtString"`
	OrganizationID    string    `json:"organizationID"`
	NetworkServerID   string    `json:"networkServerID"`
}

type Gateways struct {
	TotalCount string            `json:"totalCount"`
	Result     []*GatewaysResult `json:"result"`
}
