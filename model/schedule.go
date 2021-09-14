package model

import "gorm.io/datatypes"

//Schedule model
type Schedule struct {
	CommonUUID
	CommonName
	CommonEnable
	CommonThingClass
	CommonThingType
	CommonCreated
	DataStore datatypes.JSON `json:"data_store"`
}
