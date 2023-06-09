package schedule

import (
	"github.com/NubeIO/rubix-os/utils/boolean"
	log "github.com/sirupsen/logrus"
	"time"
)

// CheckEventScheduleEntry checks if there is a EventScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CheckEventScheduleEntry(entry EventScheduleEntry, timezone string) (ScheduleCheckerResult, error) {
	result := ScheduleCheckerResult{}
	result.Payload = entry.Value
	result.IsActive = false
	result.IsException = false

	var err error

	now := time.Now().UTC()
	result.CheckTime = now.Unix()

	// make a ScheduleCheckerResult to be combined later
	toCombine := ScheduleCheckerResult{}
	toCombine.CheckTime = now.Unix()

	for _, StartStopPair := range entry.Dates {
		toCombine.IsActive = false
		toCombine.PeriodStart = 0
		toCombine.PeriodStop = 0
		toCombine.NextStart = 0
		toCombine.NextStop = 0
		toCombine.ErrorFlag = false
		toCombine.AlertFlag = false
		toCombine.ErrorStrings = []string{}

		// parse timezone, start and stop timestamps
		/*
			location, err := time.LoadLocation(timezone)
			if err != nil {
				result.ErrorFlag = true
				result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Timezone")
				return result, err
			}
		*/
		timeParseFormatCONSTANT := "2006-01-02T15:04" // This is the format string for the Schedule JSON see: https://pkg.go.dev/time#pkg-constants
		// entryStart, err1 := time.ParseInLocation(timeParseFormatCONSTANT, StartStopPair.Start, location)
		// entryStop, err2 := time.ParseInLocation(timeParseFormatCONSTANT, StartStopPair.End, location)
		entryStart, err1 := time.Parse(timeParseFormatCONSTANT, StartStopPair.Start)
		entryStop, err2 := time.Parse(timeParseFormatCONSTANT, StartStopPair.End)
		if err1 != nil || err2 != nil {
			toCombine.ErrorFlag = true
			toCombine.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Start/End Timestamps")
			continue
		}
		// log.Println("pair #: ", key, "    entryStart: ", entryStart.In(location), "entryStop: ", entryStop.In(location))

		// Check if the schedule is active now.  Calculate PeriodStart, PeriodStop, NextStart, and NextStop where applicable
		if now.Before(entryStop) { // schedule has not yet started, or is active.
			if now.After(entryStart) {
				toCombine.IsActive = true
			}
			// PeriodStartStop
			toCombine.PeriodStart = entryStart.Unix()
			toCombine.PeriodStop = entryStop.Unix()

			// combine current StartStopPair with results
			result, err = CombineScheduleCheckerResults(result, toCombine, timezone)
			if err != nil {
				log.Errorf("CheckEventScheduleEntry %v\n", err)
			}
		}
	}
	AddHumanReadableDatetimes(&result, timezone)
	return result, nil
}

// CheckEventScheduleCollection checks if there is a EventScheduleEntry in the provided EventScheduleCollection that matches the specified schedule Name and is currently within the scheduled period.
func CheckEventScheduleCollection(scheduleMap TypeEvents, scheduleName, timezone string) ScheduleCheckerResult {
	finalResult := ScheduleCheckerResult{}
	for _, scheduleEntry := range scheduleMap {
		// fmt.Println(fmt.Sprintf("CheckEventScheduleCollection() scheduleEntry: %+v", scheduleEntry))
		if scheduleEntry.Enable == nil || boolean.IsTrue(scheduleEntry.Enable) {
			if scheduleName == "ANY" || scheduleName == "ALL" || scheduleEntry.Name == scheduleName {
				// fmt.Println(fmt.Sprintf("EVENT SCHEDULE: %v %+v", key, scheduleEntry))
				singleResult, err := CheckEventScheduleEntry(scheduleEntry, timezone)
				singleResult.Name = scheduleName

				// fmt.Println(fmt.Sprintf("finalResult: %+v  singleResult: %+v", finalResult, singleResult))
				if err != nil {
					log.Errorf("CheckEventScheduleEntry %v\n", err)
				}

				finalResult, err = CombineScheduleCheckerResults(finalResult, singleResult, timezone)
				// fmt.Println(fmt.Sprintf("finalResult: %+v", finalResult))
				if err != nil {
					log.Errorf("CheckEventScheduleEntry %v\n", err)
				}
			}
		}
	}
	AddHumanReadableDatetimes(&finalResult, timezone)
	return finalResult
}

// EventCheck checks all Event Schedules in the payload for active periods. It returns a combined ScheduleCheckerResult of all Event Schedules.
func EventCheck(events TypeEvents, scheduleName, timezone string) (ScheduleCheckerResult, error) {
	results := CheckEventScheduleCollection(events, scheduleName, timezone)
	return results, nil
}
