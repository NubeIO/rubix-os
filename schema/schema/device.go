package schema

type Host struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Host IP Address"`
	Default  string `json:"default" default:"0.0.0.0"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Ip struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"Host IP Address"`
	Default  string `json:"default" default:"0.0.0.0"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type DeviceObjectId struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Object ID"`
	Default  int    `json:"default" default:"2508"`
	Min      int    `json:"minLength" default:"0"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type AddressId struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Address ID"`
	Default  int    `json:"default" default:"1"`
	Min      int    `json:"minLength" default:"0"`
	Max      int    `json:"maxLength" default:"255"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type AddressLength struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Address Length"`
	Default  int    `json:"default" default:"1"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type NetworkNumber struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Network Number"`
	Default  int    `json:"default" default:"0"`
	Min      int    `json:"minLength" default:"0"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type DeviceMac struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:"Device MS/TP MAC Address"`
	Default  int    `json:"default" default:"0"`
	Min      int    `json:"minLength" default:"0"`
	Max      int    `json:"maxLength" default:"255"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}
