package pkgmodel

import (
	"github.com/NubeDev/flow-framework/model"
)

type ServerPing struct {
	Version        string `json:"version"`
	UpTimeDate     string `json:"up_time_date"`
	UpMin          string `json:"up_min"`
	UpHour         string `json:"up_hour"`
	DeploymentMode string `json:"deployment_mode"`
	Bacnet         struct {
		Enabled                    bool    `json:"enabled"`
		Status                     bool    `json:"status"`
		UsePreSetEthernetInterface bool    `json:"use_pre_set_ethernet_interface"`
		PreSetEthernetInterface    string  `json:"pre_set_ethernet_interface"`
		DefaultPointCov            float32 `json:"default_point_cov"`
	} `json:"bacnet"`
	Mqtt struct {
		Enabled bool `json:"enabled"`
		Status  bool `json:"status"`
	} `json:"mqtt"`
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
	Value    *float64
	Priority int
}

type BacnetPoint struct {
	AddressUUID          string          `json:"address_uuid,omitempty"`
	ObjectType           string          `json:"object_type,omitempty"`
	ObjectName           string          `json:"object_name,omitempty"`
	Address              int             `json:"address,omitempty"`
	RelinquishDefault    float64         `json:"relinquish_default"`
	EventState           string          `json:"event_state,omitempty"`
	Units                string          `json:"units,omitempty"`
	Description          string          `json:"description,omitempty"`
	Enable               bool            `json:"enable,omitempty"`
	Fault                bool            `json:"fault,omitempty"`
	DataRound            float64         `json:"data_round,omitempty"`
	DataOffset           float64         `json:"data_offset,omitempty"`
	UseNextAvailableAddr bool            `json:"use_next_available_address,omitempty"`
	COV                  float32         `json:"cov,omitempty"`
	Priority             *model.Priority `json:"priority_array_write,omitempty"`
}
