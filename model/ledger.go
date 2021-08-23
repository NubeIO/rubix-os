package model

//PointSubscriberLedger record holder for subscription/subscriber
type PointSubscriberLedger struct {
	UUID	string 	`json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
	PointUUID    		string  `json:"point_uuid" binding:"required" gorm:"TYPE:varchar(255) REFERENCES points;not null;default:null"`
	GatewayUUID  		string 	`json:"gateway_uuid"`
	SubscriberUUID  	string 	`json:"subscriber_uuid"`
	SubscriptionUUID  	string 	`json:"subscription_uuid"`
}

//PointSubscriptionLedger record holder for subscription/subscriber
type PointSubscriptionLedger struct {
	UUID	string 	`json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
	PointUUID    		string  `json:"point_uuid" binding:"required" gorm:"TYPE:varchar(255) REFERENCES points;not null;default:null"`
	GatewayUUID  		string 	`json:"gateway_uuid"`
	SubscriberUUID  	string 	`json:"subscriber_uuid"`
}
