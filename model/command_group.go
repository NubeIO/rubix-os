package model




// CommandGroup TODO add in later
//CommandGroup is for issuing global schedule writes or global point writes (as in send a value to any point associated with this group)
type CommandGroup struct {
	CommonUUID
	CommonName
	CommonEnable
	CommonDescription
	CommandUse 				string  `json:"command_use"`  //common uses will be point write to many points, master schedules or schedule grouping
	GatewayUUID     		string 	`json:"gateway_uuid" gorm:"TYPE:string REFERENCES gateways;not null;default:null"`
	WriteValue    			string 	`json:"write_value"`
	WritePriority 			string  `json:"write_priority"` //TODO maybe remove this and just use the writeJSON as we need things like Schedules aswell
	WritePriorityArray 		string  `json:"write_priority_array"`	//TODO maybe remove this and just use the writeJSON as we need things like Schedules aswell
	WriteJSON     			string 	`json:"write_json"` //TODO add data model in later
	StartDate        		string 	`json:"start_date"` // START at 25:11:2021:13:00
	EndDate        			string 	`json:"end_date"` // START at 25:11:2021:13:30
	Value					string 	`json:"value"`
	Priority				string 	`json:"priority"`
	CommonCreated
}


type WritePriorityArray struct {
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



