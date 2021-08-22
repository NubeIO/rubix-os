package model

/*

Rubix V2

Network
- Any new network is called a ProducerNetwork for example modbus or lora
- Any ProducerNetwork doesn't need to use point-mapping to get the data into the cloud
- Any ProducerNetwork can be added to a one or more SubscriberNetworks, The ProducerNetwork will keep a leger or which Subscriber are reading/writing to its points

The SubscriberNetwork is the remote rubix device.
So when the ProducerNetwork (the producer network producers data ie: lora) network has a connection with the SubscriberNetwork the ProducerNetwork keeps a ledger of the SubscriberPoints

ProducerNetwork and SubscriberNetwork Jobs
- publish any CRUD updates to all subscribers (ie when a point is deleted or the name is updated)
- publish any COV events

ProducerNetwork settings
- COV will be set in the producer

SubscriberNetwork settings (these settings are not like 2-way meaning that in the SubscriberNetwork if the COV is updated it will not affect the ProducerNetwork setting)
- as this would be considered a normal point in the SubscriberNetwork this point will have all the same settings ie: history, cov and so on

CommandGroup
- is for issuing global schedule writes or global point writes (as in send a value to any point associated with this group)

TimeOverride
- where a point value can be overridden for a duration of time



REST calls
ProducerNetwork
- can call all attributes

*/



type SubscriberType struct {
	Network   		string `json:"network"`
	Job   			string `json:"job"`
	Point   		string `json:"point"`
	Alarm   		string `json:"alarm"`

}

type SubscriberApplication struct {
	Local   		 string `json:"local"`
	Mapping   		 string `json:"mapping"`
	Remote   		string `json:"remote"`
	Plugin   		string `json:"plugin"`


}


func NewSubscriberTypeEnum() *SubscriberType {
	return &SubscriberType{
		Network: 	"network",
		Job:     	"job",
		Point:     	"point",
		Alarm:   	"alarm",
	}
}
func NewSubscriberApplicationEnum() *SubscriberApplication {
	return &SubscriberApplication{
		Mapping: 	"mapping",
	}
}



//Subscriber could be a local network, job or alarm and so on
type Subscriber struct {
	CommonSubscriber
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	GatewayUUID     		string `json:"gateway_uuid" gorm:"TYPE:string REFERENCES gateways;not null;default:null"`
	PointUUID    			string  `json:"point_uuid" binding:"required" gorm:"TYPE:varchar(255) REFERENCES points;unique;default:null"`
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

