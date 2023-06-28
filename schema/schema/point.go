package schema

type ObjectId struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Object ID"`
	Default  int    `json:"default" default:"1"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type ObjectType struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Object Type"`
	Options  []string `json:"enum" default:"[\"analog_input\",\"analog_value\",\"analog_output\",\"binary_input\",\"binary_value\",\"binary_output\",\"multi_state_input\",\"multi_state_value\",\"multi_state_output\"]"`
	EnumName []string `json:"enumNames" default:"[\"Analog Input (AI)\",\"Analog Value (AV)\",\"Analog Output (AO)\",\"Binary Input (BI)\",\"Binary value (BV)\",\"Binary Output (BO)\",\"Multi State Input (MSI)\",\"Multi State Value (MSV)\",\"Multi State Output (MSO)\"]"`
	Default  string   `json:"default" default:"analog_value"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type WriteMode struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Write Mode"`
	Options  []string `json:"enum" default:"[\"read_only\",\"write_once\",\"write_once_read_once\",\"write_always\",\"write_once_then_read\",\"write_and_maintain\"]"`
	EnumName []string `json:"enumNames" default:"[\"read only\",\"write once\",\"write once read once\",\"write always\",\"write once then read\",\"write and maintain\"]"`
	Default  string   `json:"default" default:"read_only"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type WritePriority struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Write Priority"`
	Min      int    `json:"minLength" default:"1"`
	Max      int    `json:"maxLength" default:"16"`
	Default  int    `json:"default" default:"16"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type IoNumber struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"IO Number"`
	Options  []string `json:"enum" default:"[\"UI1\",\"UI2\",\"UI3\",\"UI4\",\"UI5\",\"UI6\",\"UI7\",\"UI8\",\"UO1\",\"UO2\",\"UO3\",\"UO4\",\"UO5\",\"UO6\",\"DO1\",\"DO2\"]"`
	EnumName []string `json:"enumNames" default:"[\"UI1\",\"UI2\",\"UI3\",\"UI4\",\"UI5\",\"UI6\",\"UI7\",\"UI8\",\"UO1\",\"UO2\",\"UO3\",\"UO4\",\"UO5\",\"UO6\",\"DO1\",\"DO2\"]"`
	Default  string   `json:"default" default:"UI1"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

const (
	IOTypeRaw        = "raw"
	IOTypeTherm10kT2 = "thermistor_10k_type_2"
	IOTypeDigital    = "digital"
	IOTypeVDC        = "voltage_dc"
	IOTypeCurrent    = "current"
)

type IoType struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"IO Type"`
	Options  []string `json:"enum" default:"[\"digital\",\"voltage_dc\",\"thermistor_10k_type_2\",\"current\",\"raw\"]"`
	EnumName []string `json:"enumNames" default:"[\"digital\",\"voltage_dc\",\"thermistor_10k_type_2\",\"current\",\"raw\"]"`
	Default  string   `json:"default" default:"raw"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type ScaleEnable struct {
	Type     string `json:"type" default:"boolean"`
	Title    string `json:"title" default:"Scale Enable"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type ScaleInMin struct {
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:"Scale: Input/Device Min"`
	Default  float64 `json:"default" default:"0"`
	ReadOnly bool    `json:"readOnly" default:"false"`
}

type ScaleInMax struct {
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:"Scale: Input/Device Max"`
	Default  float64 `json:"default" default:"0"`
	ReadOnly bool    `json:"readOnly" default:"false"`
}

type ScaleOutMin struct {
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:"Scale: Output/Point Min"`
	Default  float64 `json:"default" default:"0"`
	ReadOnly bool    `json:"readOnly" default:"false"`
}

type ScaleOutMax struct {
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:"Scale: Output/Point Max"`
	Default  float64 `json:"default" default:"0"`
	ReadOnly bool    `json:"readOnly" default:"false"`
}

type MultiplicationFactor struct {
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:"Multiplication Factor"`
	Default  float64 `json:"default" default:"1"`
	ReadOnly bool    `json:"readOnly" default:"false"`
}

type Offset struct {
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:"Offset"`
	Default  float64 `json:"default" default:"0"`
	ReadOnly bool    `json:"readOnly" default:"false"`
}

type Fallback struct {
	Type     string   `json:"type" default:"number"`
	Title    string   `json:"title" default:"Fallback"`
	Default  *float64 `json:"default" default:""`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type Decimal struct {
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:"Round To Decimals"`
	Default  float64 `json:"default" default:"2"`
	ReadOnly bool    `json:"readOnly" default:"false"`
}

type HistoryType struct {
	Type    string   `json:"type" default:"string"`
	Title   string   `json:"title" default:"History Type"`
	Options []string `json:"enum" default:"[\"COV\",\"INTERVAL\",\"COV_AND_INTERVAL\"]"`
	Default string   `json:"default" default:"INTERVAL"`
}

type HistoryInterval struct {
	Type    string `json:"type" default:"number"`
	Title   string `json:"title" default:"History Interval"`
	Default *int   `json:"default" default:"15"`
}

type HistoryCOVThreshold struct {
	Type    string   `json:"type" default:"number"`
	Title   string   `json:"title" default:"History COV Threshold"`
	Default *float64 `json:"default" default:"0.01"`
}
