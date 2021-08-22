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

}

