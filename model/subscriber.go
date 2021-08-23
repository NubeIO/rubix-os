package model



//type SubscriberType struct {
//	Network   		string `json:"network"`
//	Job   			string `json:"job"`
//	Point   		string `json:"point"`
//	Alarm   		string `json:"alarm"`
//
//}
//
//type SubscriberApplication struct {
//	Local   		 string `json:"local"`
//	Mapping   		 string `json:"mapping"`
//	Remote   		string `json:"remote"`
//	Plugin   		string `json:"plugin"`
//
//
//}
//
//
//func NewSubscriberTypeEnum() *SubscriberType {
//	return &SubscriberType{
//		Network: 	"network",
//		Job:     	"job",
//		Point:     	"point",
//		Alarm:   	"alarm",
//	}
//}
//func NewSubscriberApplicationEnum() *SubscriberApplication {
//	return &SubscriberApplication{
//		Mapping: 	"mapping",
//	}
//}



//Subscriber could be a local network, job or alarm and so on
type Subscriber struct {
	CommonSubscriber
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	ThingUUID 				string `json:"thing_uuid"`
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

