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
	ProducerUUID    	string  `json:"producer_uuid" gorm:"TYPE:varchar(255) REFERENCES producers;not null;default:null"`
	SubscriptionUUID    string  `json:"subscription_uuid"` // to track which subscription wrote the current value
	PresentValue 		float64   `json:"present_value"` //point value
	ValueRaw     		[]byte    `json:"value_raw"`     //modbus array [0, 11]
	CreatedAt    		time.Time `json:"created_on"`
}

//SubscriptionHistory for storing the history
type SubscriptionHistory struct {
	CommonUUID
	SubscriptionUUID    string  `json:"subscription_uuid" gorm:"TYPE:varchar(255) REFERENCES subscriptions;not null;default:null"`
	PresentValue 		float64   `json:"present_value"` //point value
	ValueRaw     		[]byte    `json:"value_raw"`     //modbus array [0, 11]
	CreatedAt    		time.Time `json:"created_on"`
}
