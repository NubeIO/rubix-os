package model





//Subscriber could be a local network, job or alarm and so on
type Subscriber struct {
	CommonSubscriber
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	FromUUID 				string `json:"from_uuid"`
	ToUUID 					string `json:"to_uuid"`
	GatewayUUID     		string `json:"gateway_uuid" gorm:"TYPE:string REFERENCES gateways;not null;default:null"`
	//PointUUID    			string  `json:"point_uuid" binding:"required" gorm:"TYPE:varchar(255) REFERENCES points;unique;default:null"`
	//JobUUID    				string  `json:"job_uuid" binding:"required" gorm:"TYPE:varchar(255) REFERENCES jobs;unique;default:null"`
}

//type SubscriberNetwork struct {
//	Subscriber
//	SubscriberGlobalUUID  string `json:"subscriber_global_uuid"` //this would be the rubix uuid
//	SubscriberNetworkUUID string `json:"subscriber_network_uuid"` //this would be the network uuid
//	SubscriberDeviceUUID  string `json:"subscriber_device_uuid"` //this would be the device uuid
//	SubscriberPointUUID   string `json:"subscriber_point_uuid"` //this would be the point uuid
//
//}

