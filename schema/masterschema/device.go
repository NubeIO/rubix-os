package masterschema

import "github.com/NubeIO/rubix-os/schema/schema"

type DeviceSchema struct {
	UUID           schema.UUID                     `json:"uuid"`
	Name           schema.Name                     `json:"name"`
	Description    schema.Description              `json:"description"`
	Enable         schema.Enable                   `json:"enable"`
	Host           schema.Host                     `json:"host"`
	Port           schema.Port                     `json:"port"`
	DeviceObjectId schema.DeviceObjectId           `json:"device_object_id"`
	NetworkNumber  schema.NetworkNumber            `json:"network_number"`
	DeviceMac      schema.DeviceMac                `json:"device_mac"`
	MaxADPU        MaxADPU                         `json:"max_adpu"`
	Segmentation   Segmentation                    `json:"segmentation"`
	FastPollRate   schema.FastPollRate             `json:"fast_poll_rate"`
	NormalPollRate schema.NormalPollRate           `json:"normal_poll_rate"`
	SlowPollRate   schema.SlowPollRate             `json:"slow_poll_rate"`
	HistoryEnable  schema.HistoryEnableDefaultTrue `json:"history_enable"`
}
type Segmentation struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"device segmentation"`
	Options  []string `json:"enum" default:"[]"`
	EnumName []string `json:"enumNames" default:"[]"`
	Default  string   `json:"default" default:"no_segmentation"`
}

type MaxADPU struct {
	Type     string   `json:"type" default:"number"`
	Title    string   `json:"title" default:"device max-adpu"`
	Options  []int    `json:"enum" default:"[]"`
	EnumName []string `json:"enumNames" default:"[]"`
	Default  int      `json:"default" default:"1476"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	m.Segmentation.Options = []string{"segmentation_both", "no_segmentation", "segmentation_transmit", "segmentation_receive"}
	m.Segmentation.EnumName = []string{"segmentation-both", "no-segmentation", "segmentation-transmit", "segmentation-receive"}
	m.MaxADPU.Options = []int{50, 128, 206, 480, 1024, 1476}
	m.MaxADPU.EnumName = []string{"50", "128", "206", "480", "1024", "1476"}
	m.MaxADPU.Default = 1476
	m.Port.Default = 47808
	m.Host.Default = "192.168.15.10"
	schema.Set(m)
	return m
}
