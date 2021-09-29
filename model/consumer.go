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
	StreamUUID          string             `json:"stream_uuid,omitempty" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ProducerUUID        string             `json:"producer_uuid,omitempty"`
	ProducerThingName   string             `json:"producer_thing_name,omitempty"`
	ProducerThingUUID   string             `json:"producer_thing_uuid,omitempty"` // this is the remote point UUID
	ProducerThingClass  string             `json:"producer_thing_class,omitempty"`
	ProducerThingRef    string             `json:"producer_thing_reference,omitempty"`
	ProducerThingType   string             `json:"producer_thing_type,omitempty"`
	ConsumerApplication string             `json:"consumer_application,omitempty"`
	CurrentWriterUUID   string             `json:"current_writer_uuid,omitempty"` // this could come from any flow-network on any instance
	DataStore           datatypes.JSON     `json:"data_store,omitempty"`
	Writers             []*Writer          `json:"writers,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	ConsumerHistories   []*ConsumerHistory `json:"consumer_histories,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Tags                []*Tag             `json:"tags,omitempty" gorm:"many2many:consumers_tags;constraint:OnDelete:CASCADE"`
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
