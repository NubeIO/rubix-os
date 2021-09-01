package model

import (
	"gorm.io/datatypes"
	"time"
)


type ProducerType struct {
	Network   		string `json:"network"`
	Job   			string `json:"job"`
	Point   		string `json:"point"`
	Alarm   		string `json:"alarm"`
}

type ProducerUse struct {
	Local   		 string `json:"local"`
	Remote   		string `json:"remote"`
	Plugin   		string `json:"plugin"`

}

//WriterClone list of all the consumers
// a consumer
type WriterClone struct { //TODO the WriterClone needs to publish a COV event as for example we have 2x mqtt broker then the cov for a point maybe different when not going over the internet
	CommonUUID
	WriterType 			string  `json:"writer_type"` //point, schedule, job, network
	ProducerUUID 		string  `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"` // is the producer UUID
	WriterUUID 			string 	`json:"writer_uuid"`  // is the remote consumer UUID, ie: whatever is subscribing to this producer
	DataStore 			datatypes.JSON  `json:"data_store"`
	WriterSettings 		datatypes.JSON  `json:"producer_settings"` //like cov for a point or whatever is needed  #TODO this is why it needs settings
	CommonCreated
}


//Producer a producer is a placeholder to register an object to enable consumers to
// A producer for example is a point, Something that makes data, and the subscriber would have a consumer to it, Like grafana reading and writing to it from edge to cloud or wires over rest(peer to peer)
type Producer struct {
	CommonProducer
	CurrentWriterCloneUUID  string  `json:"current_writer_clone_uuid"`
	ProducerType 			string  `json:"producer_type"` //point, schedule, job, network
	EnableHistory 			bool 	`json:"enable_history"`
	ProducerApplication 	string 	`json:"producer_application"`
	StreamUUID     			string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ProducerThingUUID 		string  `json:"producer_thing_uuid"` //a producer_uuid is the point uuid
	WriterClone				[]WriterClone `json:"writer_clones" gorm:"constraint:OnDelete:CASCADE;"`
	ProducerHistory			[]ProducerHistory `json:"producer_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

//ProducerHistory for storing the history
type ProducerHistory struct {
	CommonUUID
	ProducerUUID    		string  	`json:"producer_uuid" gorm:"TYPE:varchar(255) REFERENCES producers;not null;default:null"`
	CurrentWriterCloneUUID  string  	`json:"current_writer_clone_uuid"`
	DataStore 			datatypes.JSON  `json:"data_store"`
	Timestamp    		time.Time 		`json:"timestamp"`

}