package model

//type StreamList struct { //TODO add is in so multiple flow networks can tap into an existing stream
//	StreamUUID      string `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;null;default:null"`
//	FlowNetworkUUID string `json:"flow_network_uuid" gorm:"TYPE:string REFERENCES flow_network;null;default:null"`
//}

type Stream struct {
	CommonUUID
	CommonName
	CommonDescription
	FlowNetworks []*FlowNetwork `gorm:"many2many:streams_flow_networks;"`
	IsConsumer   bool           `json:"is_consumer"`
	CommonEnable
	Producer     []Producer     `json:"producers" gorm:"constraint:OnDelete:CASCADE;"`
	Consumer     []Consumer     `json:"consumer" gorm:"constraint:OnDelete:CASCADE;"`
	CommandGroup []CommandGroup `json:"command_group" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}
