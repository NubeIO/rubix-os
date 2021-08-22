package model




// Ops TODO add in later
//Ops Means operations supported by a network, device, point and so on (example point supports point-write)
type Ops struct {

}

// SubscriberPoints TODO add in later
//SubscriberPoints is the list of points that this points is subscribed to. A Subscriber is a remote device that can pub/sub to a producer.
//type SubscriberPoints struct {
//	SubscriberConnectionUUID 				string `json:"subscriber_connection_uuid"` //this is the remote device UUID
//	SubscriberPointUUID						string `json:"subscriber_point_uuid"` //this is the remote point UUID
//	SubscriberWriteable 					bool `json:"subscriber_writeable"`
//	SubscriberPointWriteValue  				string `json:"subscriber_point_write_value"`
//	SubscriberPointWritePriority  			string `json:"subscriber_point_write_priority"`
//}

// TimeOverride TODO add in later
//TimeOverride where a point value can be overridden for a duration of time
type TimeOverride struct {
	PointUUID     		string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	StartDate        	string `json:"start_date"` // START at 25:11:2021:13:00
	EndDate        		string `json:"end_date"` // START at 25:11:2021:13:30
	Value				string `json:"value"`
	Priority			string `json:"priority"`

}

// CommandGroup TODO add in later
//CommandGroup is for issuing global schedule writes or global point writes (as in send a value to any point associated with this group)
type CommandGroup struct {
	PointUUID     	string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	//Points  []Point //list of the points in the command group, the point must be a writable point
}



//MathOperation same as in lora and point-server TODO add in later
type MathOperation struct {
	Calc string  //x + 1
	X float64
}


//Scale point value limits TODO add in later
type Scale struct {
	High float64
	Low float64
}

//Units this will be for point value conversion TODO add in later
type Units struct { // for example from temp c to temp f
     From string //https://github.com/martinlindhe/unit
     To string
}


//CommonPoint if a point is writable or not
type CommonPoint struct {
	Writeable 		bool   `json:"writeable"`
	Cov  			float64 `json:"cov"`
	ObjectType    	string `json:"object_type"`
	FallbackValue 	float64 `json:"fallback_value"` //is nullable
}

//Point table
type Point struct {
	CommonProducer
	CommonPoint
	DeviceUUID     		string `json:"device_uuid" gorm:"TYPE:string REFERENCES devices;not null;default:null"`
	Subscriber			[]Subscriber `json:"subscriber" gorm:"constraint:OnDelete:CASCADE;"`

	//Subscriptions		[]Subscriptions `json:"subscriptions" gorm:"constraint:OnDelete:CASCADE;"`
	//PriorityArrayModel 			PriorityArrayModel `json:"priority_array" gorm:"constraint:OnDelete:CASCADE"`
	//PointStore 					PointStore `json:"point_store" gorm:"constraint:OnDelete:CASCADE"`
}

//type PointSubscriber struct {
//	PointUUID    		string  `json:"point_uuid" binding:"required" gorm:"TYPE:varchar(255) REFERENCES points;not null;default:null"`
//	Subscriber
//}


type PriorityArrayModel struct {
	PointUUID     	string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	CommonCreated
	P1  			string `json:"_1"` //would be better if we stored the TS and where it was written from, for example from a Remote Subscriber
	P2  			string `json:"_2"`
	P3  			string `json:"_3"`
	P4  			string `json:"_4"`
	P5  			string `json:"_5"`
	P6  			string `json:"_6"`
	P7  			string `json:"_7"`
	P8  			string `json:"_8"`
	P9  			string `json:"_9"`
	P10  			string `json:"_10"`
	P11  			string `json:"_11"`
	P12  			string `json:"_12"`
	P13  			string `json:"_13"`
	P14  			string `json:"_14"`
	P15  			string `json:"_15"`
	P16  			string `json:"_16"`
}
