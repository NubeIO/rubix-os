package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type PointSchema struct {
	UUID          schema.UUID          `json:"uuid"`
	Name          schema.Name          `json:"name"`
	Description   schema.Description   `json:"description"`
	Enable        schema.Enable        `json:"enable"`
	ObjectId      schema.ObjectId      `json:"object_id"`
	ObjectType    schema.ObjectType    `json:"object_type"`
	WriteMode     schema.WriteMode     `json:"write_mode"`
	WritePriority schema.WritePriority `json:"write_priority"`
	PollPriority  schema.PollPriority  `json:"poll_priority"`
	PollRate      schema.PollRate      `json:"poll_rate"`
}

func GetPointSchema() *PointSchema {
	m := &PointSchema{}
	schema.Set(m)
	return m
}
