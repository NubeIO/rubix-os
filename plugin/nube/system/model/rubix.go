package system_model

import "github.com/NubeIO/flow-framework/plugin/defaults"

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"3"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"a"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
	Default  string `json:"default" default:"na"`
}

type Network struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"system"`
	} `json:"plugin_name"`
}

type Device struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
}

type Point struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	//ObjectType  struct {
	//	Type     string   `json:"type" default:"array"`
	//	Required bool     `json:"required" default:"true"`
	//	Options  []string `json:"options" default:"[\"analogValue\",\"analogOutput\"]"`
	//	Default  string   `json:"default" default:"analogValue"`
	//} `json:"object_type"`
	//AddressID struct {
	//	Type     string `json:"type" default:"int"`
	//	Required bool   `json:"required" default:"true"`
	//	Default  int    `json:"default" default:"1"`
	//} `json:"address_id"`
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
