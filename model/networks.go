package model

import "time"

type IPType struct {
	REST string `json:"rest"`
	UDP  string `json:"udp"`
	MQTT string `json:"mqttClient"`
}

//IPNetwork type ip based network
type IPNetwork struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Token    string `json:"token"`
	IPType
}

type Network struct {
	CommonUUID
	CommonNameUnique
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	Manufacture      string `json:"manufacture,omitempty"`
	Model            string `json:"model,omitempty"`
	WriteableNetwork bool   `json:"writeable_network,omitempty"` //is this a network that supports write or its read only like lora
	CommonThingClass
	CommonThingRef
	CommonThingType
	TransportType    string        `json:"transport_type,omitempty"  gorm:"type:varchar(255);not null"` //serial
	PluginConfId     string        `json:"plugin_conf_id,omitempty" gorm:"TYPE:varchar(255) REFERENCES plugin_confs;not null;default:null"`
	PluginPath       string        `json:"plugin_name,omitempty"`
	NetworkInterface string        `json:"network_interface"`
	IP               string        `json:"ip"`
	Port             *int          `json:"port"`
	NetworkMask      *int          `json:"network_mask"`
	AddressID        string        `json:"address_id"`
	AddressUUID      string        `json:"address_uuid"`
	SerialPort       *string       `json:"serial_port,omitempty" gorm:"type:varchar(255);unique"`
	SerialBaudRate   *uint         `json:"serial_baud_rate,omitempty"` //9600
	SerialStopBits   *uint         `json:"serial_stop_bits,omitempty"` //1 or 2
	SerialParity     *string       `json:"serial_parity,omitempty"`    //odd, even, none
	SerialDataBits   *uint         `json:"serial_data_bits,omitempty"` //7 or 8
	SerialTimeout    *int          `json:"serial_timeout,omitempty"`
	SerialConnected  *bool         `json:"serial_connected,omitempty"`
	Host             *string       `json:"host,omitempty"`
	Devices          []*Device     `json:"devices,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Tags             []*Tag        `json:"tags,omitempty" gorm:"many2many:networks_tags;constraint:OnDelete:CASCADE"`
	MaxPollRate      time.Duration `json:"max_poll_rate,omitempty"` //Time is millis
}
