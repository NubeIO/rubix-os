package model

/*
The SubscriberNetwork is the remote rubix device.
So when the ProducerNetwork (the producer network producers data ie: lora) network has a connection with the SubscriberNetwork the ProducerNetwork keeps a ledger of the SubscriberPoints

ProducerNetwork and SubscriberNetwork Jobs
- publish any CRUD updates to all subscribers (ie when a point is deleted or the name is updated)
- publish any COV events

ProducerNetwork settings
- COV will be set in the producer

SubscriberNetwork settings (these settings are not like 2-way meaning that in the SubscriberNetwork if the COV is updated it will not affect the ProducerNetwork setting)
- as this would be considered a normal point in the SubscriberNetwork this point will have all the same settings ie: history, cov and so on

REST calls
ProducerNetwork
- can call all attributes

 */

//SubscriberNetwork is a remote device that can pub/sub to a producer
type SubscriberNetwork struct {
	SubscriberGlobalUUID  string `json:"subscriber_global_uuid"` //this would be the rubix uuid
	SubscriberNetworkUUID string `json:"subscriber_network_uuid"` //this would be the network uuid
	SubscriberDeviceUUID  string `json:"subscriber_device_uuid"` //this would be the device uuid
	SubscriberPointUUID   string `json:"subscriber_point_uuid"` //this would be the point uuid
}



//SerialNetwork type serial
type SerialNetwork struct {
	Port     string `json:"port"` //dev/tty/USB0
	BaudRate int `json:"baud_rate"` //9600
	StopBits int `json:"stop_bits"`
	Parity   int `json:"parity"`
	DataBits int `json:"data_bits"`
	Timeout  int `json:"timeout"`

}

type IPType struct {
	REST  	 string `json:"rest"`
	UDP     string `json:"udp"`
	MQTT     string `json:"mqtt"`

}



//IPNetwork type ip based network
type IPNetwork struct {
	IP  	 	string `json:"ip"`
	Port     	string `json:"port"`
	User     	string `json:"user"`
	Password    string `json:"password"`
	Token 		string `json:"token"`
	IPType

}


type Network struct {
	CommonUUID
	CommonNameUnique
	Common
	Created
	Manufacture 	string `json:"manufacture"`
	Model 			string `json:"model"`
	NetworkType		string `json:"network_type"`
	Device 			[]Device `json:"devices" gorm:"constraint:OnDelete:CASCADE;"`
}

