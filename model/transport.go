package model

var NetworkRTUSpeed = struct {
	_9600   int
	_19200  int
	_38400  int
	_57600  int
	_115200 int
	_128000 int
}{
	_9600:   9600,
	_19200:  19200,
	_38400:  38400,
	_57600:  57600,
	_115200: 115200,
	_128000: 128000,
}

var NetworkRTUPorts = struct {
	ttyRS4851 string
	ttyRS4852 string
	ttyXBEE1  string
	ttyXBEE2  string
	ttyAMA0   string
	ttyAMA1   string
	ttyACM0   string
	ttyACM1   string
	ttyUSB0   string
	ttyUSB1   string
	ttyUSB2   string
	ttyUSB3   string
	ttyUSB4   string
	ttyUSB5   string
	ttyUSB6   string
}{
	ttyRS4851: "/dev/ttyRS485-1",
	ttyRS4852: "/dev/ttyRS485-2",
	ttyXBEE1:  "/dev/ttyXBEE-1",
	ttyXBEE2:  "/dev/ttyXBEE-2",
	ttyAMA0:   "/dev/ttyAMA0",
	ttyAMA1:   "/dev/ttyAMA1",
	ttyACM0:   "/dev/ttyACM0",
	ttyACM1:   "/dev/ttyACM1",
	ttyUSB0:   "/dev/ttyUSB0",
	ttyUSB1:   "/dev/ttyUSB1",
	ttyUSB2:   "/dev/ttyUSB2",
	ttyUSB3:   "/dev/ttyUSB3",
	ttyUSB4:   "/dev/ttyUSB4",
	ttyUSB5:   "/dev/ttyUSB5",
	ttyUSB6:   "/dev/ttyUSB6",
}

var SerialParity = struct {
	None string `json:"none"`
	Odd  string `json:"odd"`
	Even string `json:"even"`
}{
	None: "none",
	Odd:  "odd",
	Even: "even",
}

type SerialConnection struct {
	CommonUUID
	CommonEnable
	SerialPort  string `json:"serial_port" gorm:"type:varchar(255);unique"`
	BaudRate    uint   `json:"baud_rate"` //9600
	StopBits    uint   `json:"stop_bits"`
	Parity      string `json:"parity"`
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
