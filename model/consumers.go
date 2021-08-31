package model


type ConsumersType struct {
	Network   		string `json:"network"`
	Job   			string `json:"job"`
	Point   		string `json:"point"`
	Alarm   		string `json:"alarm"`

}

type ConsumersApplication struct {
	Local   		 string `json:"local"`
	Remote   		string `json:"remote"`
	Plugin   		string `json:"plugin"`


}

//Writer could be a local network, job or alarm and so on
type Writer struct {
	CommonUUID
	PresentValue       			float64  `json:"present_value"` // for common use of points
	WriteValue       			float64  `json:"write_value"` // for common use of points
	WriteCloneUUID 				string `json:"write_clone_uuid"`
	ConsumerUUID 				string `json:"consumer_uuid" gorm:"TYPE:string REFERENCES consumers;not null;default:null"`
	ConsumerThingUUID 			string `json:"consumer_thing_uuid"` // this is the consumer child point UUID
	ConsumerCOV 				float64 `json:"consumer_cov"`
	CommonCreated
}



//Consumer could be a local network, job or alarm and so on
type Consumer struct {
	CommonConsumer
	PresentValue       			float64  `json:"present_value"` // for common use of points
	ProducerUUID  				string  `json:"producer_uuid"`
	ProducerThingUUID 			string 	`json:"producer_thing_uuid"` // this is the remote point UUID
	ConsumerType 				string  `json:"consumer_type"`
	ConsumerApplication 		string 	`json:"consumer_application"`
	StreamUUID     				string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	Writer						[]Writer `json:"writers" gorm:"constraint:OnDelete:CASCADE;"`
	ConsumerHistory				[]ConsumerHistory `json:"consumer_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

