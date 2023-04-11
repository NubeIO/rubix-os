package main

import (
	"fmt"
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
		inst.systemErrorMsg(fmt.Sprintf("Schedule Checks: GetSchedules %s", err.Error()))
		return
	} else {
		inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: run schedule checks, schedule count: %d", len(schedules)))
	}
	for _, sch := range schedules {
		ScheduleJSON, err := schedule.DecodeSchedule(sch.Schedule)
		if err != nil {
			inst.systemErrorMsg(fmt.Sprintf("Schedule Checks: issue on DecodeSchedule %v\n", err))
			return
		}
		if !boolean.IsTrue(sch.Enable) {
			inst.systemDebugMsg("Schedule Checks: runSchedule() sch is not enabled so skip logic. name:", sch.Name)
			continue
		}
		scheduleNameToCheck := "ALL" // TODO: we may need a way to specify the schedule name that is being checked for.

		var timezone = ScheduleJSON.Config.TimeZone
		if timezone == "" {
			timezone = sch.TimeZone
		}
		_, err = time.LoadLocation(timezone)
		if timezone == "" || err != nil {
			log.Error("Schedule Checks: CheckWeeklyScheduleCollection(): no timezone pass in from user")
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
			inst.systemErrorMsg(fmt.Sprintf("Schedule Checks: issue on WeeklyCheck %v\n", err))
		} else {
			inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: weekly schedule: %s  is-active: %t", weeklyResult.Name, weeklyResult.IsActive))
		}
		// CHECK EVENT SCHEDULES
		eventResult, err := schedule.EventCheck(ScheduleJSON.Schedules.Events, scheduleNameToCheck, timezone) // This will check for any active schedules with defined name.
		if err != nil {
			inst.systemErrorMsg(fmt.Sprintf("Schedule Checks: issue on eventResult %s", err.Error()))
		} else {
			inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: event schedule: %s  is-active: %t", eventResult.Name, eventResult.IsActive))
		}
		inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: eventResult: %+v", eventResult))

		// Combine Event and Weekly schedule results.
		weeklyAndEventResult, err := schedule.CombineScheduleCheckerResults(weeklyResult, eventResult, timezone)
		if err != nil {
			inst.systemErrorMsg(fmt.Sprintf("Schedule Checks: issue on weeklyAndEventResult %s", err.Error()))
		} else {
			inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: weekly & event schedule: %s  is-active: %t", weeklyAndEventResult.Name, weeklyAndEventResult.IsActive))
		}
		inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: weeklyAndEventResult: %+v", weeklyAndEventResult))

		// CHECK EXCEPTION SCHEDULES
		// inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: exception schedule: %+v", ScheduleJSON.Schedules.Exceptions))
		exceptionResult, err := schedule.ExceptionCheck(ScheduleJSON.Schedules.Exceptions, scheduleNameToCheck, timezone) // This will check for any active schedules with defined name.
		if err != nil {
			inst.systemErrorMsg(fmt.Sprintf("Schedule Checks: issue on exceptionResult %s", err.Error()))
		} else {
			inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: exception schedule: %s  is-active: %t", exceptionResult.Name, exceptionResult.IsActive))
		}
		if exceptionResult.CheckIfEmpty() {
			inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: exception schedule is empty: %s", exceptionResult.Name))
		}
		inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: exceptionResult: %+v", exceptionResult))

		finalResult, err := schedule.ApplyExceptionSchedule(weeklyAndEventResult, exceptionResult, timezone) // This applies the exception schedule to mask the combined weekly and event schedules.
		if err != nil {
			inst.systemErrorMsg(fmt.Sprintf("Schedule Checks: final-result: %s", err.Error()))
		}
		inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: final-result: %s  is-active: %t timezone: %s", finalResult.Name, finalResult.IsActive, timezone))
		inst.systemDebugMsg(fmt.Sprintf("Schedule Checks: finalResult: %+v", finalResult))
		if sch != nil {
			inst.store.Set(sch.Name, finalResult, -1)
			sch.IsActive = boolean.New(finalResult.IsActive)
			sch.ActiveWeekly = boolean.New(weeklyResult.IsActive)
			sch.ActiveException = boolean.New(exceptionResult.IsActive)
			sch.ActiveEvent = boolean.New(eventResult.IsActive)
			sch.Payload = finalResult.Payload

			sch.PeriodStart = finalResult.PeriodStart
			if finalResult.PeriodStart == 0 {
				sch.PeriodStartString = ""
			} else {
				sch.PeriodStartString = finalResult.PeriodStartString
			}

			sch.PeriodStop = finalResult.PeriodStop
			if finalResult.PeriodStop == 0 {
				sch.PeriodStopString = ""
			} else {
				sch.PeriodStopString = finalResult.PeriodStopString
			}

			sch.NextStart = finalResult.NextStart
			if finalResult.NextStart == 0 {
				sch.NextStartString = ""
			} else {
				sch.NextStartString = finalResult.NextStartString
			}

			sch.NextStop = finalResult.NextStop
			if finalResult.NextStop == 0 {
				sch.NextStopString = ""
			} else {
				sch.NextStopString = finalResult.NextStopString
			}

			_, err = inst.db.UpdateScheduleAllProps(sch.UUID, sch)
			if err != nil {
				inst.systemErrorMsg(fmt.Sprintf("Schedule Checks: issue on UpdateSchedule %s, error: %v", sch.UUID, err))
			}
		}
	}
}
