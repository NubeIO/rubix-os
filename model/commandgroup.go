package model

// CommandGroup TODO add in later
//CommandGroup is for issuing global schedule writes or global point writes (as in send a value to any point associated with this group)
type CommandGroup struct {
	CommonUUID
	CommonName
	CommonEnable
	CommonDescription
	CommandUse         string `json:"command_use,omitempty"` //common uses will be point write to many points, master schedules or schedule grouping
	StreamUUID         string `json:"stream_uuid,omitempty" gorm:"TYPE:string REFERENCES streams;null;default:null"`
	WriteValue         string `json:"write_value,omitempty"`
	WritePriority      string `json:"write_priority,omitempty"`       //TODO maybe remove this and just use the writeJSON as we need things like Schedules aswell
	WritePriorityArray string `json:"write_priority_array,omitempty"` //TODO maybe remove this and just use the writeJSON as we need things like Schedules aswell
	WriteJSON          string `json:"write_json,omitempty"`           //TODO add data model in later
	StartDate          string `json:"start_date,omitempty"`           // START at 25:11:2021:13:00
	EndDate            string `json:"end_date,omitempty"`             // START at 25:11:2021:13:30
	Value              string `json:"value,omitempty"`
	Priority           string `json:"priority,omitempty"`
	CommonCreated
}

type CommandSlaves struct {
	CommonUUID
	CommonName
	CommonEnable
	CommonDescription
}
