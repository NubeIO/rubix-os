package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type PointSchema struct {
	UUID          schema.UUID        `json:"uuid"`
	Name          schema.Name        `json:"name"`
	Description   schema.Description `json:"description"`
	Enable        schema.Enable      `json:"enable"`
	ObjectId      schema.ObjectId    `json:"object_id"`
	ObjectType    schema.ObjectType  `json:"object_type"`
	WriteMode     schema.ObjectType  `json:"write_mode"`
	WritePriority schema.ObjectType  `json:"write_priority"`
}

func GetPointSchema() *PointSchema {
	m := &PointSchema{}
	schema.Set(m)
	return m
}
