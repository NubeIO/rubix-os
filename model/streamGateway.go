package model



type StreamList struct { //TODO add is in so multiple flow networks can tap into an existing stream
	UUID			string 	`json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
}


type Stream struct {
	CommonUUID
	CommonName
	CommonDescription
	StreamListUUID 		string `json:"stream_list_uuid" gorm:"TYPE:varchar(255) REFERENCES stream_lists;not null;default:null"`
	IsConsumer  	bool   `json:"is_consumer"`
	CommonEnable
	Producer			[]Producer `json:"producers" gorm:"constraint:OnDelete:CASCADE;"`
	Consumer			[]Consumer `json:"consumer" gorm:"constraint:OnDelete:CASCADE;"`
	CommandGroup		[]CommandGroup `json:"command_group" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}