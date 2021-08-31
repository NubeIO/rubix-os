package model

import (
	"gorm.io/datatypes"
	"time"
)

//HistorySettings point history settings TODO add in later
//type HistorySettings struct {
//	Type string //cov, interval, cov_interval
//	Duration int //15min
//	SizeLimit int //max amount of records to keep, the newest will override the oldest record
//
//}


//ProducerHistory for storing the history
type ProducerHistory struct {
	CommonUUID
	ProducerUUID    		string  	`json:"producer_uuid" gorm:"TYPE:varchar(255) REFERENCES producers;not null;default:null"`
	CurrentWriterCloneUUID  string  	`json:"current_writer_clone_uuid"`
	DataStore 			datatypes.JSON  `json:"data_store"`
	Timestamp    		time.Time 		`json:"timestamp"`

}

//ConsumerHistory for storing the history
type ConsumerHistory struct {
	CommonUUID
	ConsumerUUID    	string  		`json:"consumer_uuid" gorm:"TYPE:varchar(255) REFERENCES consumers;not null;default:null"`
	ProducerUUID    	string
	DataStore 			datatypes.JSON  `json:"data_store"`
	Timestamp    		time.Time 		`json:"timestamp"`
}
