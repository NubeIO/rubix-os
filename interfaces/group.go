package interfaces

import "github.com/NubeIO/rubix-os/schema/schema"

type GroupProperties struct {
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
}

func GetGroupProperties() *GroupProperties {
	m := &GroupProperties{}
	m.Name.Min = 0
	schema.Set(m)
	return m
}

type GroupSchema struct {
	Required   []string         `json:"required"`
	Properties *GroupProperties `json:"properties"`
}

func GetGroupSchema() *GroupSchema {
	m := &GroupSchema{
		Required:   []string{},
		Properties: GetGroupProperties(),
	}
	return m
}
