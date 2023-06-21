package schema

type UUID struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"UUID"`
	ReadOnly bool   `json:"readOnly" default:"true"`
}

type AddressUUID struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Address UUID"`
	Min      int    `json:"minLength" default:"1"`
	Max      int    `json:"maxLength" default:"100"`
	Default  string `json:"default" default:""`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Name struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Name"`
	Min      int    `json:"minLength" default:"2"`
	Max      int    `json:"maxLength" default:"200"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Model struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Model"`
	Options  []string `json:"enum" default:"[]"`
	EnumName []string `json:"enumNames" default:"[]"`
	Default  string   `json:"default" default:""`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type Username struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Username"`
	Min      int    `json:"minLength" default:"2"`
	Max      int    `json:"maxLength" default:"50"`
	Default  string `json:"default" default:"admin"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Password struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Password"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Token struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Token"`
	Min      int    `json:"minLength" default:"0"`
	Max      int    `json:"maxLength" default:"200"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Description struct {
	Type  string `json:"type" default:"string"`
	Title string `json:"title" default:"Description"`
}

type Enable struct {
	Type     string `json:"type" default:"boolean"`
	Title    string `json:"title" default:"Enable"`
	Default  bool   `json:"default" default:"true"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type HistoryEnable struct {
	Type     string `json:"type" default:"boolean"`
	Title    string `json:"title" default:"History Enable"`
	Default  bool   `json:"default" default:"false"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type HistoryEnableDefaultTrue struct {
	Type     string `json:"type" default:"boolean"`
	Title    string `json:"title" default:"History Enable"`
	Default  bool   `json:"default" default:"true"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Product struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Product"`
	Options  []string `json:"enum" default:"[\"RubixCompute\",\"RubixCompute5\",\"RubixComputeIO\",\"Edge28\",\"Nuc\",\"Server\"]"`
	EnumName []string `json:"enumNames" default:"[\"RubixCompute\",\"RubixCompute5\",\"RubixComputeIO\",\"Edge28\",\"Nuc\",\"Server\"]"`
	Help     string   `json:"help" default:"a nube product type or a general linux server"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type Interface struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"Network Interface"`
	Options  []string `json:"enum" default:"[]"`
	Default  string   `json:"default" default:"eth0"`
	Help     string   `json:"help" default:"host network interface card, eg eth0"`
	ReadOnly bool     `json:"readOnly" default:"false"`
}

type Netmask struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Netmask"`
	Default  string `json:"default" default:"255.255.255.0"`
	Help     string `json:"help" default:"ip netmask address eg, 255.255.255.0"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type SubNetMask struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Subnet Mask"`
	Min      int    `json:"minLength" default:"8"`
	Max      int    `json:"maxLength" default:"30"`
	Default  int    `json:"default" default:"24"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Gateway struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Gateway"`
	Help     string `json:"help" default:"ip gateway address eg, 192.168.15.1"`
	ReadOnly bool   `json:"readOnly" default:"false"`
	Default  string `json:"default" default:"192.168.15.1"`
}

type HTTPS struct {
	Type     string `json:"type" default:"boolean"`
	Title    string `json:"title" default:"Enable HTTPS"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Port struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Port"`
	Min      int    `json:"minLength" default:"2"`
	Max      int    `json:"maxLength" default:"65535"`
	Default  int    `json:"default" default:"1660"`
	Help     string `json:"help" default:"ip port, eg port 1660 192.168.15.10:1660"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type PluginName struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Plugin"`
	ReadOnly bool   `json:"readOnly" default:"true"`
}

type OptionOneOf struct {
	Const string `json:"const"`
	Title string `json:"Title"`
}

type OptionOneOfInt struct {
	Const int    `json:"const"`
	Title string `json:"Title"`
}
