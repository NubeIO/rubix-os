package model

import (
	"time"
)

//https://github.com/asaskevich/govalidator TODO add in a validator

type Job struct {
	CommonUUID
	CommonName
	CommonDescription
	Frequency   				string 	  `json:"frequency,omitempty" sql:"frequency"`
	StartDate   				time.Time `json:"start_date,omitempty" sql:"start_date"`
	EndDate     				time.Time `json:"end_date,omitempty" sql:"end_date"`
	CommonEnable
	DestroyAfterCompleted   	bool      `json:"destroy_after_completed" sql:"destroy_after_completed"`
	PluginConfId     			string  `json:"plugin_conf_id" gorm:"TYPE:varchar(255) REFERENCES plugin_confs;null;default:null"`
	CommonCreated

}
//
//type JobSubscriber struct {
//	JobUUID    		string  `json:"job_uuid" binding:"required" gorm:"TYPE:varchar(255) REFERENCES jobs;not null;default:null"`
//	Subscriber
//}