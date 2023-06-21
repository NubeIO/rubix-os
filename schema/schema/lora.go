package schema

type SerialPortLora struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Serial Port"`
	Options  []string `json:"enum" default:"[\"/dev/ttyAMA0\",\"/dev/ttyRS485-1\",\"/dev/ttyRS485-2\",\"/data/socat/loRa1\",\"/dev/ttyUSB0\",\"/dev/ttyUSB1\",\"/dev/ttyUSB2\",\"/dev/ttyUSB3\",\"/dev/ttyUSB4\",\"/data/socat/serialBridge1\",\"/dev/ttyACM0\"]"`
	EnumName []string `json:"enumNames" default:"[\"/dev/ttyAMA0\",\"/dev/ttyRS485-1\",\"/dev/ttyRS485-2\",\"/data/socat/loRa1\",\"/dev/ttyUSB0\",\"/dev/ttyUSB1\",\"/dev/ttyUSB2\",\"/dev/ttyUSB3\",\"/dev/ttyUSB4\",\"/data/socat/serialBridge1\",\"/dev/ttyACM0\"]"`
	Default  string   `json:"default" default:"/data/socat/LoRa1"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}
