package model

import "gorm.io/datatypes"

//Schedule model
type Schedule struct {
	CommonUUID
	CommonName
	CommonEnable
	DataStore datatypes.JSON `json:"data_store"`
}

//
////ScheduleBody model
//type ScheduleBody struct {
//	CommonUUID
//	CommonName
//	CommonEnable
//	DataStore datatypes.JSON `json:"data_store"`
//
//}
//
