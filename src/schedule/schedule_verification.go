package schedule

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
)

func ScheduleTest() {
	json, err := ioutil.ReadFile("./schTestNewJSON1.json")
	if err != nil {
		log.Errorf("ReadFile %v\n", err)
	}
	ScheduleJSON, err := DecodeSchedule(json)
	log.Println("decodeSchedule: ", ScheduleJSON)

	// scheduleNameToCheck :=  //TODO: we need a way to specify the schedule name that is being checked for.
	scheduleNameToCheck := "HVAC"

	timezone := ScheduleJSON.Config.TimeZone
	_, err = time.LoadLocation(timezone)
	if err != nil || timezone == "" { // If timezone field is not assigned or invalid, get timezone from System Time
		log.Error("CheckWeeklyScheduleCollection(): invalid schedule timezone. checking with system time.")
		systemTimezone := strings.Split((*utilstime.SystemTime()).HardwareClock.Timezone, " ")[0]
		// fmt.Println("systemTimezone 2: ", systemTimezone)
		if systemTimezone == "" {
			zone, _ := utilstime.GetHardwareTZ()
			timezone = zone
		} else {
			timezone = systemTimezone
		}
	}

	weeklyResult, err := WeeklyCheck(ScheduleJSON.Schedules.Weekly, scheduleNameToCheck, timezone) // This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("Schedule Checks: issue on WeeklyCheck %v\n", err)
	}
	// fmt.Println("weeklyResult")
	// fmt.Printf("%+v\n", weeklyResult)

	// CHECK EVENT SCHEDULES
	// eventResult, err := schedule.EventCheck(decodeSchedule.Events, "ANY")  //This will check for any active schedules with any name
	eventResult, err := EventCheck(ScheduleJSON.Schedules.Events, scheduleNameToCheck, timezone) // This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("Schedule Checks: issue on EventCheck %v\n", err)
	}
	// fmt.Println("eventResult")
	// fmt.Printf("%+v\n", eventResult)

	// Combine Event and Weekly schedule results.
	weeklyAndEventResult, err := CombineScheduleCheckerResults(weeklyResult, eventResult, timezone)
	// fmt.Println("weeklyAndEventResult")
	// fmt.Printf("%+v\n", weeklyAndEventResult)

	// CHECK EXCEPTION SCHEDULES
	// exceptionResult, err := schedule.ExceptionCheck(decodeSchedule.Exceptions, "ANY")  //This will check for any active schedules with any name
	exceptionResult, err := ExceptionCheck(ScheduleJSON.Schedules.Exceptions, scheduleNameToCheck, timezone) // This will check for any active schedules with defined name.
	if err != nil {
		log.Errorf("Schedule Checks: issue on ExceptionCheck %v\n", err)
	}
	// fmt.Println("exceptionResult")
	// fmt.Printf("%+v\n", exceptionResult)
	if exceptionResult.CheckIfEmpty() {
		// fmt.Println("Exception schedule is empty")
	}

	finalResult, err := ApplyExceptionSchedule(weeklyAndEventResult, exceptionResult, timezone) // This applies the exception schedule to mask the combined weekly and event schedules.
	if err != nil {
		log.Errorf("Schedule Checks: issue on ApplyExceptionSchedule %v\n", err)
	}
	fmt.Println("finalResult")
	fmt.Printf("%+v\n", finalResult)

	// fmt.Println("schedule is ", finalResult.IsActive)
}
