package linkmodel

type EnableStruct struct {
	Type     string `json:"type" default:"bool"`
	Required bool   `json:"required" default:"true"`
	Options  bool   `json:"options" default:"false"`
	Default  bool   `json:"default" default:"true"`
}

type NameStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"true"`
	Min      int    `json:"min" default:"2"`
	Max      int    `json:"max" default:"100"`
	Default  string `json:"default" default:"name_tmp"`
}

type DescriptionStruct struct {
	Type     string `json:"type" default:"string"`
	Required bool   `json:"required" default:"false"`
	Min      int    `json:"min" default:"0"`
	Max      int    `json:"max" default:"100"`
	Default  string `json:"default" default:"desc_tmp"`
}

type Dropdown struct {
	Type     string   `json:"type" default:"array"`
	Required bool     `json:"required" default:"true"`
	Options  []string `json:"options" default:"[]"`
	Default  string   `json:"default" default:""`
}

type SchemaNetwork struct {
	Enable      EnableStruct `json:"enable"`
	AddressUUID Dropdown     `json:"address_uuid"`
	PluginName  struct {
		Type     string `json:"type" default:"string"`
		Required bool   `json:"required" default:"true"`
		Default  string `json:"default" default:"networklinker"`
	} `json:"plugin_name"`
}

type SchemaDevice struct {
	Enable      EnableStruct `json:"enable"`
	Name        NameStruct   `json:"name"`
	AddressUUID Dropdown     `json:"address_uuid"`
}

type SchemaPoint struct {
	Enable EnableStruct `json:"enable"`
	Name   NameStruct   `json:"name"`
}
