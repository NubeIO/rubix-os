package main

import (
	"github.com/NubeDev/flow-framework/src/jobs"
	"github.com/NubeDev/flow-framework/src/schedule"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) run() {
	schedules, err := i.db.GetSchedules()
	if err != nil {
		log.Infof("system-plugin-schedule: db get scheudles %v\n", err)
	}
	for _, sch := range schedules {
		check, err := schedule.WeeklyCheck(sch.Schedules, sch.Name)
		if err != nil {
			return
		}
		log.Infof("system-plugin-schedule: schedule is active %v\n", check.IsActive)
		log.Infof("system-plugin-schedule: schedule payload %v\n", check.Payload)
	}
}

func (i *Instance) schedule() {
	j, ok := jobs.GetJobService()
	if ok {
		_, err := j.Every(30).Second().Tag("foo").Do(i.run)
		if err != nil {
			log.Infof("system-plugin-schedule: error on create job %v\n", err)
		}
	}
}
