package model

import (
	"gorm.io/datatypes"
	"time"
)

//ProducerHistory for storing the history
type ProducerHistory struct {
	ID           int    `json:"id" gorm:"AUTO_INCREMENT;primary_key;index"`
	ProducerUUID string `json:"producer_uuid" gorm:"TYPE:varchar(255) REFERENCES producers;not null;default:null"`
	CommonCurrentProducer
	DataStore  datatypes.JSON `json:"data_store"`
	Timestamp  time.Time      `json:"timestamp"`
	WriterUUID string         `json:"writer_uuid,omitempty"`
}
