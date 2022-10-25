package pgmodel

import (
	"gorm.io/datatypes"
	"time"
)

type CommonUUID struct {
	UUID string `json:"uuid" sql:"uuid" gorm:"type:varchar(255);unique;primaryKey"`
}

type CommonName struct {
	Name string `json:"name"`
}

type CommonNameUnique struct {
	Name string `json:"name"  gorm:"type:varchar(255);unique;not null"`
}

type CommonDescription struct {
	Description string `json:"description,omitempty"`
}

type CommonEnable struct {
	Enable *bool `json:"enable"`
}

type CommonCreated struct {
	CreatedAt time.Time `json:"created_on,omitempty"`
	UpdatedAt time.Time `json:"updated_on,omitempty"`
}

type CommonFault struct {
	InFault      bool      `json:"fault,omitempty"`
	MessageLevel string    `json:"message_level,omitempty"`
	MessageCode  string    `json:"message_code,omitempty"`
	Message      string    `json:"message,omitempty"`
	LastOk       time.Time `json:"last_ok,omitempty"`
	LastFail     time.Time `json:"last_fail,omitempty"`
}

type CommonWriter struct {
	CommonUUID
	WriterThingClass string         `json:"writer_thing_class,omitempty"`
	WriterThingType  string         `json:"writer_thing_type,omitempty"`
	WriterThingUUID  string         `json:"writer_thing_uuid,omitempty"`
	WriterThingName  string         `json:"writer_thing_name,omitempty"`
	DataStore        datatypes.JSON `json:"data_store,omitempty"`
	CommonCreated
}

type FlowNetworkClone struct {
	CommonUUID
	CommonName
	CommonDescription
	GlobalUUID string `json:"global_uuid,omitempty"`
	ClientId   string `json:"client_id,omitempty"`
	ClientName string `json:"client_name,omitempty"`
	SiteId     string `json:"site_id,omitempty"`
	SiteName   string `json:"site_name,omitempty"`
	DeviceId   string `json:"device_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	CommonCreated
	StreamClones []*StreamClone `json:"stream_clones" gorm:"constraint:OnDelete:CASCADE;"`
}

type StreamClone struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonCreated
	FlowNetworkCloneUUID string      `json:"flow_network_clone_uuid" gorm:"TYPE:string REFERENCES flow_network_clones;not null;default:null"`
	Consumers            []*Consumer `json:"consumers" gorm:"constraint:OnDelete:CASCADE;"`
	Tags                 []*Tag      `json:"tags" gorm:"many2many:stream_clones_tags;constraint:OnDelete:CASCADE"`
}

// Consumer could be a local network, job or alarm and so on
type Consumer struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	ProducerUUID       string    `json:"producer_uuid,omitempty"`
	ProducerThingName  string    `json:"producer_thing_name,omitempty"`
	ProducerThingUUID  string    `json:"producer_thing_uuid,omitempty"` // this is the remote point UUID
	ProducerThingClass string    `json:"producer_thing_class,omitempty"`
	ProducerThingType  string    `json:"producer_thing_type,omitempty"`
	ProducerThingRef   string    `json:"producer_thing_ref,omitempty"`
	StreamCloneUUID    string    `json:"stream_clone_uuid,omitempty" gorm:"TYPE:string REFERENCES stream_clones;not null;default:null"`
	Writers            []*Writer `json:"writers,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Tags               []*Tag    `json:"tags,omitempty" gorm:"many2many:consumers_tags;constraint:OnDelete:CASCADE"`
	CommonCreated
}

type Writer struct {
	CommonWriter
	ConsumerUUID string   `json:"consumer_uuid,omitempty" gorm:"TYPE:string REFERENCES consumers;not null;default:null"`
	PresentValue *float64 `json:"present_value,omitempty"`
}

type Network struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	Devices []*Device `json:"devices,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Tags    []*Tag    `json:"tags,omitempty" gorm:"many2many:networks_tags;constraint:OnDelete:CASCADE"`
}

type Device struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	NetworkUUID string   `json:"network_uuid,omitempty" gorm:"TYPE:varchar(255) REFERENCES networks;not null;default:null"`
	Points      []*Point `json:"points,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Tags        []*Tag   `json:"tags,omitempty" gorm:"many2many:devices_tags;constraint:OnDelete:CASCADE"`
}

type Point struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	DeviceUUID string `json:"device_uuid,omitempty" gorm:"TYPE:string REFERENCES devices;not null;default:null"`
	Tags       []*Tag `json:"tags,omitempty" gorm:"many2many:points_tags;constraint:OnDelete:CASCADE"`
}

type Tag struct {
	Tag          string         `json:"tag" gorm:"type:varchar(255);unique;not null;default:null;primaryKey"`
	Networks     []*Network     `json:"networks,omitempty" gorm:"many2many:networks_tags;constraint:OnDelete:CASCADE"`
	Devices      []*Device      `json:"devices,omitempty" gorm:"many2many:devices_tags;constraint:OnDelete:CASCADE"`
	Points       []*Point       `json:"points,omitempty" gorm:"many2many:points_tags;constraint:OnDelete:CASCADE"`
	StreamClones []*StreamClone `json:"stream_clones,omitempty" gorm:"many2many:stream_clones_tags;constraint:OnDelete:CASCADE"`
	Consumers    []*Consumer    `json:"consumers,omitempty" gorm:"many2many:consumers_tags;constraint:OnDelete:CASCADE"`
}

type History struct {
	ID        int       `json:"id" gorm:"primary_key"`
	UUID      string    `json:"uuid" gorm:"primary_key"`
	Value     float64   `json:"value" gorm:"primary_key"`
	Timestamp time.Time `json:"timestamp" gorm:"primary_key"`
}

type HistoryData struct {
	Value            float64   `json:"value"`
	Timestamp        time.Time `json:"timestamp"`
	RubixNetworkUUID string    `json:"rubix_network_uuid"`
	RubixNetworkName string    `json:"rubix_network_name"`
	RubixDeviceUUID  string    `json:"rubix_device_uuid"`
	RubixDeviceName  string    `json:"rubix_device_name"`
	RubixPointUUID   string    `json:"rubix_point_uuid"`
	RubixPointName   string    `json:"rubix_point_name"`
	TagData
	FlowNetworkCloneData
}

type TagData struct {
	NetworkTag string `json:"network_tag,omitempty"`
	DeviceTag  string `json:"device_tag,omitempty"`
	PointTag   string `json:"point_tag,omitempty"`
}

type FlowNetworkCloneData struct {
	GlobalUUID string `json:"global_uuid,omitempty"`
	ClintId    string `json:"client_id,omitempty"`
	ClintName  string `json:"client_name,omitempty"`
	SiteId     string `json:"site_id,omitempty"`
	SiteName   string `json:"site_name,omitempty"`
	DeviceId   string `json:"device_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
}
