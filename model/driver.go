package model


/*

protocol
Communication protocol used for devices on a network

bacnet
ASHRAE building automation and control protocol

bluetooth
Short range wireless communication protocol

coap
Constrained Application Protocol

dali
Digital Addressable Lighting Interface protocol for lighting

ftp
File Transfer Protocol

haystack
Haystack HTTP protocol for exchanging tagged data

http
Hypertext Transfer Protocol which is foundation of the web

imap
Internet Message Access Protocol for retrieving email

knx
KNX protocol commonly used for lighting systems

LoRa
LoRa is a wireless network

LoRaWAN
LoRaWAN is a wireless network same as lora but for wide area

modbus
Register based communication protocol used with industrial devices

mqtt
Message Queuing Telemetry Transport publish/subscribe protocol

obix
XML based Open Building Information eXchange protocol

pop3
Post Office Protocol version 3 for retrieving email

smtp
Simple Mail Transfer Protocol for sending email

snmp
Simple Network Management Protocol for managing IP devices

sox
Sedona Framework UDP based communication protocol

zigbee
Low power wireless communication protocol for home automation

zwave
Low power wireless communication protocol for home automation

*/


//Driver the plugin when added will populate this and then based of its requirements will be able to use the network table or not
type Driver struct {
	PluginUUID 		string `json:"plugin_uuid"`
	TransportType 	string `json:"transport_type"` //serial, TCP, UDP, IP, MQTT
	Protocol 		string `json:"protocol"` //lora, bacnet, rest
	WriteableNetwork WriteableNetwork `json:"writeable_network"` //is this a network that supports write or its read only like lora
}

type DriverType string

var DriverTypeEnum = struct {
	Serial   	string
	IP   		string

}{
	Serial:   	"serial",
	IP:   		"IP",
}


type ProtocolType string

var ProtocolTypeEnum = struct {
	Client   				string
	Server   				string
	WirelessGateway   		string
	Stream   				string


}{
	Client:   "client",
	Server:   "server",
	WirelessGateway: "wireless",
	Stream:  "gateway",
}


type Protocol string

var ProtocolEnum = struct {
	REST   				string
	BACnet   			string
	Modbus   			string
	ModbusMaster   		string
	MQTT   				string
	LoRa   				string
	LoRaWAN   			string

}{
	REST:   			"rest",
	BACnet:   			"BACnet",
	Modbus:   			"Modbus",
	ModbusMaster:   	"ModbusMaster",
	MQTT:   			"MQTT",
	LoRa:   			"LoRa",
	LoRaWAN:   			"LoRaWAN",

}

type WriteableNetwork bool



