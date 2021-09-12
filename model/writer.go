package model

import "gorm.io/datatypes"

//Writer could be a local network, job or alarm and so on
type Writer struct {
	CommonUUID
	CommonThingClass
	CommonThingType
	CommonThingUUID
	CloneUUID         string         `json:"clone_uuid"`
	ConsumerUUID      string         `json:"consumer_uuid" gorm:"TYPE:string REFERENCES consumers;not null;default:null"`
	ConsumerThingUUID string         `json:"consumer_thing_uuid"` // this is the consumer child point UUID or an API-WRITER
	DataStore         datatypes.JSON `json:"data_store"`
	WriterSettings    datatypes.JSON `json:"producer_settings"` //like cov for a point or whatever is needed  #TODO this is why it needs settings
	CommonCreated
}

//WriterBody could be a local network, job or alarm and so on
type WriterBody struct {
	Action           string   `json:"action"` //read, write and so on
	AskRefresh       bool     `json:"ask_refresh"`
	CommonThingClass          //point, job
	CommonThingType           // for example temp, rssi, voltage
	Priority         Priority `json:"priority"`
	Point            Point    `json:"point"`
}

//WriterBulk could be a local network, job or alarm and so on
type WriterBulk struct {
	WriterUUID       string      `json:"writer_uuid"`
	Action           string      `json:"action"` //read, write and so on
	AskRefresh       bool        `json:"ask_refresh"`
	CommonThingClass             //point, job
	CommonThingType              // for example temp, rssi, voltage
	CommonValue      CommonValue `json:"common_value"`
	Priority         Priority    `json:"priority"`
}
