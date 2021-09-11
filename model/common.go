package model

import (
	"github.com/NubeIO/null"
	"time"
)

type CommonDescription struct {
	Description string `json:"description,omitempty"`
}

type CommonName struct {
	Name string `json:"name"`
}

type CommonThing struct {
	ThingType string `json:"thing_type"`
}

type CommonNameUnique struct {
	Name string `json:"name"  gorm:"type:varchar(255);unique;not null"`
}

type CommonModulePath struct {
	ModulePath string `json:"module_path"  gorm:"type:varchar(255);unique;not null"`
}

type CommonHelp struct {
	Help string `json:"help"`
}

type CommonNodeType struct {
	NodeType string `json:"node_type"`
}

type CommonType struct {
	ObjectType string `json:"object_type"`
}

type CommonAction struct {
	Action string `json:"action"`
}

type CommonEnable struct {
	Enable bool `json:"enable"`
}

type CommonID struct {
	ID string `json:"id"`
}

type CommonIDUnique struct {
	Name string `json:"id"  gorm:"type:varchar(255);unique;not null"`
}

type CommonUUID struct {
	UUID string `json:"uuid" sql:"uuid"  gorm:"type:varchar(255);unique;primaryKey"`
}

type CommonRubixUUID struct {
	RubixUUID string `json:"rubix_uuid"`
}

type CommonCreated struct {
	CreatedAt time.Time `json:"created_on"`
	UpdatedAt time.Time `json:"updated_on"`
}

type CommonHistory struct {
	EnableHistory bool `json:"history_enable"`
}

type CommonValue struct {
	Value    null.Float `json:"value"`
	ValueRaw string     `json:"value_raw"`
}

type CommonFault struct {
	InFault      bool      `json:"fault"`
	MessageLevel string    `json:"message_level"`
	MessageCode  string    `json:"message_code"`
	Message      string    `json:"message"`
	LastOk       time.Time `json:"last_ok"`
	LastFail     time.Time `json:"last_fail"`
}

type CommonIP struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	HTTP  bool   `json:"http"`
	HTTPS bool   `json:"https"`
}

type CommonStore struct {
	CommonValue
	CommonFault
}

type CommonProducerPermissions struct {
	Blacklist bool `json:"blacklist"`
	ReadOnly  bool `json:"read_only"`
	AllowCRUD bool `json:"allow_crud"` //not sure if this will be used, but it will allow the producer to update the producer
}

type CommonCurrentProducer struct {
	ThingWriterUUID string `json:"thing_writer_uuid"`
}

type CommonConsumer struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
}
