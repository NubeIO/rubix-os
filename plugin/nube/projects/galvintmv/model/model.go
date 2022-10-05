package model

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
	Max      int    `json:"max" default:"100"`
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
	Options  []string `json:"options" default:"[\"serial\",\"ip\",\"LoRa\"]"`
	Default  string   `json:"default" default:"serial"`
}

type Network struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"mb_net"`
		DisplayName string `json:"display_name" default:"Network Name"`
	} `json:"name"`
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
		Min      int    `json:"min" default:"0"`
		Max      int    `json:"max" default:"100"`
		Default  string `json:"default" default:"/dev/ttyRS485-1"`
	} `json:"serial_port"`
	SerialBaudRate struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[9600, 38400, 57600, 115200]"`
		Default  int    `json:"default" default:"38400"`
	} `json:"serial_baud_rate"`
	SerialParity struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"odd\",\"even\",\"none\"]"`
		Default  string   `json:"default" default:"none"`
	} `json:"serial_parity"`
	SerialDataBits struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[7, 8]"`
		Default  int    `json:"default" default:"8"`
	} `json:"serial_data_bits"`
	SerialStopBits struct {
		Type     string `json:"type" default:"array"`
		Required bool   `json:"required" default:"false"`
		Options  []int  `json:"options" default:"[1, 2]"`
		Default  int    `json:"default" default:"1"`
	} `json:"serial_stop_bits"`
	MaxPollRate struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Max Poll Rate (seconds)"`
		Default     *float64 `json:"default" default:"0.1"`
	} `json:"max_poll_rate"`
	SerialTimeout struct {
		Type        string   `json:"type" default:"int"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Polling Timeout"`
		Default     *float64 `json:"default" default:"2"`
	} `json:"serial_timeout"`
	AutoMappingNetworksSelection struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"disable\",\"self-mapping\",\"rubix-io-to-bacnetserver\"]"`
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
	AutoMappingEnableHistories struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
		Default  *bool  `json:"default" default:"false"`
	} `json:"auto_mapping_enable_histories"`
}

type Device struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"mb_dev"`
		DisplayName string `json:"display_name" default:"Device Name"`
	} `json:"name"`
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
		Default  string `json:"default" default:""`
	} `json:"host"`
	Port struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"false"`
		Options  int    `json:"options" default:"502"`
		Default  int    `json:"default" default:""`
	} `json:"port"`
	FastPollRate struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Fast Poll Rate (seconds)"`
		Default     *float64 `json:"default" default:"5"`
	} `json:"fast_poll_rate"`
	NormalPollRate struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Normal Poll Rate (seconds)"`
		Default     *float64 `json:"default" default:"30"`
	} `json:"normal_poll_rate"`
	SlowPollRate struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"true"`
		Options     int      `json:"options" default:"1"`
		DisplayName string   `json:"display_name" default:"Slow Poll Rate (seconds)"`
		Default     *float64 `json:"default" default:"120"`
	} `json:"slow_poll_rate"`
	ZeroMode struct {
		Type     string `json:"type" default:"bool"`
		Required bool   `json:"required" default:"true"`
		Options  bool   `json:"options" default:"false"`
		Default  *bool  `json:"default" default:"true"`
	} `json:"zero_mode"`
}

type Point struct {
	Enable EnableStruct `json:"enable"`
	Name   struct {
		Type        string `json:"type" default:"string"`
		Required    bool   `json:"required" default:"true"`
		Default     string `json:"default" default:"mb_pnt"`
		DisplayName string `json:"display_name" default:"Point Name"`
	} `json:"name"`
	Description DescriptionStruct `json:"description"`
	ObjectType  struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"read_coil\",\"write_coil\",\"read_discrete_input\",\"read_register\",\"read_holding\",\"write_holding\"]"`
		Default  string   `json:"default" default:"read_holding"`
	} `json:"object_type"`
	AddressId struct {
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
	DataType struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"digital\",\"uint16\",\"int16\",\"uint32\",\"int32\",\"uint64\",\"int64\",\"float32\",\"float64\"]"`
		Default  string   `json:"default" default:"uint16"`
	} `json:"data_type"`
	WriteMode struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"read_once\",\"read_only\",\"write_once\",\"write_once_read_once\",\"write_always\",\"write_once_then_read\",\"write_and_maintain\"]"`
		Default  string   `json:"default" default:"read_only"`
	} `json:"write_mode"`
	PollPriority struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"high\",\"normal\",\"low\"]"`
		Default  string   `json:"default" default:"normal"`
	} `json:"poll_priority"`
	PollRate struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
		Options  []string `json:"options" default:"[\"fast\",\"normal\",\"slow\"]"`
		Default  string   `json:"default" default:"normal"`
	} `json:"poll_rate"`
	ObjectEncoding struct {
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"false"`
		Options  []string `json:"options" default:"[\"leb_lew\",\"leb_bew\",\"beb_lew\",\"beb_bew\"]"`
		Default  string   `json:"default" default:"beb_lew"`
	} `json:"object_encoding"`
	ScaleEnable struct {
		Type        string `json:"type" default:"bool"`
		Required    bool   `json:"required" default:"true"`
		Default     *bool  `json:"default" default:"false"`
		DisplayName string `json:"display_name" default:"Scale/Limit Enable"`
	} `json:"scale_enable"`
	ScaleInMin struct {
		Type        string `json:"type" default:"float"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"0"`
		DisplayName string `json:"display_name" default:"Scale: Input Min"`
	} `json:"scale_in_min"`
	ScaleInMax struct {
		Type        string `json:"type" default:"float"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"0"`
		DisplayName string `json:"display_name" default:"Scale: Input Max"`
	} `json:"scale_in_max"`
	ScaleOutMin struct {
		Type        string `json:"type" default:"float"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"0"`
		DisplayName string `json:"display_name" default:"Scale/Limit: Output Min"`
	} `json:"scale_out_min"`
	ScaleOutMax struct {
		Type        string `json:"type" default:"float"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"0"`
		DisplayName string `json:"display_name" default:"Scale/Limit: Output Max"`
	} `json:"scale_out_max"`
	MultiplicationFactor struct {
		Type        string `json:"type" default:"float"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"1"`
		DisplayName string `json:"display_name" default:"Multiplication Factor"`
	} `json:"multiplication_factor"`
	Offset struct {
		Type        string `json:"type" default:"float"`
		Required    bool   `json:"required" default:"false"`
		Default     string `json:"default" default:"0"`
		DisplayName string `json:"display_name" default:"Offset"`
	} `json:"offset"`
	Fallback struct {
		Type        string   `json:"type" default:"float"`
		Required    bool     `json:"required" default:"false"`
		Default     *float64 `json:"default" default:""`
		DisplayName string   `json:"display_name" default:"Fallback Value"`
		Nullable    bool     `json:"nullable" default:"true"`
	} `json:"fallback"`
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
