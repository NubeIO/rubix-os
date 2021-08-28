package model

/*

Rubix V2

Network
- Any new network is called a ProducerNetwork for example modbus or lora
- Any ProducerNetwork doesn't need to use point-mapping to get the data into the cloud
- Any ProducerNetwork can be added to a one or more ProducerNetworks, The ProducerNetwork will keep a leger or which Producer are reading/writing to its points

The ProducerNetwork is the remote rubix device.
So when the ProducerNetwork (the producer network producers data ie: lora) network has a connection with the ProducerNetwork the ProducerNetwork keeps a ledger of the ProducerPoints

ProducerNetwork and ProducerNetwork Jobs
- publish any CRUD updates to all producers (ie when a point is deleted or the name is updated)
- publish any COV events

ProducerNetwork settings
- COV will be set in the producer

ProducerNetwork settings (these settings are not like 2-way meaning that in the ProducerNetwork if the COV is updated it will not affect the ProducerNetwork setting)
- as this would be considered a normal point in the ProducerNetwork this point will have all the same settings ie: history, cov and so on

CommandGroup
- is for issuing global schedule writes or global point writes (as in send a value to any point associated with this group)

TimeOverride
- where a point value can be overridden for a duration of time



REST calls
ProducerNetwork
- can call all attributes

*/




type Stream struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonCreated
	IsSubscription  	bool   `json:"is_subscription"`
	FlowNetworkUUID 	string `json:"flow_network_uuid" gorm:"TYPE:varchar(255) REFERENCES flow_networks;not null;default:null"`
	Producer			[]Producer `json:"producers" gorm:"constraint:OnDelete:CASCADE;"`
	Subscription		[]Subscription `json:"subscription" gorm:"constraint:OnDelete:CASCADE;"`
	CommandGroup		[]CommandGroup `json:"command_group" gorm:"constraint:OnDelete:CASCADE;"`
}