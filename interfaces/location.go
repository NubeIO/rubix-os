package interfaces

import "github.com/NubeIO/rubix-os/schema/schema"

type LocationProperties struct {
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
	Address     schema.Address     `json:"address"`
	City        schema.City        `json:"city"`
	State       schema.State       `json:"state"`
	Zip         schema.Zip         `json:"zip"`
	Country     schema.Country     `json:"country"`
	Lat         schema.Lat         `json:"lat"`
	Lon         schema.Lon         `json:"lon"`
	Timezone    schema.Timezone    `json:"timezone"`
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

type LocationGroupHostName struct {
	LocationName string `json:"location_name"`
	GroupName    string `json:"group_name"`
	HostName     string `json:"host_name"`
}
