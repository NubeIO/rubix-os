package model



//SubscriberList list of all the subscriptions
// a subscription
type SubscriberList struct {
	CommonUUID
	ProducerUUID 		string  `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"` // is the producer UUID
	SubscriptionUUID 	string `json:"subscription_uuid"`  // is the remote subscription UUID, ie: whatever is subscribing to this producer
	CommonCreated
}


//Producer a producer is a placeholder to register an object to enable subscriptions to
type Producer struct {
	CommonProducer
	ProducerType 			string  `json:"producer_type"`
	ProducerApplication 	string 	`json:"producer_application"`
	COV 					float64 	`json:"cov"`
	StreamUUID     			string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ProducerThingUUID 		string `json:"producer_thing_uuid"` //a producer_uuid is the point uuid
	SubscriberList			[]SubscriberList `json:"subscribers_list" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

