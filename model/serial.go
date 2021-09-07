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
