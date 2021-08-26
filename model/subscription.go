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


//Subscription could be a local network, job or alarm and so on
type Subscription struct {
	CommonSubscription
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	ToUUID 					string `json:"to_uuid"`
	IsRemote 				bool 	`json:"is_remote"`
	RemoteRubixUUID			string 	`json:"remote_rubix_uuid"`
	StreamUUID     			string `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;null;default:null"`
}

