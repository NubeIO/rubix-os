package csmodel

import "github.com/NubeIO/flow-framework/plugin/defaults"

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

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"true"`
	Options  bool   `json:"options" default:"false"`
	Default  *bool  `json:"default" default:"true"`
}

type NameNet struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"lorawan"`
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

type Interface struct {
	Type     string   `json:"type" default:"string"`
	Required bool     `json:"required" default:"true"`
	Options  []string `json:"options" default:"[]"`
	Default  string   `json:"default" default:""`
}

type SchemaNetwork struct {
	Enable      EnableStruct      `json:"enable"`
	Name        NameNet           `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"lorawan"`
	} `json:"plugin_name"`
}

type SchemaDevice struct {
	Enable         EnableStruct      `json:"enable"`
	Name           NameDev           `json:"name"`
	Description    DescriptionStruct `json:"description"`
	DeviceObjectId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"1"`
	} `json:"device_object_id"`
}

type SchemaPoint struct {
	Enable      EnableStruct      `json:"enable"`
	Name        NamePnt           `json:"name"`
	Description DescriptionStruct `json:"description"`
	ObjectType  struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"analog_input\",\"analog_value\",\"analog_output\",\"binary_input\",\"binary_value\",\"binary_output\"]"`
		Default  string   `json:"default" default:"analog_value"`
	} `json:"object_type"`
	ObjectId struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Default  int    `json:"default" default:"0"`
	} `json:"object_id"`
}
