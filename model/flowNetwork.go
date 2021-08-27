package model

//FlowNetwork flow network
type FlowNetwork struct {
	CommonUUID
	CommonName
	CommonDescription
	IsRemote 	bool
	FlowIP 		string
	FlowPort 	string
	FlowHTTPS 	bool
	FlowUsername string
	FlowPassword string
	MqttIP 		string
	MqttPort 	string
	MqttHTTPS 	bool
	MqttUsername string
	MqttPassword string
	Stream		[]Stream `json:"streams" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}
