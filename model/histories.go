package model

import "time"

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
	ProducerUUID    	string  	`json:"producer_uuid" gorm:"TYPE:varchar(255) REFERENCES producers;not null;default:null"`
	PresentValue 		float64   	`json:"present_value"` //these fields are support as points is the most common use case for histories
	SLWriteUUID 		string 		`json:"sl_write_uuid"` //writer UUID
	ValueRaw     		[]byte    	`json:"value_raw"`     //used as example modbus array [0, 11]
	ValueStore     		[]byte    	`json:"value_store"`  //used to story priority_array or any generic data
	Timestamp    		time.Time 	`json:"timestamp"`

}

//ConsumerHistory for storing the history
type ConsumerHistory struct {
	CommonUUID
	ConsumerUUID    string  	`json:"consumer_uuid" gorm:"TYPE:varchar(255) REFERENCES consumers;not null;default:null"`
	PresentValue 		float64   	`json:"present_value"` //these fields are support as points is the most common use case for histories
	WriteValue       	float64    	`json:"write_value"` // for common use of points
	ValueRaw     		[]byte    	`json:"value_raw"`     //used as example modbus array [0, 11]
	ValueStore     		[]byte    	`json:"value_store"`  //used to story priority_array or any generic data
	Timestamp    		time.Time 	`json:"timestamp"`
}
