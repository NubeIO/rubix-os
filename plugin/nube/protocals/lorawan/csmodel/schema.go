package csmodel

import "github.com/NubeIO/rubix-os/plugin/defaults"

func GetNetworkSchema() *SchemaNetwork {
	m := &SchemaNetwork{}
	defaults.Set(m)
	return m
}

func GetDeviceSchema() *SchemaDevice {
	device := &SchemaDevice{}
	defaults.Set(device)
	return device
}

func GetPointSchema() *SchemaPoint {
	point := &SchemaPoint{}
	defaults.Set(point)
	return point
}

type NameNet struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"100"`
	Default  string `json:"default" default:"lorawan"`
}

type NameDev struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"100"`
	Default  string `json:"default" default:"device"`
}

type NamePnt struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"100"`
	Default  string `json:"default" default:"point"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type Interface struct {
	Type     string   `json:"type" default:"string"`
	Required bool     `json:"required" default:"true"`
	Options  []string `json:"options" default:"[]"`
	Default  string   `json:"default" default:""`
}

type SchemaNetwork struct {
	Name        NameNet           `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"lorawan"`
	} `json:"plugin_name"`
}

type SchemaDevice struct {
	Name        NameDev           `json:"name"`
	Description DescriptionStruct `json:"description"`
}

type SchemaPoint struct {
	Name        NamePnt           `json:"name"`
	Description DescriptionStruct `json:"description"`
}
