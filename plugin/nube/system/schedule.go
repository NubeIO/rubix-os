package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/flow-framework/src/schedule"
	"github.com/NubeIO/flow-framework/src/utilstime"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (i *Instance) run() {
	getSch, err := i.db.GetOneScheduleByArgs(api.Args{Name: utils.NewStringAddress("HVAC")})
	if err != nil || getSch == nil {
		log.Errorf("system-plugin-schedule: issue on GetLatestProducerHistoryByProducerName %v\n", err)
		return
	}
	ScheduleJSON, err := schedule.DecodeSchedule(getSch.Schedule)
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
		return
	}
	scheduleNameToCheck := "ALL" //TODO: we need a way to specify the schedule name that is being checked for.

	timezone := ScheduleJSON.Config.TimeZone
	_, err = time.LoadLocation(timezone)
	if err != nil || timezone == "" { // If timezone field is not assigned or invalid, get timezone from System Time
		log.Error("system-plugin-schedule: CheckWeeklyScheduleCollection(): invalid schedule timezone. checking with system time.")
		systemTimezone := strings.Split((*utilstime.SystemTime()).HardwareClock.Timezone, " ")[0]
		//fmt.Println("systemTimezone 2: ", systemTimezone)
		if systemTimezone == "" {
			zone, _ := utilstime.GetHardwareTZ()
			timezone = zone
		} else {
			timezone = systemTimezone
		}
	}

	//CHECK WEEKLY SCHEDULES
	weeklyResult, err := schedule.WeeklyCheck(ScheduleJSON.Schedules.Weekly, scheduleNameToCheck, timezone) //This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
	}
	// CHECK EVENT SCHEDULES
	eventResult, err := schedule.EventCheck(ScheduleJSON.Schedules.Events, scheduleNameToCheck, timezone) //This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on EventCheck %v\n", err)
	}
	//Combine Event and Weekly schedule results.
	weeklyAndEventResult, err := schedule.CombineScheduleCheckerResults(weeklyResult, eventResult)
	// CHECK EXCEPTION SCHEDULES
	exceptionResult, err := schedule.ExceptionCheck(ScheduleJSON.Schedules.Exceptions, scheduleNameToCheck, timezone) //This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on ExceptionCheck %v\n", err)
	}
	if exceptionResult.CheckIfEmpty() {
		log.Println("system-plugin-schedule: Exception schedule is empty")
	}

	finalResult, err := schedule.ApplyExceptionSchedule(weeklyAndEventResult, exceptionResult) //This applies the exception schedule to mask the combined weekly and event schedules.
	if err != nil {
		log.Errorf("system-plugin-schedule: issue on ApplyExceptionSchedule %v\n", err)
	}

	log.Printf("system-plugin-schedule: finalResult: %+v\n", finalResult.IsActive)

	if getSch != nil {
		i.store.Set(getSch.Name, finalResult, -1)
		s := new(model.Schedule)
		if finalResult.IsActive {
			s.IsActive = utils.NewTrue()
		} else {
			s.IsActive = utils.NewFalse()
		}
		if utils.IsTrue(s.IsActive) != utils.IsTrue(getSch.IsActive) {
			log.Printf("system-plugin-schedule: UPDATE SCHEDULE IN DB %v\n", getSch.Name)
			_, err = i.db.UpdateSchedule(getSch.UUID, s)
			if err != nil {
				log.Errorf("system-plugin-schedule: issue on UpdateSchedule %v\n", getSch.UUID)
			}
		}
	}
}

//
//func (i *Instance) run() {
//	class, err := i.db.GetWriters(api.Args{WriterThingClass: utils.NewStringAddress("schedule")})
//	if err != nil {
//		log.Errorf("system-plugin-schedule: issue on GetWriters %v\n", err)
//	}
//
//	//class, err := i.db.GetLatestProducerHistoryByProducerName("aa")
//	//if err != nil {
//	//	log.Errorf("system-plugin-schedule: issue on GetWriters %v\n", err)
//	//}
//	for _, v := range class {
//		decodeSchedule, err := schedule.DecodeSchedule(v.DataStore)
//		if err != nil {
//			log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
//		}
//		schUUID := v.WriterThingUUID
//		for _, v := range decodeSchedule.Weekly {
//
//			scheduleNameToCheck := v.Name //TODO: we need a way to specify the schedule name that is being checked for.

//			getSch, err := i.db.GetSchedule(schUUID)
//			if err != nil {
//				log.Errorf("system-plugin-schedule: issue on GetSchedule %v\n", err)
//			}
//			weeklyResult, err := schedule.WeeklyCheck(decodeSchedule.Weekly, scheduleNameToCheck) //This will check for any active schedules with defined name.
//			if err != nil {
//				log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
//			}
//			// CHECK EVENT SCHEDULES
//			eventResult, err := schedule.EventCheck(decodeSchedule.Events, scheduleNameToCheck) //This will check for any active schedules with defined name.
//			if err != nil {
//				log.Errorf("system-plugin-schedule: issue on EventCheck %v\n", err)
//			}
//			//Combine Event and Weekly schedule results.
//			weeklyAndEventResult, err := schedule.CombineScheduleCheckerResults(weeklyResult, eventResult)
//			// CHECK EXCEPTION SCHEDULES
//			exceptionResult, err := schedule.ExceptionCheck(decodeSchedule.Exceptions, scheduleNameToCheck) //This will check for any active schedules with defined name.
//			if err != nil {
//				log.Errorf("system-plugin-schedule: issue on ExceptionCheck %v\n", err)
//			}
//			if exceptionResult.CheckIfEmpty() {
//				log.Println("Exception schedule is empty")
//			}
//
//			finalResult, err := schedule.ApplyExceptionSchedule(weeklyAndEventResult, exceptionResult) //This applies the exception schedule to mask the combined weekly and event schedules.
//			if err != nil {
//				log.Errorf("system-plugin-schedule: issue on ApplyExceptionSchedule %v\n", err)
//			}
//			log.Println("finalResult")
//			log.Printf("%+v\n", finalResult.IsActive)
//			i.store.Set(finalResult.Name, finalResult, -1)
//			s := new(model.Schedule)
//			if finalResult.IsActive {
//				s.IsActive = utils.NewTrue()
//			} else {
//				s.IsActive = utils.NewFalse()
//			}
//			if getSch != nil {
//				fmt.Println(utils.IsTrue(s.IsActive), utils.IsTrue(getSch.IsActive))
//				if utils.IsTrue(s.IsActive) != utils.IsTrue(getSch.IsActive) {
//					log.Printf("system-plugin-schedule: UPDATE SCHEDULE IN DB %v\n", getSch.Name)
//					_, err = i.db.UpdateSchedule(getSch.UUID, s)
//					if err != nil {
//						log.Errorf("system-plugin-schedule: issue on UpdateSchedule %v\n", getSch.UUID)
//					}
//				}
//
//			}
//		}
//	}
//}
//
