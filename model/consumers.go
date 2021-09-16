package model

import (
	"gorm.io/datatypes"
	"time"
)

//Consumer could be a local network, job or alarm and so on
type Consumer struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	StreamUUID          string             `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ProducerUUID        string             `json:"producer_uuid"`
	ProducerThingName   string             `json:"producer_thing_name"`
	ProducerThingUUID   string             `json:"producer_thing_uuid"` // this is the remote point UUID
	ProducerThingClass  string             `json:"producer_thing_class"`
	ProducerThingType   string             `json:"producer_thing_type"`
	ConsumerApplication string             `json:"consumer_application"`
	CurrentWriterUUID   string             `json:"current_writer_uuid"` // this could come from any flow-network on any instance
	DataStore           datatypes.JSON     `json:"data_store"`
	Writers             []*Writer          `json:"writers" gorm:"constraint:OnDelete:CASCADE;"`
	ConsumerHistories   []*ConsumerHistory `json:"consumer_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

//ConsumerHistory for storing the history
type ConsumerHistory struct {
	CommonUUID
	ConsumerUUID string         `json:"consumer_uuid" gorm:"TYPE:varchar(255) REFERENCES consumers;not null;default:null"`
	ProducerUUID string         `json:"producer_uuid"`
	DataStore    datatypes.JSON `json:"data_store"`
	Timestamp    time.Time      `json:"timestamp"`
}
