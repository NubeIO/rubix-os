package model

// CommandGroup TODO add in later
//CommandGroup is for issuing global schedule writes or global point writes (as in send a value to any point associated with this group)
type CommandGroup struct {
	CommonUUID
	CommonName
	CommonEnable
	CommonDescription
	CommandUse         string `json:"command_use"` //common uses will be point write to many points, master schedules or schedule grouping
	StreamUUID         string `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;null;default:null"`
	WriteValue         string `json:"write_value"`
	WritePriority      string `json:"write_priority"`       //TODO maybe remove this and just use the writeJSON as we need things like Schedules aswell
	WritePriorityArray string `json:"write_priority_array"` //TODO maybe remove this and just use the writeJSON as we need things like Schedules aswell
	WriteJSON          string `json:"write_json"`           //TODO add data model in later
	StartDate          string `json:"start_date"`           // START at 25:11:2021:13:00
	EndDate            string `json:"end_date"`             // START at 25:11:2021:13:30
	Value              string `json:"value"`
	Priority           string `json:"priority"`
	CommonCreated
}

type CommandSlaves struct {
	CommonUUID
	CommonName
	CommonEnable
	CommonDescription
}
