package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type ObjectType struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"object type"`
	Options  []string `json:"enum" default:"[\"analog_input\",\"analog_value\",\"analog_output\",\"binary_input\",\"binary_value\",\"binary_output\"]"`
	EnumName []string `json:"enumNames" default:"[\"analog input\",\"analog value\",\"analog output\",\"binary input\",\"binary value\",\"binary output\"]"`
	Default  string   `json:"default" default:"analog_input"`
}

type PointSchema struct {
	UUID        schema.UUID        `json:"uuid"`
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
	Enable      schema.Enable      `json:"enable"`
	ObjectId    schema.ObjectId    `json:"object_id"`
	ObjectType  ObjectType         `json:"object_type"`
	//WriteMode     schema.ObjectType  `json:"write_mode"`
	//WritePriority schema.ObjectType  `json:"write_priority"`
}

func GetPointSchema() *PointSchema {
	m := &PointSchema{}
	schema.Set(m)
	return m
}
