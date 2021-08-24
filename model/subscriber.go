package model





//Subscriber could be a local network, job or alarm and so on
type Subscriber struct {
	CommonSubscriber
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	FromUUID 				string `json:"from_uuid"`
	ToUUID 					string 	`json:"to_uuid"`
	IsRemote 				bool 	`json:"is_remote"`
	RemoteRubixUUID			string 	`json:"remote_rubix_uuid"`
	GatewayUUID     		string `json:"gateway_uuid" gorm:"TYPE:string REFERENCES gateways;not null;default:null"`
}

