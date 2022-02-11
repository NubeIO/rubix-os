package model

import "gorm.io/datatypes"

type CommonWriter struct {
	CommonUUID
	CommonSyncUUID
	WriterThingClass string         `json:"writer_thing_class,omitempty"`
	WriterThingType  string         `json:"writer_thing_type,omitempty"`
	WriterThingUUID  string         `json:"writer_thing_uuid,omitempty"`
	WriterThingName  string         `json:"writer_thing_name,omitempty"`
	DataStore        datatypes.JSON `json:"data_store,omitempty"`
	CommonCreated
}

//Writer could be a local network, job or alarm and so on...
type Writer struct {
	CommonWriter
	ConsumerUUID string `json:"consumer_uuid,omitempty" gorm:"TYPE:string REFERENCES consumers;not null;default:null"`
}

//WriterClone list of all the consumers
type WriterClone struct { //TODO the WriterClone needs to publish a COV event as for example we have 2x mqtt broker then the cov for a point maybe different when not going over the internet
	CommonWriter
	ProducerUUID      string `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"`
	FlowFrameworkUUID string `json:"flow_framework_uuid,omitempty"`
	CommonSourceUUID
}

type SyncWriter struct {
	Writer            Writer
	ProducerUUID      string
	FlowFrameworkUUID string
}

//WriterBody could be a local network, job or alarm and so on
type WriterBody struct {
	Action     string       `json:"action,omitempty"` //read, write and so on
	AskRefresh bool         `json:"ask_refresh,omitempty"`
	Priority   Priority     `json:"priority,omitempty"`
	Schedule   ScheduleData `json:"schedule,omitempty"`
}

//WriterBulk could be a local network, job or alarm and so on
type WriterBulk struct {
	WriterUUID  string      `json:"writer_uuid,omitempty"`
	Action      string      `json:"action,omitempty"` //read, write and so on
	AskRefresh  bool        `json:"ask_refresh,omitempty"`
	CommonValue CommonValue `json:"common_value,omitempty"`
	Priority    Priority    `json:"priority,omitempty"`
}

type SyncWriterCOV struct {
	WriterUUID      string
	OriginalValue   *float64
	PresentValue    *float64
	CurrentPriority *int
	Priority        *Priority
	Schedule        *datatypes.JSON
}
