package interfaces

import "github.com/NubeIO/lib-schema/schema"

type LocationProperties struct {
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
}

func GetLocationProperties() *LocationProperties {
	m := &LocationProperties{}
	m.Name.Min = 0
	schema.Set(m)
	return m
}

type LocationSchema struct {
	Required   []string            `json:"required"`
	Properties *LocationProperties `json:"properties"`
}

func GetLocationSchema() *LocationSchema {
	m := &LocationSchema{
		Required:   []string{},
		Properties: GetLocationProperties(),
	}
	return m
}
