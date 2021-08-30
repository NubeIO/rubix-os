package client

import "github.com/NubeDev/flow-framework/model"

type Token struct {
	ID  	int  	`json:"id"`
	Token  	string 	`json:"token"`
	UUID  	string 	`json:"uuid"`
}


type ResponsePlugins struct {
	Response 	Plugins 	`json:"response"`
	Status     	string 			`json:"status"`
	Count     	int 			`json:"count"`
}

type ResponsePoint struct {
	Points 		Points 			`json:"response"`
}

type ResponseProducer struct {
	Producer 		model.Producer 			`json:"response"`
}

//type Producer struct {
//	PresentValue 			float64   `json:"present_value"` //these fields are support as points is the most common use case for histories
//	WriteValue       		float64    `json:"write_value"` // for common use of points
//	ProducerType 			string  `json:"producer_type"`
//	ProducerApplication 	string 	`json:"producer_application"`
//	StreamUUID     			string 	`json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
//	ProducerThingUUID 		string  `json:"producer_thing_uuid"` //a producer_uuid is the point uuid
//}


type Points struct {
	model.CommonUUID
	model.CommonName
	model.CommonDescription
	model.CommonEnable
	DeviceUUID     	string `json:"device_uuid"`
	Writeable 		bool   `json:"writeable"`
	Cov  			float64 `json:"cov"`
	ObjectType    	string `json:"object_type"`
	FallbackValue 	float64 `json:"fallback_value"` //is nullable
}
type Plugins struct {
	Items []model.PluginConf
}


type ResponseBody struct {
	Response 	ResponseCommon 	`json:"response"`
	Status     	string 			`json:"status"`
	Count     	string 			`json:"count"`
}

type ResponseCommon struct {
	UUID  			string 	`json:"uuid"`
	Name  			string 	`json:"name"`
	NetworkUUID  	string 	`json:"network_uuid"`
	DeviceUUID  	string 	`json:"device_uuid"`
	PointUUID  		string 	`json:"point_uuid"`
	StreamUUID  	string 	`json:"stream_uuid"`
	GlobalUUID  	string 	`json:"global_uuid"`

}


type Stream struct {
	Name  		string 	`json:"name"`
	IsRemote  	bool 	`json:"is_remote"`
}


type Subscription struct {
	Name                  	string `json:"name"`
	Enable                	bool   `json:"enable"`
	ProducerType        	string `json:"producer_type"`
	ProducerApplication 	string `json:"producer_application"`
	StreamUUID  	string 	`json:"stream_uuid"`
	ToUUID 					string `json:"to_uuid"`
	IsRemote 				bool 	`json:"is_remote"`
	RemoteRubixUUID			string 	`json:"remote_rubix_uuid"`
}
