package model

import "gorm.io/datatypes"

//WriterClone list of all the consumers
// a consumer
type WriterClone struct { //TODO the WriterClone needs to publish a COV event as for example we have 2x mqtt broker then the cov for a point maybe different when not going over the internet
	CommonUUID
	CommonThingClass
	CommonThingType
	ProducerUUID   string         `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"` // is the producer UUID
	WriterUUID     string         `json:"writer_uuid"`                                                                 // is the remote consumer UUID, ie: whatever is subscribing to this producer
	DataStore      datatypes.JSON `json:"data_store"`
	WriterSettings datatypes.JSON `json:"producer_settings"` //like cov for a point or whatever is needed  #TODO this is why it needs settings
	CommonCreated
}
