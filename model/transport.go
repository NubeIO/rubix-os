package model

type SerialConnection struct {
	CommonUUID
	CommonEnable
	SerialPort  string `json:"serial_port"`
	BaudRate    uint   `json:"baud_rate"`
	StopBits    uint   `json:"stop_bits"`
	Parity      uint   `json:"parity"`
	DataBits    uint   `json:"data_bits"`
	Timeout     int    `json:"timeout"`
	Connected   bool   `json:"connected"`
	Error       bool   `json:"error"`
	NetworkUUID string `json:"network_uuid" gorm:"TYPE:varchar(255) REFERENCES networks"`
}

type IpConnection struct {
	CommonUUID
	Host        string `json:"host"`
	Port        int    `json:"port"`
	NetworkUUID string `json:"network_uuid" gorm:"TYPE:varchar(255) REFERENCES networks"`
}

type TransportBody struct {
	NetworkType      string           `json:"network_type"`
	TransportType    string           `json:"transport_type"`
	IpConnection     IpConnection     `json:"ip_connection"`
	SerialConnection SerialConnection `json:"serial_connection"`
}
