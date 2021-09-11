package model

import (
	"gorm.io/datatypes"
	"time"
)

type ProducerType struct {
	Network string `json:"network"`
	Job     string `json:"job"`
	Point   string `json:"point"`
	Alarm   string `json:"alarm"`
}

type ProducerUse struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
	Plugin string `json:"plugin"`
}

//WriterClone list of all the consumers
// a consumer
type WriterClone struct { //TODO the WriterClone needs to publish a COV event as for example we have 2x mqtt broker then the cov for a point maybe different when not going over the internet
	CommonUUID
	WriterType     string         `json:"writer_type"`                                                                 //point, schedule, job, network
	ProducerUUID   string         `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"` // is the producer UUID
	WriterUUID     string         `json:"writer_uuid"`                                                                 // is the remote consumer UUID, ie: whatever is subscribing to this producer
	DataStore      datatypes.JSON `json:"data_store"`
	WriterSettings datatypes.JSON `json:"producer_settings"` //like cov for a point or whatever is needed  #TODO this is why it needs settings
	CommonCreated
}

//Producer a producer is a placeholder to register an object to enable consumers to
// A producer for example is a point, Something that makes data, and the subscriber would have a consumer to it, Like grafana reading and writing to it from edge to cloud or wires over rest(peer to peer)
type Producer struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	//ProducerType        string `json:"producer_type"` //TODO REMOVE
	CommonCurrentProducer //if the point for example is read only the writer uuid would be the point uuid, ie: itself, so in this case there is no writer or writer clone
	CommonThing
	EnableHistory       bool              `json:"enable_history"`
	ProducerApplication string            `json:"producer_application"`
	StreamUUID          string            `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ProducerThingUUID   string            `json:"producer_thing_uuid"` //a producer_uuid is the point uuid
	PublishWithName     bool              `json:"publish_with_name"`   //publish with the point name and the type as an example TODO add these in for when we do MQTT
	PublishAttributes   bool              `json:"publish_attributes"`  //publish all fields from the producer WARNING this will increase network data TODO add these in for when we do MQTT
	WriterClone         []WriterClone     `json:"writer_clones" gorm:"constraint:OnDelete:CASCADE;"`
	ProducerHistory     []ProducerHistory `json:"producer_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

//ProducerHistory for storing the history
type ProducerHistory struct {
	CommonUUID
	ProducerUUID string `json:"producer_uuid" gorm:"TYPE:varchar(255) REFERENCES producers;not null;default:null"`
	CommonCurrentProducer
	DataStore datatypes.JSON `json:"data_store"`
	Timestamp time.Time      `json:"timestamp"`
}

//ProducerBody could be a local network, job or alarm and so on
type ProducerBody struct {
	CommonThing
	FlowNetworkUUID string      `json:"flow_network_uuid"`
	ProducerUUID    string      `json:"producer_uuid,omitempty"`
	StreamUUID      string      `json:"stream_uuid,omitempty"`
	Payload         interface{} `json:"payload"`
}
