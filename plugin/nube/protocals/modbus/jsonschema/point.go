package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type PointSchema struct {
	UUID           schema.UUID           `json:"uuid"`
	Name           schema.Name           `json:"name"`
	Description    schema.Description    `json:"description"`
	Enable         schema.Enable         `json:"enable"`
	AddressId      schema.AddressId      `json:"address_id"`
	AddressLength  schema.AddressLength  `json:"address_length"`
	DataType       schema.DataType       `json:"data_type"`
	ObjectEncoding schema.ObjectEncoding `json:"object_encoding"`
}

func GetPointSchema() *PointSchema {
	m := &PointSchema{}
	schema.Set(m)
	return m
}
