package model

import "gorm.io/datatypes"

type ConsumersType struct {
	Network   		string `json:"network"`
	Job   			string `json:"job"`
	Point   		string `json:"point"`
	Alarm   		string `json:"alarm"`

}

type ConsumersUse struct {
	Local   		 string `json:"local"`
	Remote   		string `json:"remote"`
	Plugin   		string `json:"plugin"`
}



//Writer could be a local network, job or alarm and so on
type Writer struct {
	CommonUUID
	WriterType 					string  `json:"writer_type"` //point, schedule, job, network
	WriteCloneUUID 				string `json:"write_clone_uuid"`
	ConsumerUUID 				string `json:"consumer_uuid" gorm:"TYPE:string REFERENCES consumers;not null;default:null"`
	ConsumerThingUUID 			string `json:"consumer_thing_uuid"` // this is the consumer child point UUID
	DataStore 					datatypes.JSON  `json:"data_store"`
	WriterSettings 				datatypes.JSON  `json:"producer_settings"` //like cov for a point or whatever is needed  #TODO this is why it needs settings
	CommonCreated
}



//Consumer could be a local network, job or alarm and so on
type Consumer struct {
	CommonConsumer
	CurrentWriterCloneUUID  	string  `json:"current_writer_clone_uuid"` // this could come from any flow-network on any instance
	ProducerUUID  				string  `json:"producer_uuid"`
	ProducerThingUUID 			string 	`json:"producer_thing_uuid"` // this is the remote point UUID
	ConsumerType 				string  `json:"consumer_type"`
	ConsumerApplication 		string 	`json:"consumer_application"`
	StreamUUID     				string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	Writer						[]Writer `json:"writers" gorm:"constraint:OnDelete:CASCADE;"`
	DataStore 					datatypes.JSON  `json:"data_store"`
	ConsumerHistory				[]ConsumerHistory `json:"consumer_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

