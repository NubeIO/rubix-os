package master

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeIO/flow-framework/plugin/defaults"
	"github.com/gin-gonic/gin"
)

type NameNet struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"bacnet-master"`
}

type NameDev struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"device"`
}

type NamePnt struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"point"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"true"`
}

type Network struct {
	Name        NameNet           `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"bacnetmaster"`
	} `json:"plugin_name"`
	NetworkInterface struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"200"`
		Default  string `json:"default" default:"eth0"`
	} `json:"network_interface"`
	AutoMappingNetworksSelection struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"disable\",\"self-mapping\",\"rubix-io-to-bacnetserver\"]"`
		Default  string   `json:"default" default:""`
	} `json:"auto_mapping_networks_selection"`
	AutoMappingFlowNetworkName struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"200"`
		Default  string `json:"default" default:"local"`
	} `json:"auto_mapping_flow_network_name"`
	AutoMappingFlowNetworkUUID struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"200"`
		Default  string `json:"default" default:""`
	} `json:"auto_mapping_flow_network_uuid"`
	Enable EnableStruct `json:"enable"`
}

type Device struct {
	Name        NameDev           `json:"name"`
	Description DescriptionStruct `json:"description"`
	Host        struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Options  string `json:"options" default:"192.168.15.10"`
		Default  string `json:"default" default:"192.168.15.10"`
	} `json:"host"`
	Port struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Options  int    `json:"options" default:"47808"`
		Default  int    `json:"default" default:"47808"`
	} `json:"port"`
	DeviceObjectId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"1"`
	} `json:"device_object_id"`
	NetworkNumber struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"0"`
	} `json:"network_number"`
	DeviceMac struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"255"`
		Default  int    `json:"default" default:"0"`
	} `json:"device_mac"`
	MaxADPU struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"true"`
		Options  []int  `json:"options" default:"[50, 128, 206, 480, 1024, 1476]"`
		Default  int    `json:"default" default:"1024"`
	} `json:"max_adpu"`
	Segmentation struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"segmentation_both\",\"no_segmentation\",\"segmentation_transmit\",\"segmentation_receive\"]"`
		Default  string   `json:"default" default:"no_segmentation"`
	} `json:"segmentation"`
}

type Point struct {
	Name        NamePnt           `json:"name"`
	Description DescriptionStruct `json:"description"`
	ObjectType  struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"analog_input\",\"analog_value\",\"analog_output\",\"binary_input\",\"binary_value\",\"binary_output\",\"multi_state_input\",\"multi_state_value\",\"multi_state_output\"]"`
		Default  string   `json:"default" default:"analog_output"`
	} `json:"object_type"`
	ObjectId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"1"`
	} `json:"object_id"`
	WriteMode struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"read_only\",\"write_only\"]"`
		Default  string   `json:"default" default:"read_only"`
	} `json:"write_mode"`
	WritePriority struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16]"`
		Default  int    `json:"default" default:"16"`
	} `json:"write_priority"`
}

func GetNetworkSchema() *Network {
	network := &Network{}
	defaults.Set(network)
	return network
}

func GetDeviceSchema() *Device {
	device := &Device{}
	defaults.Set(device)
	return device
}

func GetPointSchema() *Point {
	point := &Point{}
	defaults.Set(point)
	return point
}

func BodyWhoIs(ctx *gin.Context) (dto *bacnet.WhoIsOpts, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}
