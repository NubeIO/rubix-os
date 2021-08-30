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
	SubscriptionUUID 			string `json:"subscription_uuid" gorm:"TYPE:string REFERENCES subscriptions;not null;default:null"`
	SubscriptionThingUUID 		string `json:"subscription_thing_uuid"` // this is the subscription child point UUID
	PresentValue 				float64 `json:"present_value"` //these fields are support as points is the most common use case for histories
	WriteValue       			float64  `json:"write_value"` // for common use of points
	SubscriptionCOV 			float64 `json:"subscription_cov"`
	CommonCreated
}



//Subscription could be a local network, job or alarm and so on
type Subscription struct {
	CommonSubscription
	ProducerUUID  				string  `json:"producer_uuid"`
	ProducerThingUUID 			string `json:"producer_thing_uuid"` // this is the remote point UUID
	SubscriptionType 			string  `json:"subscription_type"`
	SubscriptionApplication 	string `json:"subscription_application"`
	StreamUUID     				string `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	SubscriptionList			[]SubscriptionList `json:"subscription_list" gorm:"constraint:OnDelete:CASCADE;"`
	SubscriptionHistory			[]SubscriptionHistory `json:"subscription_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

