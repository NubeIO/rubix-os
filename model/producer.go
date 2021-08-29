package model



//ProducerSubscriptionList list of all the subscriptions
// a subscription
type ProducerSubscriptionList struct {
	CommonUUID
	ProducerUUID 		string  `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"` // is the producer UUID
	SubscriptionUUID 	string 	`json:"subscription_uuid"`  // is the remote subscription UUID, ie: whatever is subscribing to this producer
	SubscriptionCOV 	float64 `json:"subscription_cov"`
	CommonCreated
}


//Producer a producer is a placeholder to register an object to enable subscriptions to
type Producer struct {
	CommonProducer
	ProducerType 			string  `json:"producer_type"`
	ProducerApplication 	string 	`json:"producer_application"`
	StreamUUID     			string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	ProducerThingUUID 		string  `json:"producer_thing_uuid"` //a producer_uuid is the point uuid
	ProducerSubscriptionList	[]ProducerSubscriptionList `json:"subscription_list" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

