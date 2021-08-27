package model



type SubscriptionsType struct {
	Network   		string `json:"network"`
	Job   			string `json:"job"`
	Point   		string `json:"point"`
	Alarm   		string `json:"alarm"`

}

type SubscriptionsApplication struct {
	Local   		 string `json:"local"`
	Remote   		string `json:"remote"`
	Plugin   		string `json:"plugin"`


}

//SubscriptionList could be a local network, job or alarm and so on
type SubscriptionList struct {
	CommonUUID
	SubscriptionUUID string `json:"subscription_uuid" gorm:"TYPE:string REFERENCES subscriptions;not null;default:null"`
	ToThingUUID 			 string 	`json:"to_thing_uuid"`

}



//Subscription could be a local network, job or alarm and so on
type Subscription struct {
	CommonSubscription
	SubscriptionType 			string  `json:"subscription_type"`
	SubscriptionApplication 	string `json:"subscription_application"`
	ToSubscriberUUID 			string `json:"to_subscriber_uuid"`
	StreamUUID     				string `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	SubscriptionList			[]SubscriptionList `json:"subscription_list" gorm:"constraint:OnDelete:CASCADE;"`
}

