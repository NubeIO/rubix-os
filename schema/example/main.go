package main

import (
	"encoding/json"
	"github.com/NubeIO/rubix-os/schema/schema"
	"os"
)

func GetDeviceSchema() *Device {
	device := &Device{}
	schema.Set(device)
	return device
}

type Help string

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"50"`
	Default  string `json:"default" default:"lora"`
}

type Device struct {
	Name        NameStruct `json:"name"`
	AddressUUID struct {
		Title       string `json:"title" default:"name"`
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Min         int    `json:"minLength" default:"8"`
		Max         int    `json:"maxLength" default:"8"`
		DisplayName string `json:"display_name" default:"Address UUID"`
	} `json:"address_uuid"`
	SerialBaudRate struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"true"`
		Options  []int  `json:"options" default:"[38400, 999, 1234, 888]"`
		Default  int    `json:"default" default:"38400"`
		Help     Help   `json:"help" default:"this is help"`
	} `json:"serial_baud_rate"`
	Model struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"123\",\"33\",\"TH\",\"abc\",\"Test\"]"`
		Default  string   `json:"default" default:"THLM"`
	} `json:"model"`
}

func PrintJOSN(x interface{}) {
	ioWriter := os.Stdout
	w := json.NewEncoder(ioWriter)
	w.SetIndent("", "    ")
	w.Encode(x)
}

func main() {
	PrintJOSN(GetDeviceSchema())
}
