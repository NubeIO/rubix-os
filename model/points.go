package model

import (
	"github.com/NubeIO/null"
	"gorm.io/datatypes"
)

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

var ObjectType = struct {
	analogInput  string
	analogOutput string
	analogValue  string
	binaryInput  string
	binaryOutput string
	binaryValue  string
}{
	analogInput:  "analogInput",
	analogOutput: "analogOutput",
	analogValue:  "analogValue",
	binaryInput:  "binaryInput",
	binaryOutput: "binaryOutput",
	binaryValue:  "binaryValue",
}

//Point table
type Point struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	PresentValue  float64        `json:"present_value"` //point value, read only
	WriteValue    null.Float     `json:"write_value"`   //TODO add in logic if user writes to below priority 16
	ValueRaw      datatypes.JSON `json:"value_raw"`
	Fallback      null.Float     `json:"fallback"`
	DeviceUUID    string         `json:"device_uuid" gorm:"TYPE:string REFERENCES devices;not null;default:null"`
	Writeable     bool           `json:"writeable"`
	Cov           float32        `json:"cov"`
	ObjectType    string         `json:"object_type"`    //binaryInput, coil, if type os input dont return the priority array  TODO decide if we just stick to bacnet object types, as a binaryOut is the sample as a coil in modbus
	AddressId     int            `json:"address_id"`     // for example a modbus address or bacnet address
	AddressOffset int            `json:"address_offset"` // for example a modbus address offset
	AddressUUID   string         `json:"address_uuid"`   // for example a droplet id (so a string)
	PointType     string         `json:"point_type"`     // for example temp, rssi, voltage
	IsProducer    bool           `json:"is_producer"`
	IsConsumer    bool           `json:"is_consumer"`
	CommonFault
	Priority Priority `json:"priority,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	CommonCreated
}

type Priority struct {
	PointUUID string     `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	P1        null.Float `json:"_1"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	P2        null.Float `json:"_2"`
	P3        null.Float `json:"_3"`
	P4        null.Float `json:"_4"`
	P5        null.Float `json:"_5"`
	P6        null.Float `json:"_6"`
	P7        null.Float `json:"_7"`
	P8        null.Float `json:"_8"`
	P9        null.Float `json:"_9"`
	P10       null.Float `json:"_10"`
	P11       null.Float `json:"_11"`
	P12       null.Float `json:"_12"`
	P13       null.Float `json:"_13"`
	P14       null.Float `json:"_14"`
	P15       null.Float `json:"_15"`
	//P16       null.Float `json:"_16"` //removed and added to the point to save one DB write

}
