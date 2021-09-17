package pkgmodel

import (
	"github.com/NubeIO/null"
	"time"
)



type Devices struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
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
	} `json:"result"`
}

type DeviceProfiles struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id                string    `json:"id"`
		Name              string    `json:"name"`
		OrganizationID    string    `json:"organizationID"`
		NetworkServerID   string    `json:"networkServerID"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
		NetworkServerName string    `json:"networkServerName"`
	} `json:"result"`
}

type Applications struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id                 string `json:"id"`
		Name               string `json:"name"`
		Description        string `json:"description"`
		OrganizationID     string `json:"organizationID"`
		ServiceProfileID   string `json:"serviceProfileID"`
		ServiceProfileName string `json:"serviceProfileName"`
	} `json:"result"`
}


type ServiceProfiles struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id                string    `json:"id"`
		Name              string    `json:"name"`
		OrganizationID    string    `json:"organizationID"`
		NetworkServerID   string    `json:"networkServerID"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
		NetworkServerName string    `json:"networkServerName"`
	} `json:"result"`
}


type Users struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id         string      `json:"id"`
		Email      string      `json:"email"`
		SessionTTL int         `json:"sessionTTL"`
		IsAdmin    interface{} `json:"isAdmin"`
		IsActive   interface{} `json:"isActive"`
		CreatedAt  time.Time   `json:"createdAt"`
		UpdatedAt  time.Time   `json:"updatedAt"`
	} `json:"result"`
}


type Gateways struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id              string    `json:"id"`
		Name            string    `json:"name"`
		Description     string    `json:"description"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
		FirstSeenAt     time.Time `json:"firstSeenAt"`
		LastSeenAt      time.Time `json:"lastSeenAt"`
		OrganizationID  string    `json:"organizationID"`
		NetworkServerID string    `json:"networkServerID"`
		Location        struct {
			Latitude  int    `json:"latitude"`
			Longitude int    `json:"longitude"`
			Altitude  int    `json:"altitude"`
			Source    string `json:"source"`
			Accuracy  int    `json:"accuracy"`
		} `json:"location"`
	} `json:"result"`
}





type Server struct {
	Ip                string `json:"ip"`
	Port              int    `json:"port"`
	DeviceId          string `json:"device_id"`
	LocalObjName      string `json:"local_obj_name"`
	ModelName         string `json:"model_name"`
	VendorId          string `json:"vendor_id"`
	VendorName        string `json:"vendor_name"`
	EnableIpByNicName bool   `json:"enable_ip_by_nic_name"`
	IpByNicName       string `json:"ip_by_nic_name"` //eth0
}

// MqttPayload payload from the bacnet server
type MqttPayload struct {
	Value    null.Float
	Priority int
}









type BacnetPoint struct {
	ObjectType           string  `json:"object_type,omitempty"`
	ObjectName           string  `json:"object_name,omitempty"`
	Address              int     `json:"address,omitempty"`
	RelinquishDefault    float64 `json:"relinquish_default"`
	EventState           string  `json:"event_state,omitempty"`
	Units                string  `json:"units,omitempty"`
	Description          string  `json:"description,omitempty"`
	Enable               bool    `json:"enable,omitempty"`
	Fault                bool    `json:"fault,omitempty"`
	DataRound            float64 `json:"data_round,omitempty"`
	DataOffset           float64 `json:"data_offset,omitempty"`
	UseNextAvailableAddr bool    `json:"use_next_available_address,omitempty"`
	COV                  float32 `json:"cov,omitempty"`
	Priority             `json:"priority_array_write,omitempty"`
}

type Priority struct {
	P1  null.Float `json:"_1,omitempty"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	P2  null.Float `json:"_2,omitempty"`
	P3  null.Float `json:"_3,omitempty"`
	P4  null.Float `json:"_4,omitempty"`
	P5  null.Float `json:"_5,omitempty"`
	P6  null.Float `json:"_6,omitempty"`
	P7  null.Float `json:"_7,omitempty"`
	P8  null.Float `json:"_8,omitempty"`
	P9  null.Float `json:"_9,omitempty"`
	P10 null.Float `json:"_10,omitempty"`
	P11 null.Float `json:"_11,omitempty"`
	P12 null.Float `json:"_12,omitempty"`
	P13 null.Float `json:"_13,omitempty"`
	P14 null.Float `json:"_14,omitempty"`
	P15 null.Float `json:"_15,omitempty"`
	P16 null.Float `json:"_16"` //removed and added to the point to save one DB write

}
