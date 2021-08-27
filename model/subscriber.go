package model



//SubscriberList list of all the subscriptions
type SubscriberList struct {
	CommonUUID
	SubscriberUUID 		string  `json:"subscriber_list" gorm:"TYPE:string REFERENCES subscribers;not null;default:null"`
	FromThingUUID 		string `json:"from_thing_uuid"`

}


//Subscriber a subscriber is a placeholder to register an object to enable subscriptions to
type Subscriber struct {
	CommonSubscriber
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string 	`json:"subscriber_application"`
	COV 					int 	`json:"cov"`
	FromThingUUID 			string 	`json:"from_thing_uuid"`
	SubscriptionUUID 		string 	`json:"to_subscription_uuid"`
	StreamUUID     			string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	SubscriberList			[]SubscriberList `json:"subscriber_list" gorm:"constraint:OnDelete:CASCADE;"`

}

