package model

import (
	"github.com/NubeIO/flow-framework/plugin/defaults"
	"github.com/NubeIO/null"
)

type Priority struct {
	P1  null.Float `json:"_1,omitempty"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	P2  null.Float `json:"_2,omitempty"`
	P3  null.Float `json:"_3,omitempty"`
	P4  null.Float `json:"_4,omitempty"`
	P5  null.Float `json:"_5,omitempty"`
	P6  null.Float `json:"_6,omitempty"`
	P7  null.Float `json:"_7,omitempty"`
	P8  null.Float `json:"_8,omitempty"`
	P9  null.Float `json:"_9,omitempty"`
	P10 null.Float `json:"_10,omitempty"`
	P11 null.Float `json:"_11,omitempty"`
	P12 null.Float `json:"_12,omitempty"`
	P13 null.Float `json:"_13,omitempty"`
	P14 null.Float `json:"_14,omitempty"`
	P15 null.Float `json:"_15,omitempty"`
	P16 null.Float `json:"_16"` //removed and added to the point to save one DB write
}

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"3"`
	Max      int    `json:"max" default:"20"`
	Default  string `json:"default" default:"mb"`
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
	Options  []string `json:"options" default:"[\"serial\",\"ip\"]"`
	Default  string   `json:"default" default:"serial"`
}

type Network struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"modbus"`
	} `json:"plugin_name"`
	TransportType TransportTypeStruct `json:"transport_type"`
	SerialPort    struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Min      int    `json:"min" default:"3"`
		Max      int    `json:"max" default:"20"`
		Default  string `json:"default" default:"/dev/ttyAMA0"`
	} `json:"serial_port"`
	SerialBaudRate struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[9600, 38400]"`
		Default  int    `json:"default" default:"38400"`
	} `json:"serial_baud_rate"`
	SerialParity struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"odd\",\"even\",\"none\"]"`
		Default  string   `json:"default" default:"none"`
	} `json:"serial_parity"`
	SerialStopBits struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[1, 2]"`
		Default  int    `json:"default" default:"1"`
	} `json:"serial_stop_bits"`
	SerialDataBits struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[7, 8]"`
		Default  int    `json:"default" default:"8"`
	} `json:"serial_data_bits"`
}

type Device struct {
	Name          NameStruct          `json:"name"`
	Description   DescriptionStruct   `json:"description"`
	TransportType TransportTypeStruct `json:"transport_type"`
	AddressId     struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"1"`
	} `json:"address_id"`
	Host struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"false"`
		Options  string `json:"options" default:"192.168.15.10"`
		Default  string `json:"default" default:"192.168.15.10"`
	} `json:"host"`
	Port struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Options  int    `json:"options" default:"502"`
		Default  int    `json:"default" default:"502"`
	} `json:"port"`
	ZeroMode struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
		Default  bool   `json:"default" default:"false"`
	} `json:"zero_mode"`
}

type Point struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	AddressId   struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"1"`
	} `json:"address_id"`
	AddressLength struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Options  int    `json:"options" default:"1"`
		Default  int    `json:"default" default:"1"`
	} `json:"address_length"`
	ObjectType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"read_coil\",\"read_coils\",\"read_discrete_input\",\"read_discrete_inputs\",\"write_coil\",\"read_registers\",\"read_holding\",\"read_holdings\",\"write_holding\",\"write_holdings\",\"read_int_16\",\"write_int_16\",\"read_uint_16\",\"write_uint_16\",\"read_int_32\",\"write_int_32\",\"read_uint_32\",\"write_uint_32\",\"read_float_32\",\"write_float_32\",\"read_float_64\",\"write_float_64\"]"`
		Default  string   `json:"default" default:"read_coil"`
	} `json:"object_type"`
	ObjectEncoding struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"leb_lew\",\"leb_lew\",\"beb_lew\",\"beb_bew\"]"`
		Default  string   `json:"default" default:"beb_lew"`
	} `json:"object_encoding"`
	IsOutput struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
		Default  bool   `json:"default" default:"false"`
	} `json:"is_output"`
}

func GetNetworkSchema() *Network {
	network := &Network{}
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
