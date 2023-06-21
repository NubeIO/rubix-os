package schema

type DataType struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Data Type"`
	Options  []string `json:"enum" default:"[\"digital\",\"uint16\",\"int16\",\"uint32\",\"int32\",\"uint64\",\"int64\",\"float32\",\"float64\"]"`
	EnumName []string `json:"enumNames" default:"[\"digital\",\"uint16\",\"int16\",\"uint32\",\"int32\",\"uint64\",\"int64\",\"float32\",\"float64\"]"`
	Default  string   `json:"default" default:"uint16"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type ObjectEncoding struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Object Encoding (Endianness)"`
	Options  []string `json:"enum" default:"[\"beb_lew\",\"beb_bew\",\"leb_lew\",\"leb_bew\"]"`
	EnumName []string `json:"enumNames" default:"[\"beb_lew\",\"beb_bew\",\"leb_lew\",\"leb_bew\"]"`
	Default  string   `json:"default" default:"beb_lew"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type ObjectTypeModbus struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Object Type"`
	Options  []string `json:"enum" default:"[\"coil\",\"discrete_input\",\"input_register\",\"holding_register\",\"read_coil\",\"write_coil\",\"read_discrete_input\",\"read_register\",\"read_holding\",\"write_holding\"]"`
	EnumName []string `json:"enumNames" default:"[\"Coil\",\"Discrete Input\",\"Input Register\",\"Holding Register\",\"Read Coil\",\"Write Coil\",\"Read Discrete Input\",\"Read Input Register\",\"Read Holding Register\",\"Write Holding Register\"]"`
	Default  string   `json:"default" default:"coil"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type SerialPortModbus struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Serial Port"`
	Options  []string `json:"enum" default:"[\"/dev/ttyAMA0\",\"/dev/ttyRS485-1\",\"/dev/ttyRS485-2\",\"/data/socat/loRa1\",\"/dev/ttyUSB0\",\"/dev/ttyUSB1\",\"/dev/ttyUSB2\",\"/dev/ttyUSB3\",\"/dev/ttyUSB4\",\"/data/socat/serialBridge1\",\"/dev/ttyACM0\"]"`
	EnumName []string `json:"enumNames" default:"[\"/dev/ttyAMA0\",\"/dev/ttyRS485-1\",\"/dev/ttyRS485-2\",\"/data/socat/loRa1\",\"/dev/ttyUSB0\",\"/dev/ttyUSB1\",\"/dev/ttyUSB2\",\"/dev/ttyUSB3\",\"/dev/ttyUSB4\",\"/data/socat/serialBridge1\",\"/dev/ttyACM0\"]"`
	Default  string   `json:"default" default:"/dev/ttyRS485-2"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}
