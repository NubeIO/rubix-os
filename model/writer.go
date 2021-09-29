package model

import "gorm.io/datatypes"

//Writer could be a local network, job or alarm and so on
type Writer struct {
	CommonUUID
	WriterThingClass  string         `json:"writer_thing_class,omitempty"`
	WriterThingType   string         `json:"writer_thing_type,omitempty"`
	WriterThingUUID   string         `json:"writer_thing_uuid,omitempty"`
	CloneUUID         string         `json:"clone_uuid,omitempty"`
	ConsumerUUID      string         `json:"consumer_uuid,omitempty" gorm:"TYPE:string REFERENCES consumers;not null;default:null"`
	ConsumerThingUUID string         `json:"consumer_thing_uuid,omitempty"` // this is the consumer child point UUID or an API-WRITER
	DataStore         datatypes.JSON `json:"data_store,omitempty"`
	WriterSettings    datatypes.JSON `json:"producer_settings,omitempty"` //like cov for a point or whatever is needed  #TODO this is why it needs settings
	CommonCreated
}

//WriterBody could be a local network, job or alarm and so on
type WriterBody struct {
	Action           string   `json:"action,omitempty"` //read, write and so on
	AskRefresh       bool     `json:"ask_refresh,omitempty"`
	CommonThingClass          //point, job,
	CommonThingType           // for example api, rssi, voltage
	Priority         Priority `json:"priority,omitempty"`
	Point            Point    `json:"point,omitempty"`
}

//WriterBulk could be a local network, job or alarm and so on
type WriterBulk struct {
	WriterUUID       string      `json:"writer_uuid,omitempty"`
	Action           string      `json:"action,omitempty"` //read, write and so on
	AskRefresh       bool        `json:"ask_refresh,omitempty"`
	CommonThingClass             //point, job
	CommonThingType              // for example temp, rssi, voltage
	CommonValue      CommonValue `json:"common_value,omitempty"`
	Priority         Priority    `json:"priority,omitempty"`
}
