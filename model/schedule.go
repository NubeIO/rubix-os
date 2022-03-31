package model

import (
	"gorm.io/datatypes"
)

type Schedule struct {
	CommonUUID
	CommonNameUnique
	CommonEnable
	CommonThingClass
	CommonThingType
	IsActive *bool          `json:"is_active"`
	IsGlobal *bool          `json:"is_global"`
	Schedule datatypes.JSON `json:"schedule"`
	CommonCreated
}

type ScheduleData struct {
	Schedules Schedules      `json:"schedules,omitempty"`
	Config    datatypes.JSON `json:"config,omitempty"`
}

type ScheduleDataWithConfig struct {
	Schedules Schedules      `json:"schedules,omitempty"`
	Config    ScheduleConfig `json:"config,omitempty"`
}

type ScheduleConfig struct {
	ScheduleNames datatypes.JSON `json:"names"`
	TimeZone      string         `json:"timezone"`
}

type WeeklyMap map[string]Weekly
type EventsMap map[string]Events
type ExceptionMap map[string]Exception

type Schedules struct {
	Events    EventsMap    `json:"events"`
	Weekly    WeeklyMap    `json:"weekly"`
	Exception ExceptionMap `json:"exception"`
}

type Events struct {
	Name  string `json:"name"`
	Dates []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"dates"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

type Weekly struct {
	Name  string   `json:"name"`
	Days  []string `json:"days"`
	Start string   `json:"start"`
	End   string   `json:"end"`
	Value float64  `json:"value"`
	Color string   `json:"color"`
}

type Exception struct {
	Name  string `json:"name"`
	Dates []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"dates"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}
