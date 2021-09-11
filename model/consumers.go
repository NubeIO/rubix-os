package model

import (
	"gorm.io/datatypes"
	"time"
)

//Writer could be a local network, job or alarm and so on
type Writer struct {
	CommonUUID
	WriterType        string         `json:"writer_type"` //point, schedule, job, network
	WriteCloneUUID    string         `json:"write_clone_uuid"`
	ConsumerUUID      string         `json:"consumer_uuid" gorm:"TYPE:string REFERENCES consumers;not null;default:null"`
	ConsumerThingUUID string         `json:"consumer_thing_uuid"` // this is the consumer child point UUID
	DataStore         datatypes.JSON `json:"data_store"`
	WriterSettings    datatypes.JSON `json:"producer_settings"` //like cov for a point or whatever is needed  #TODO this is why it needs settings
	CommonCreated
}

//Consumer could be a local network, job or alarm and so on
type Consumer struct {
	CommonConsumer
	CurrentWriterCloneUUID string            `json:"current_writer_clone_uuid"` // this could come from any flow-network on any instance
	ProducerUUID           string            `json:"producer_uuid"`
	ProducerThingUUID      string            `json:"producer_thing_uuid"` // this is the remote point UUID
	ThingType              string            `json:"thing_type"`
	ConsumerThingType      string            `json:"consumer_thing_type"`
	ConsumerApplication    string            `json:"consumer_application"`
	StreamUUID             string            `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	DataStore              datatypes.JSON    `json:"data_store"`
	Writer                 []Writer          `json:"writers" gorm:"constraint:OnDelete:CASCADE;"`
	ConsumerHistory        []ConsumerHistory `json:"consumer_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

//ConsumerHistory for storing the history
type ConsumerHistory struct {
	CommonUUID
	ConsumerUUID string `json:"consumer_uuid" gorm:"TYPE:varchar(255) REFERENCES consumers;not null;default:null"`
	ProducerUUID string
	DataStore    datatypes.JSON `json:"data_store"`
	Timestamp    time.Time      `json:"timestamp"`
}

//WriterBody could be a local network, job or alarm and so on
type WriterBody struct {
	Action      string      `json:"action"` //read, write and so on
	AskRefresh  bool        `json:"ask_refresh"`
	CommonValue CommonValue `json:"common_value"`
	Priority    Priority    `json:"priority"`
	Point       Point       `json:"point"`
}

//WriterBulk could be a local network, job or alarm and so on
type WriterBulk struct {
	WriterUUID  string      `json:"writer_uuid"`
	Action      string      `json:"action"` //read, write and so on
	AskRefresh  bool        `json:"ask_refresh"`
	CommonValue CommonValue `json:"common_value"`
	Priority    Priority    `json:"priority"`
}
