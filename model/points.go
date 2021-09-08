package model

import "github.com/NubeIO/null"

// Ops TODO add in later
//Ops Means operations supported by a network, device, point and so on (example point supports point-write)
type Ops struct {
}

// TimeOverride TODO add in later
//TimeOverride where a point value can be overridden for a duration of time
type TimeOverride struct {
	PointUUID string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	StartDate string `json:"start_date"` // START at 25:11:2021:13:00
	EndDate   string `json:"end_date"`   // START at 25:11:2021:13:30
	Value     string `json:"value"`
	Priority  string `json:"priority"`
}

//MathOperation same as in lora and point-server TODO add in later
type MathOperation struct {
	Calc string //x + 1
	X    float64
}

//Scale point value limits TODO add in later
type Scale struct {
	High float64
	Low  float64
}

//Units this will be for point value conversion TODO add in later
type Units struct { // for example from temp c to temp f
	From string //https://github.com/martinlindhe/unit
	To   string
}

//CommonPoint if a point is writable or not
type CommonPoint struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	PresentValue  float64 `json:"present_value"` //point value
	WriteValue    float64 `json:"write_value"`   //TODO add in logic if user writes to below priority 16
	ValueRaw      []byte  `json:"value_raw"`     //modbus array [0, 11]
	Fallback      float64 `json:"fallback"`
	Writeable     bool    `json:"writeable"`
	Cov           float64 `json:"cov"`
	ObjectType    string  `json:"object_type"`    //binaryInput, coil, if type os input dont return the priority array  TODO decide if we just stick to bacnet object types, as a binaryOut is the sample as a coil in modbus
	AddressId     int     `json:"address_id"`     // for example a modbus address or bacnet address
	AddressOffset int     `json:"address_offset"` // for example a modbus address offset
	AddressCode   string  `json:"address_code"`   // for example a droplet id (so a string)
	PointType     string  `json:"point_type"`     // for example temp, rssi, voltage

}

//Point table
type Point struct {
	CommonPoint
	DeviceUUID string `json:"device_uuid" gorm:"TYPE:string REFERENCES devices;not null;default:null"`
	CommonCreated
	Priority Priority `json:"priority" gorm:"constraint:OnDelete:CASCADE"`
}

type Priority struct {
	PointUUID string     `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	P1        null.Float `json:"_1"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	P2        null.Float `json:"_2"`
	P3        null.Float `json:"_3"`
	//P4  			float64 `json:"_4"`
	//P5  			float64 `json:"_5"`
	//P6  			float64 `json:"_6"`
	//P7  			float64 `json:"_7"`
	//P8  			float64 `json:"_8"`
	//P9  			float64 `json:"_9"`
	//P10  			float64 `json:"_10"`
	//P11  			float64 `json:"_11"`
	//P12  			float64 `json:"_12"`
	//P13  			float64 `json:"_13"`
	//P14  			float64 `json:"_14"`
	//P15  			float64 `json:"_15"`
	//P16  			float64 `json:"_16"` removed and added to the point to save one DB write
	//CommonCreated
}
