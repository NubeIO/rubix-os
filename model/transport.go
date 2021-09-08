package model

type SerialConnection struct {
	SerialPort  string `json:"serial_port"`
	Enable      bool   `json:"enable"`
	Port        string `json:"port"`
	BaudRate    int    `json:"baud_rate"`
	StopBits    int    `json:"stop_bits"`
	Parity      int    `json:"parity"`
	DataBits    int    `json:"data_bits"`
	Timeout     int    `json:"timeout"`
	Connected   bool   `json:"connected"`
	Error       bool   `json:"error"`
	NetworkUUID string `json:"network_uuid" gorm:"TYPE:varchar(255) REFERENCES networks;null;default:null"`
}

type IpConnection struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	NetworkUUID string `json:"network_uuid" gorm:"TYPE:varchar(255) REFERENCES networks;null;default:null"`
}

type TransportBody struct {
	NetworkType      string           `json:"network_type"`
	TransportType    string           `json:"transport_type"`
	IpConnection     IpConnection     `json:"ip_connection"`
	SerialConnection SerialConnection `json:"serial_connection"`
}
