package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/src/jobs"
	"github.com/NubeDev/flow-framework/src/schedule"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) run() {
	schedules, err := i.db.GetSchedules()
	if err != nil {
		log.Infof("system-plugin-schedule: db get scheudles %v\n", err)
	}
	for _, s := range schedules {
		decodeSchedule, err := schedule.DecodeSchedule(s.Schedules)
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
		}
		for schKey, week := range decodeSchedule.Weekly {
			sch, err := schedule.WeeklyCheck(decodeSchedule.Weekly, week.Name)
			if err != nil {
				log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
			}
			fmt.Println(schKey)
			//TODO we need to now update the schedule
			log.Infof("system-plugin-schedule: schedule is active %v\n", week.Name)
			log.Infof("system-plugin-schedule: schedule is active %v\n", sch.NextStart)
			log.Infof("system-plugin-schedule: schedule is active %v\n", sch.IsActive)
			log.Infof("system-plugin-schedule: schedule payload %v\n", sch.Payload)
		}
	}
}

func (i *Instance) schedule() {
	j, ok := jobs.GetJobService()
	if ok {
		_, err := j.Every(30).Second().Do(i.run)
		if err != nil {
			log.Infof("system-plugin-schedule: error on create job %v\n", err)
		}
	}
}
