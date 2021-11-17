package serial_model

import (
	"github.com/NubeDev/flow-framework/plugin/defaults"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
)

// Message is a message wrapper with the channel, sender and recipient.
type Message struct {
	Msg         plugin.Message
	ChannelName string
	IsSend      bool
}

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"3"`
	Max      int    `json:"max" default:"20"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"80"`
}

type Network struct {
	Name        NameStruct        `json:"name"`
	Description DescriptionStruct `json:"description"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"lora"`
	} `json:"plugin_name"`
	TransportType struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"serial"`
	} `json:"transport_type"`
	SerialPort struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Min      int    `json:"min" default:"3"`
		Max      int    `json:"max" default:"20"`
	} `json:"serial_port"`
	BaudRate struct {
		Type     string `json:"type" default:"int"`
		Required bool   `json:"required" default:"true"`
		Default  int    `json:"default" default:"9600"`
	} `json:"baud_rate"`
}

func GetNetworkSchema() *Network {
	network := &Network{}
	defaults.Set(network)
	return network
}
