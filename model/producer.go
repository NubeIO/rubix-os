package model



//ProducerList list of all the subscriptions
type ProducerList struct {
	CommonUUID
	ProducerUUID 		string  `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"`
	ProducerThingUUID 	string `json:"subscription_uuid"`
	CommonCreated
}


//Producer a producer is a placeholder to register an object to enable subscriptions to
type Producer struct {
	CommonProducer
	ProducerType 			string  `json:"producer_type"`
	ProducerApplication 	string 	`json:"producer_application"`
	COV 					int 	`json:"cov"`
	//FromThingUUID 			string 	`json:"from_thing_uuid"`
	ProducerThingUUID 		string `json:"producer_thing_uuid"`
	StreamUUID     			string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ProducerList			[]ProducerList `json:"subscribers_list" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

