package model





//Subscriber could be a local network, job or alarm and so on
type Subscriber struct {
	CommonSubscriber
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string 	`json:"subscriber_application"`
	COV 					int 	`json:"cov"`
	FromUUID 				string 	`json:"from_thing_uuid"`
	SubscriptionUUID 		string 	`json:"to_subscription_uuid"`
	StreamUUID     			string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
}

