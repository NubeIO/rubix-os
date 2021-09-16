package model

type Stream struct {
	CommonUUID
	CommonName
	CommonDescription
	IsConsumer bool `json:"is_consumer"`
	CommonEnable
	FlowNetworks  []*FlowNetwork  `json:"flow_networks" gorm:"many2many:flow_networks_streams;constraint:OnDelete:CASCADE"`
	Producers     []*Producer     `json:"producers" gorm:"constraint:OnDelete:CASCADE;"`
	Consumers     []*Consumer     `json:"consumers" gorm:"constraint:OnDelete:CASCADE;"`
	CommandGroups []*CommandGroup `json:"command_groups" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}
