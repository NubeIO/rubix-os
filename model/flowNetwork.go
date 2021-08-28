package model

//FlowNetwork flow network
type FlowNetwork struct {
	CommonUUID
	CommonName
	CommonDescription
	IsRemote 	bool `json:"is_remote"`
	RemoteUUID  string `json:"remote_uuid" gorm:"type:varchar(255);unique;not null"`
	FlowIP 		string `json:"flow_ip"`
	FlowPort 	string `json:"flow_port"`
	FlowHTTPS 	bool `json:"flow_https"`
	FlowUsername string `json:"flow_username"`
	FlowPassword string `json:"flow_password"`
	MqttIP 		string `json:"mqtt_ip"`
	MqttPort 	string `json:"mqtt_port"`
	MqttHTTPS 	bool `json:"mqtt_https"`
	MqttUsername string `json:"mqtt_username"`
	MqttPassword string `json:"mqtt_password"`
	CommonCreated
	Stream		[]Stream `json:"streams" gorm:"constraint:OnDelete:CASCADE;"`

}
