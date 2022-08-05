package edgemodel

import "github.com/NubeIO/flow-framework/plugin/defaults"

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"true"`
	Options  bool   `json:"options" default:"false"`
	Default  *bool  `json:"default" default:"true"`
}

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"edge"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type Network struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"edge28_net"`
		DisplayName string `json:"display_name" default:"Network Name"`
	} `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"edge28"`
	} `json:"plugin_name"`
	AutoMappingNetworksSelection struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"disable\",\"edge28-to-bacnetserver\"]"`
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
}

type Device struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"edge28_dev"`
		DisplayName string `json:"display_name" default:"Device Name"`
	} `json:"name"`
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
		Options  int    `json:"options" default:"5000"`
		Default  int    `json:"default" default:"5000"`
	} `json:"port"`
}

type Point struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"edge28_pnt"`
		DisplayName string `json:"display_name" default:"Point Name"`
	} `json:"name"`
	Description DescriptionStruct `json:"description"`
	IoNumber    struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"UI1\",\"UI2\",\"UI3\",\"UI4\",\"UI5\",\"UI6\",\"UI7\",\"DI1\",\"DI2\",\"DI3\",\"DI4\",\"DI5\",\"DI6\",\"DI7\",\"R1\",\"R2\",\"DO1\",\"DO2\",\"DO3\",\"DO4\",\"DO5\",\"UO1\",\"UO2\",\"UO3\",\"UO4\",\"UO5\",\"UO6\",\"UO7\"]"`
	} `json:"io_number"`
	IoType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"THERMISTOR_10K_TYPE2\",\"DIGITAL\",\"PERCENT\",\"0-10VDC\",\"4-20mA\",\"RESISTANCE\"]"`
	} `json:"io_type"`
	Fallback struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"false"`
		Default     *float64 `json:"default" default:""`
		DisplayName string   `json:"display_name" default:"Fallback Value"`
		Nullable    bool     `json:"nullable" default:"true"`
	} `json:"fallback"`
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

type ServerPing struct {
	State string `json:"1_state"`
}
