package model

import (
	"errors"
	"log"
	"time"
)

type Job struct {
	//sync.Mutex
	ID          *int64    `json:"id" sql:"id"`
	Frequency   string    `json:"frequency,omitempty" sql:"frequency"`
	StartDate   time.Time `json:"start_date,omitempty" sql:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty" sql:"end_date"`
	CronEntryID int       `json:"-" sql:"cron_entry_id"`
	IsActive    bool      `json:"is_active" sql:"is_active"`
}



var RemoveJob = make(chan int)

//func (j Job) Run() {
//	if time.Now().Unix() > j.StartDate.Unix()&& ((j.EndDate.Unix()>0 && j.EndDate.Unix()>time.Now().Unix())||j.EndDate.Unix()<0){
//
//		//log.Println("Hi from job: ", *j.ID)
//
//	} else if j.EndDate.UnixNano()>0 && time.Now().UnixNano() > j.EndDate.UnixNano() {
//		//inform entry id of job to remove job from cron when it is expired
//		RemoveJob <- j.CronEntryID
//		log.Printf("job %v is expired ", *j.ID)
//	}
//}


func (j *Job) FormatJobData() (err error) {
	if j.Frequency == "" {
		log.Println("invalid frequency: ", j.Frequency)
		err = errors.New("invalid frequency")
		return
	}
	return
}