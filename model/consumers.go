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
	CurrentWriterUUID   string            `json:"current_writer_uuid"` // this could come from any flow-network on any instance
	ProducerUUID        string            `json:"producer_uuid"`
	ProducerThingUUID   string            `json:"producer_thing_uuid"` // this is the remote point UUID
	ProducerThingClass  string            `json:"thing_class"`
	ProducerThingType   string            `json:"thing_type"`
	StreamUUID          string            `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ConsumerApplication string            `json:"consumer_application"`
	DataStore           datatypes.JSON    `json:"data_store"`
	Writer              []Writer          `json:"writers" gorm:"constraint:OnDelete:CASCADE;"`
	ConsumerHistory     []ConsumerHistory `json:"consumer_histories" gorm:"constraint:OnDelete:CASCADE;"`
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
