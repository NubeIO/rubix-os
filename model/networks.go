package model

type IPType struct {
	REST string `json:"rest"`
	UDP  string `json:"udp"`
	MQTT string `json:"mqttClient"`
}

//IPNetwork type ip based network
type IPNetwork struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Token    string `json:"token"`
	IPType
}

type Network struct {
	CommonUUID
	CommonNameUnique
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	Manufacture      string `json:"manufacture,omitempty"`
	Model            string `json:"model,omitempty"`
	WriteableNetwork bool   `json:"writeable_network,omitempty"` //is this a network that supports write or its read only like lora
	CommonThingClass
	CommonThingRef
	CommonThingType
	TransportType    string            `json:"transport_type,omitempty"  gorm:"type:varchar(255);not null"` //serial
	PluginConfId     string            `json:"plugin_conf_id,omitempty" gorm:"TYPE:varchar(255) REFERENCES plugin_confs;not null;default:null"`
	PluginPath       string            `json:"plugin_name,omitempty"`
	Devices          []*Device         `json:"devices,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	SerialConnection *SerialConnection `json:"serial_connection,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IpConnection     *IpConnection     `json:"ip_connection,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tags             []*Tag            `json:"tags,omitempty" gorm:"many2many:networks_tags;constraint:OnDelete:CASCADE"`
}
