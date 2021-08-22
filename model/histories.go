package model


//History point history settings TODO add in later
type History struct {
	Type string //cov, interval, cov_interval
	Duration int //15min
	SizeLimit int //max amount of records to keep, the newest will override the oldest record

}


//PointStore for storing the history
type PointStore struct {
	PointUuid     			string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	CommonCreated
	Value          			string `json:"value" sql:"DEFAULT:NULL"`
	ValueOriginal 			string `json:"value_original"`
	ValueRaw      			string `json:"value_raw"`
	Fault        			string `json:"fault"`
	FaultMessage 			string `json:"fault_message"`
	TsValue 				string `json:"ts_value"`
	TsFault 				string `json:"ts_fault"`

}

