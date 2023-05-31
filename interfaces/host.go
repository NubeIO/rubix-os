package interfaces

import "github.com/NubeIO/lib-schema/schema"

type HostProperties struct {
	Name          schema.Name        `json:"name"`
	Enable        schema.Enable      `json:"enable"`
	Description   schema.Description `json:"description"`
	IP            schema.Host        `json:"ip"`
	BiosPort      schema.Port        `json:"bios_port"`
	Port          schema.Port        `json:"port"`
	HTTPS         schema.HTTPS       `json:"https"`
	ExternalToken schema.Token       `json:"external_token"`
}

func GetHostProperties() *HostProperties {
	m := &HostProperties{}
	m.Name.Min = 0
	m.IP.Default = "0.0.0.0"
	m.BiosPort.Title = "bios port"
	m.BiosPort.Default = 1659
	m.BiosPort.ReadOnly = false
	m.Port.Default = 1660
	m.Port.ReadOnly = false
	m.ExternalToken.Title = "external token"
	schema.Set(m)
	return m
}

type HostSchema struct {
	Required   []string        `json:"required"`
	Properties *HostProperties `json:"properties"`
}

func GetHostSchema() *HostSchema {
	m := &HostSchema{
		Required:   []string{"ip", "bios_port", "port"},
		Properties: GetHostProperties(),
	}
	return m
}
