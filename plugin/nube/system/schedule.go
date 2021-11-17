package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/schedule"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

func (i *Instance) run() {
	class, err := i.db.GetWriters(api.Args{WriterThingClass: utils.NewStringAddress("schedule")})
	if err != nil {
		return
	}
	for _, v := range class {
		decodeSchedule, err := schedule.DecodeSchedule(v.DataStore)
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
		}
		for schKey, week := range decodeSchedule.Weekly {
			sch, err := schedule.WeeklyCheck(decodeSchedule.Weekly, week.Name)
			if err != nil {
				log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
			}
			if sch.IsActive {
				i.store.Set(week.Name, sch, -1)
			}
			log.Infof("system-plugin-schedule: schedule schKey %v\n", schKey)
			log.Infof("system-plugin-schedule: schedule Name %v\n", week.Name)
			log.Infof("system-plugin-schedule: schedule NextStart %v\n", time.Unix(sch.NextStart, 0))
			log.Infof("system-plugin-schedule: schedule NextStop %v\n", time.Unix(sch.NextStop, 0))
			log.Infof("system-plugin-schedule: schedule is IsActive %v\n", sch.IsActive)
			log.Infof("system-plugin-schedule: schedule Payload %v\n", sch.Payload)
		}
	}

}

//func (i *Instance) runScheduleAPI() {
//	schedules, err := i.db.GetSchedules()
//	if err != nil {
//		log.Infof("system-plugin-schedule: db get scheudles %v\n", err)
//	}
//	for _, s := range schedules {
//		decodeSchedule, err := schedule.DecodeSchedule(s.Schedules)
//		if err != nil {
//			log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
//		}
//		for schKey, week := range decodeSchedule.Weekly {
//			sch, err := schedule.WeeklyCheck(decodeSchedule.Weekly, week.Name)
//			if err != nil {
//				log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
//			}
//
//			if sch.IsActive {
//				i.store.Set(week.Name, sch, -1)
//			}
//			log.Infof("system-plugin-schedule: schedule schKey %v\n", schKey)
//			log.Infof("system-plugin-schedule: schedule Name %v\n", week.Name)
//			log.Infof("system-plugin-schedule: schedule NextStart %v\n", time.Unix(sch.NextStart, 0))
//			log.Infof("system-plugin-schedule: schedule NextStop %v\n", time.Unix(sch.NextStop, 0))
//			log.Infof("system-plugin-schedule: schedule is IsActive %v\n", sch.IsActive)
//			log.Infof("system-plugin-schedule: schedule Payload %v\n", sch.Payload)
//		}
//	}
//}
