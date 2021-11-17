package model

import (
	"github.com/NubeIO/flow-framework/plugin/defaults"
)

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

type TransportTypeStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Default  string `json:"default" default:"serial"`
}

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"false"`
}

type Network struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"lora"`
	} `json:"plugin_name"`
	TransportType TransportTypeStruct `json:"transport_type"`
	SerialPort    struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Min      int    `json:"min" default:"3"`
		Max      int    `json:"max" default:"20"`
	} `json:"serial_port"`
	BaudRate struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Default  int    `json:"default" default:"9600"`
	} `json:"baud_rate"`
}

type Device struct {
	Name          NameStruct          `json:"name"`
	Description   DescriptionStruct   `json:"description"`
	Enable        EnableStruct        `json:"enable"`
	TransportType TransportTypeStruct `json:"transport_type"`
}

type Point struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	Address     struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Min      int    `json:"min" default:"8"`
		Max      int    `json:"max" default:"8"`
	} `json:"address"`
	IoId struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Min      int    `json:"min" default:"3"`
		Max      int    `json:"max" default:"20"`
	} `json:"io_id"`
	ThingClass struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"point\"]"`
	} `json:"thing_class"`
	ThingType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"point\"]"`
	} `json:"thing_type"`
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
