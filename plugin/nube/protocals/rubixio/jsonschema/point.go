package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type IoNumber struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"io number"`
	Options  []string `json:"enum" default:"[\"UI1\",\"UI2\",\"UI3\",\"UI4\",\"UI5\",\"UI6\",\"UI7\",\"UI8\",\"UO1\",\"UO2\",\"UO3\",\"UO4\",\"UO5\",\"UO6\",\"DO1\",\"DO2\"]"`
	EnumName []string `json:"enumNames" default:"[\"UI1\",\"UI2\",\"UI3\",\"UI4\",\"UI5\",\"UI6\",\"UI7\",\"UI8\",\"UO1\",\"UO2\",\"UO3\",\"UO4\",\"UO5\",\"UO6\",\"DO1\",\"DO2\"]"`
	Default  string   `json:"default" default:"UI1"`
}

type IoType struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"io type"`
	Options  []string `json:"enum" default:"[\"digital\",\"voltage_dc\",\"thermistor_10k_type_2\",\"current\",\"raw\"]"`
	EnumName []string `json:"enumNames" default:"[\"digital\",\"voltage dc\",\"thermistor 10k-type-2\",\"current\",\"raw\"]"`
	Default  string   `json:"default" default:"digital"`
}

type PointSchema struct {
	UUID        schema.UUID        `json:"uuid"`
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
	Enable      schema.Enable      `json:"enable"`
	IoNumber    IoNumber           `json:"io_number"`
	IoType      IoType             `json:"io_type"`
}

func GetPointSchema() *DeviceSchema {
	m := &DeviceSchema{}
	schema.Set(m)
	return m
}
