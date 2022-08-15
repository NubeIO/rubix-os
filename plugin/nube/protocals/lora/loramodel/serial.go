package loramodel

import (
	"github.com/NubeIO/flow-framework/plugin/defaults"
)

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"true"`
	Options  bool   `json:"options" default:"false"`
	Default  *bool  `json:"default" default:"true"`
}

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"50"`
	Default  string `json:"default" default:"lora"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type TransportTypeStruct struct {
	Type     string   `json:"type" default:"array"`
	Required bool     `json:"required" default:"true"`
	Options  []string `json:"options" default:"[\"serial\"]"`
	Default  string   `json:"default" default:"serial"`
}

type Network struct {
	Enable      EnableStruct      `json:"enable"`
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"lora"`
	} `json:"plugin_name"`
	TransportType TransportTypeStruct `json:"transport_type"`
	SerialPort    struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"/dev/ttyAMA0\",\"/data/socat/loRa1\",\"/dev/ttyUSB0\",\"/dev/ttyUSB1\",\"/dev/ttyUSB2\",\"/dev/ttyUSB3\",\"/dev/ttyUSB4\"]"`
		Default  string   `json:"default" default:"/dev/ttyAMA0"`
	} `json:"serial_port"`
	SerialBaudRate struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"true"`
		Options  []int  `json:"options" default:"[38400]"`
		Default  int    `json:"default" default:"38400"`
	} `json:"serial_baud_rate"`
	AutoMappingNetworksSelection struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"disable\",\"lora-to-bacnetserver\"]"`
		Default  string   `json:"default" default:""`
	} `json:"auto_mapping_networks_selection"`
	AutoMappingFlowNetworkName struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"200"`
		Default  string `json:"default" default:"local"`
	} `json:"auto_mapping_flow_network_name"`
	AutoMappingFlowNetworkUUID struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"200"`
		Default  string `json:"default" default:""`
	} `json:"auto_mapping_flow_network_uuid"`
}

type Device struct {
	Enable      EnableStruct      `json:"enable"`
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	AddressUUID struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Min         int    `json:"min" default:"8"`
		Max         int    `json:"max" default:"8"`
		DisplayName string `json:"display_name" default:"Address UUID"`
	} `json:"address_uuid"`
	Model struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"THLM\",\"THL\",\"TH\",\"MicroEdge\",\"ZipHydroTap\"]"`
		Default  string   `json:"default" default:"THLM"`
	} `json:"model"`
}

type Point struct {
	Enable      EnableStruct      `json:"enable"`
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	AddressUUID struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Min         int    `json:"min" default:"8"`
		Max         int    `json:"max" default:"8"`
		DisplayName string `json:"display_name" default:"Address UUID"`
	} `json:"address_uuid"`
	IoNumber struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"true"`
	} `json:"io_number"`
	IoType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"0-10dc\",\"0-40ma\",\"thermistor\"]"`
	} `json:"io_type"`
	Eval struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"(x + 0) + 0"`
		DisplayName string `json:"display_name" default:"math expression"`
	} `json:"eval_expression"`
	ThingClass struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"point\"]"`
	} `json:"thing_class"`
	ThingType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"point\"]"`
	} `json:"thing_type"`
}

func GetNetworkSchema() *Network {
	network := &Network{}
	//ports, err := serial.GetPortsList()
	//if err != nil {
	//	log.Errorf("lora: get serial ports for schema err:%s", err.Error())
	//} else {
	//	network.SerialPort.Options = ports
	//}
	defaults.Set(network)
	return network
}

func GetDeviceSchema() *Device {
	device := &Device{}
	defaults.Set(device)
	return device
}

func GetPointSchema() *Point {
	point := &Point{}
	defaults.Set(point)
	return point
}
