package client

type Token struct {
	ID  	int  	`json:"id"`
	Token  	string 	`json:"token"`
	UUID  	string 	`json:"uuid"`
}


type ResponseBody struct {
	Response ResponseCommon 	`json:"response"`
	Status     string 			`json:"status"`
	Count     string 			`json:"count"`
}

type ResponseCommon struct {

	UUID  			string 	`json:"uuid"`
	Name  			string 	`json:"name"`
	NetworkUUID  	string 	`json:"network_uuid"`
	DeviceUUID  	string 	`json:"device_uuid"`
	PointUUID  		string 	`json:"point_uuid"`
	GatewayUUID  	string 	`json:"gateway_uuid"`


}


type Gateway struct {
	Name  		string 	`json:"name"`
	IsRemote  	bool 	`json:"is_remote"`
}
type Subscriber struct {
	Name                  	string `json:"name"`
	Enable                	bool   `json:"enable"`
	SubscriberType        	string `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	GatewayUuid           	string `json:"gateway_uuid"`
	FromUUID 				string `json:"from_uuid"`
	ToUUID 					string `json:"to_uuid"`
}


type Subscription struct {
	Name                  	string `json:"name"`
	Enable                	bool   `json:"enable"`
	SubscriberType        	string `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	GatewayUuid           	string `json:"gateway_uuid"`
	ToUUID 					string `json:"to_uuid"`
}
