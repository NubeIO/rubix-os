package edge_model

import "github.com/NubeIO/flow-framework/plugin/defaults"

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"3"`
	Max      int    `json:"max" default:"20"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type Network struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"edge28"`
	} `json:"plugin_name"`
}

type Device struct {
	Name        NameStruct        `json:"name"`
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
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	IoNumber        struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"UI1\",\"UI2\",\"UI3\",\"temp\",\"pulse\",\"motion\",\"light\",\"voltage\"]"`
	} `json:"io_number"`
	IoType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"THERMISTOR_10K_TYPE2\",\"UI2\",\"UI3\",\"temp\",\"pulse\",\"motion\",\"light\",\"voltage\"]"`
	} `json:"io_type"`
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
