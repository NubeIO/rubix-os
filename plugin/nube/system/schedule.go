package main

import (
	"github.com/NubeIO/flow-framework/src/schedule"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (inst *Instance) runSchedule() {
	schedules, err := inst.db.GetSchedules()
	if err != nil {
		log.Errorf("system-plugin-schedule: GetSchedules %s", err.Error())
		return
	} else {
		log.Infof("system-plugin-schedule: run schedule executaion, schedule count: %d", len(schedules))
	}
	for _, sch := range schedules {
		ScheduleJSON, err := schedule.DecodeSchedule(sch.Schedule)
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
			return
		}
		if !boolean.IsTrue(sch.Enable) {
			log.Infoln("system-plugin-schedule: runSchedule() sch is not enabled so skip logic name:", sch.Name)
			continue
		}
		scheduleNameToCheck := "ALL" // TODO: we need a way to specify the schedule name that is being checked for.

		var timezone = ScheduleJSON.Config.TimeZone
		if timezone == "" {
			timezone = sch.TimeZone
		}
		_, err = time.LoadLocation(timezone)
		if timezone == "" || err != nil {
			log.Error("system-plugin-schedule: CheckWeeklyScheduleCollection(): no timezone pass in from user")
			systemTimezone := strings.Split((*utilstime.SystemTime()).HardwareClock.Timezone, " ")[0]
			if systemTimezone == "" {
				zone, _ := utilstime.GetHardwareTZ()
				timezone = zone
			} else {
				timezone = systemTimezone
			}
			sch.TimeZone = timezone
		}
		// CHECK WEEKLY SCHEDULES
		weeklyResult, err := schedule.WeeklyCheck(ScheduleJSON.Schedules.Weekly, scheduleNameToCheck, timezone) // This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
		} else {
			log.Infof("system-plugin-schedule: weekly schedule: %s  is-active: %t", weeklyResult.Name, weeklyResult.IsActive)
		}
		// CHECK EVENT SCHEDULES
		eventResult, err := schedule.EventCheck(ScheduleJSON.Schedules.Events, scheduleNameToCheck, timezone) // This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on eventResult %s", err.Error())
		} else {
			log.Infof("system-plugin-schedule: event schedule: %s  is-active: %t", eventResult.Name, eventResult.IsActive)
		}
		// Combine Event and Weekly schedule results.
		weeklyAndEventResult, err := schedule.CombineScheduleCheckerResults(weeklyResult, eventResult)
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on weeklyAndEventResult %s", err.Error())
		} else {
			log.Infof("system-plugin-schedule: weekly & event schedule: %s  is-active: %t", weeklyAndEventResult.Name, weeklyAndEventResult.IsActive)
		}
		//
		// CHECK EXCEPTION SCHEDULES
		exceptionResult, err := schedule.ExceptionCheck(ScheduleJSON.Schedules.Exceptions, scheduleNameToCheck, timezone) // This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on exceptionResult %s", err.Error())
		} else {
			log.Infof("system-plugin-schedule: exception schedule: %s  is-active: %t", weeklyAndEventResult.Name, weeklyAndEventResult.IsActive)
		}
		if exceptionResult.CheckIfEmpty() {
			log.Infof("system-plugin-schedule: exception schedule is empty: %s", exceptionResult.Name)
		}

		finalResult, err := schedule.ApplyExceptionSchedule(weeklyAndEventResult, exceptionResult) // This applies the exception schedule to mask the combined weekly and event schedules.
		if err != nil {
			log.Errorf("system-plugin-schedule: final-result: %s", err.Error())
		}
		log.Infof("system-plugin-schedule: final-result: %s  is-active: %t timezone: %s", weeklyAndEventResult.Name, weeklyAndEventResult.IsActive, timezone)
		if sch != nil {
			inst.store.Set(sch.Name, finalResult, -1)
			sch.IsActive = boolean.New(finalResult.IsActive)
			sch.ActiveWeekly = boolean.New(weeklyResult.IsActive)
			sch.ActiveException = boolean.New(exceptionResult.IsActive)
			sch.ActiveEvent = boolean.New(eventResult.IsActive)
			_, err = inst.db.UpdateSchedule(sch.UUID, sch)
			if err != nil {
				log.Errorf("system-plugin-schedule: issue on UpdateSchedule %s", sch.UUID)
			}
		}
	}
}
