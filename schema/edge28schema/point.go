package edge28schema

import "github.com/NubeIO/rubix-os/schema/schema"

// Options  []string `json:"options"

type IoNumber struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"IO Number"`
	Options  []string `json:"enum" default:"[\"UI1\",\"UI2\",\"UI3\",\"UI4\",\"UI5\",\"UI6\",\"UI7\",\"DI1\",\"DI2\",\"DI3\",\"DI4\",\"DI5\",\"DI6\",\"DI7\",\"R1\",\"R2\",\"DO1\",\"DO2\",\"DO3\",\"DO4\",\"DO5\",\"UO1\",\"UO2\",\"UO3\",\"UO4\",\"UO5\",\"UO6\",\"UO7\"]"`
	EnumName []string `json:"enumNames" default:"[\"UI1\",\"UI2\",\"UI3\",\"UI4\",\"UI5\",\"UI6\",\"UI7\",\"DI1\",\"DI2\",\"DI3\",\"DI4\",\"DI5\",\"DI6\",\"DI7\",\"R1\",\"R2\",\"DO1\",\"DO2\",\"DO3\",\"DO4\",\"DO5\",\"UO1\",\"UO2\",\"UO3\",\"UO4\",\"UO5\",\"UO6\",\"UO7\"]"`
	Default  string   `json:"default" default:"UI1"`
}

type IoType struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"IO Type"`
	Options  []string `json:"enum" default:"[\"digital\",\"voltage_dc\",\"thermistor_10k_type_2\",\"current\",\"raw\"]"`
	EnumName []string `json:"enumNames" default:"[\"digital\",\"voltage dc\",\"thermistor 10k-type-2\",\"current\",\"raw\"]"`
	Default  string   `json:"default" default:"thermistor_10k_type_2"`
}

type ObjectType struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Object Type"`
	Options  []string `json:"enum" default:"[\"analog_input\"\"analog_output\",\"binary_input\",\"binary_value\",\"binary_output\"]"`
	EnumName []string `json:"enumNames" default:"[\"analog input\"\"analog output\",\"binary input\",\"binary value\",\"binary output\"]"`
	Default  string   `json:"default" default:"analog_input"`
}

type PointSchema struct {
	UUID        schema.UUID        `json:"uuid"`
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
	Enable      schema.Enable      `json:"enable"`

	IoNumber IoNumber `json:"io_number"`
	IoType   IoType   `json:"io_type"`
	// ObjectType  schema.ObjectType  `json:"object_type"`

	MultiplicationFactor schema.MultiplicationFactor `json:"multiplication_factor"`
	ScaleEnable          schema.ScaleEnable          `json:"scale_enable"`
	ScaleInMin           schema.ScaleInMin           `json:"scale_in_min"`
	ScaleInMax           schema.ScaleInMax           `json:"scale_in_max"`
	ScaleOutMin          schema.ScaleOutMin          `json:"scale_out_min"`
	ScaleOutMax          schema.ScaleOutMax          `json:"scale_out_max"`
	Offset               schema.Offset               `json:"offset"`
	Decimal              schema.Decimal              `json:"decimal"`
	Fallback             schema.Fallback             `json:"fallback"`

	HistoryEnable   schema.HistoryEnableDefaultTrue `json:"history_enable"`
	HistoryType     schema.HistoryType              `json:"history_type"`
	HistoryInterval schema.HistoryInterval          `json:"history_interval"`
}

func GetPointSchema() *PointSchema {
	m := &PointSchema{}
	schema.Set(m)
	return m
}
