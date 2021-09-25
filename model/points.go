package model

import (
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

var ObjectTypes = struct {
	//bacnet
	AnalogInput  string
	AnalogOutput string
	AnalogValue  string
	BinaryInput  string
	BinaryOutput string
	BinaryValue  string
	//modbus
	ReadCoil           string
	ReadCoils          string
	ReadDiscreteInput  string
	ReadDiscreteInputs string
	WriteCoil          string
	WriteCoils         string
	ReadRegister       string
	ReadRegisters      string
	ReadInt16          string
	ReadSingleInt16    string
	WriteSingleInt16   string
	ReadUint16         string
	ReadSingleUint16   string
	WriteSingleUint16  string
	ReadInt32          string
	ReadSingleInt32    string
	WriteSingleInt32   string
	ReadUint32         string
	ReadSingleUint32   string
	WriteSingleUint32  string
	ReadFloat32        string
	ReadSingleFloat32  string
	WriteSingleFloat32 string
	ReadFloat64        string
	ReadSingleFloat64  string
	WriteSingleFloat64 string
}{
	//bacnet
	AnalogInput:  "analogInput",
	AnalogOutput: "analogOutput",
	AnalogValue:  "analogValue",
	BinaryInput:  "binaryInput",
	BinaryOutput: "binaryOutput",
	BinaryValue:  "binaryValue",
	//modbus
	ReadCoil:           "readCoil",
	ReadCoils:          "readCoils",
	ReadDiscreteInput:  "readDiscreteInput",
	ReadDiscreteInputs: "readDiscreteInputs",
	WriteCoil:          "writeCoil",
	WriteCoils:         "writeCoils",
	ReadRegister:       "readRegister",
	ReadRegisters:      "readRegisters",
	ReadInt16:          "readInt16",
	ReadSingleInt16:    "readSingleInt16",
	WriteSingleInt16:   "writeSingleInt16",
	ReadUint16:         "readUint16",
	ReadSingleUint16:   "readSingleUint16",
	WriteSingleUint16:  "writeSingleUint16",
	ReadInt32:          "readInt32",
	ReadSingleInt32:    "readSingleInt32",
	WriteSingleInt32:   "writeSingleInt32",
	ReadUint32:         "readUint32",
	ReadSingleUint32:   "readSingleUint32",
	WriteSingleUint32:  "writeSingleUint32",
	ReadFloat32:        "readFloat32",
	ReadSingleFloat32:  "readSingleFloat32",
	WriteSingleFloat32: "writeSingleFloat32",
	ReadFloat64:        "readFloat64",
	ReadSingleFloat64:  "readSingleFloat64",
	WriteSingleFloat64: "writeSingleFloat64",
}

var ObjectEncoding = struct {
	LebBew string //LITTLE_ENDIAN, HIGH_WORD_FIRST
	LebLew string
	BebLew string
	BebBew string
}{
	LebBew: "lebBew",
	LebLew: "lebLew",
	BebLew: "bebLew",
	BebBew: "bebBew",
}

//Point table
type Point struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	PresentValue         float64        `json:"present_value"` //point value, read only
	CurrentPriority      int            `json:"current_priority"`
	WriteValue           *float64       `json:"write_value"` //TODO add in logic if user writes to below priority 16
	ValueRaw             datatypes.JSON `json:"value_raw"`
	Fallback             float64        `json:"fallback"`
	DeviceUUID           string         `json:"device_uuid" gorm:"TYPE:string REFERENCES devices;not null;default:null"`
	EnableWriteable      *bool          `json:"writeable"`
	IsOutput             *bool          `json:"is_output"`
	BoolInvert           *bool          `json:"bool_invert"`
	COV                  float32        `json:"cov"`
	ObjectType           string         `json:"object_type"`    //binaryInput, coil, if type os input dont return the priority array  TODO decide if we just stick to bacnet object types, as a binaryOut is the sample as a coil in modbus
	AddressId            int            `json:"address_id"`     // for example a modbus address or bacnet address
	AddressOffset        int            `json:"address_offset"` // for example a modbus address offset
	ZeroMode             *bool          `json:"zero_mode"`
	AddressUUID          string         `json:"address_uuid"` // for example a droplet id (so a string)
	NextAvailableAddress *bool          `json:"use_next_available_address"`
	UnitsTo              string         //with take the unit and convert to, this would affect the presentValue and the original value will be stored in the raw
	Decimal              int
	Round                *float64
	InputMin             *float64
	InputMax             *float64
	ScaleMin             *float64
	ScaleMax             *float64
	CommonThingClass           //point
	CommonThingUse             //mapping, lora, api
	CommonThingType            //temp
	IsProducer           *bool `json:"is_producer"`
	IsConsumer           *bool `json:"is_consumer"`
	CommonFault
	ThingClass   string    `json:"thing_class"`
	ThingType    string    `json:"thing_type"`
	Kind         string    `json:"kind"`
	Unit         string    `json:"unit"`
	UnitImperial string    `json:"unit_imperial"`
	EquipType    string    `json:"equip_type"`
	Priority     *Priority `json:"priority" gorm:"constraint:OnDelete:CASCADE"`
	Tags         []*Tag    `json:"tags" gorm:"many2many:points_tags;constraint:OnDelete:CASCADE"`
}

type Priority struct {
	PointUUID string   `json:"point_uuid,omitempty" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	P1        *float64 `json:"_1"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	P2        *float64 `json:"_2"`
	P3        *float64 `json:"_3"`
	P4        *float64 `json:"_4"`
	P5        *float64 `json:"_5"`
	P6        *float64 `json:"_6"`
	P7        *float64 `json:"_7"`
	P8        *float64 `json:"_8"`
	P9        *float64 `json:"_9"`
	P10       *float64 `json:"_10"`
	P11       *float64 `json:"_11"`
	P12       *float64 `json:"_12"`
	P13       *float64 `json:"_13"`
	P14       *float64 `json:"_14"`
	P15       *float64 `json:"_15"`
	P16       *float64 `json:"_16"` //removed and added to the point to save one DB write
}
