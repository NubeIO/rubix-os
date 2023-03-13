package smodel

import (
	"github.com/NubeIO/flow-framework/plugin/defaults"
)

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"true"`
	Options  bool   `json:"options" default:"false"`
	Default  *bool  `json:"default" default:"true"`
}

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"3"`
	Max      int    `json:"max" default:"100"`
	Default  string `json:"default" default:"net"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
	Default  string `json:"default" default:"na"`
}

type Network struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"net"`
		DisplayName string `json:"display_name" default:"Network Name"`
	} `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"system"`
	} `json:"plugin_name"`
}

type Device struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"dev"`
		DisplayName string `json:"display_name" default:"Device Name"`
	} `json:"name"`
	Description       DescriptionStruct `json:"description"`
	AutoMappingEnable struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
		Default  *bool  `json:"default" default:"false"`
	} `json:"auto_mapping_enable"`
	AutoMappingFlowNetworkName struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[]"`
		Default  string   `json:"default" default:""`
	} `json:"auto_mapping_flow_network_name"`
}

type Point struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"pnt"`
		DisplayName string `json:"display_name" default:"Point Name"`
	} `json:"name"`
	Description DescriptionStruct `json:"description"`
	Fallback    struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"false"`
		Default     *float64 `json:"default" default:""`
		DisplayName string   `json:"display_name" default:"Fallback Value"`
		Nullable    bool     `json:"nullable" default:"true"`
	} `json:"fallback"`
	AutoMappingEnable struct {
		Type    string `json:"type" default:"bool"`
		Options bool   `json:"options" default:"false"`
		Default *bool  `json:"default" default:"true"`
	} `json:"auto_mapping_enable"`
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
