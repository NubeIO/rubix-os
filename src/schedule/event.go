package schedule

import (
	"time"
)

//CheckEventScheduleEntry checks if there is a EventScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CheckEventScheduleEntry(entry EventScheduleEntry) (ScheduleCheckerResult, error) {
	result := ScheduleCheckerResult{}
	result.Payload = entry.Value
	result.IsActive = false

	now := time.Now().UTC()
	result.CheckTime = now.Unix()

	//make a ScheduleCheckerResult to be combined later
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

		//parse start and stop timestamps
		entryStart, err1 := time.Parse(time.RFC3339Nano, StartStopPair.Start)
		entryStop, err2 := time.Parse(time.RFC3339Nano, StartStopPair.End)
		if err1 != nil || err2 != nil {
			toCombine.ErrorFlag = true
			toCombine.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Start/End Timestamps")
			continue
		}
		//log.Println("pair #: ", i, "    entryStart: ", entryStart, "entryStop: ", entryStop)

		//Check if the schedule is active now.  Calculate PeriodStart, PeriodStop, NextStart, and NextStop where applicable
		if now.Before(entryStop) { //schedule has not yet started, or is active.
			if now.After(entryStart) {
				toCombine.IsActive = true
			}
			//PeriodStartStop
			toCombine.PeriodStart = entryStart.Unix()
			toCombine.PeriodStop = entryStop.Unix()

			//combine current StartStopPair with results
			result = CombineScheduleCheckerResults(result, toCombine, true)
		}
	}
	return result, nil
}

//CheckEventScheduleCollection checks if there is a WeeklyScheduleEntry in the provided EventScheduleCollection that matches the specified schedule Name and is currently within the scheduled period.
func CheckEventScheduleCollection(scheduleMap TypeEvents, scheduleName string) ScheduleCheckerResult {
	finalResult := ScheduleCheckerResult{}
	for _, scheduleEntry := range scheduleMap {
		if scheduleName == "ANY" || scheduleEntry.Name == scheduleName {
			//fmt.Println("EVENT SCHEDULE ", i, ": ", scheduleEntry)
			singleResult, err := CheckEventScheduleEntry(scheduleEntry)
			//fmt.Println("finalResult ", finalResult, "singleResult: ", singleResult)

			finalResult = CombineScheduleCheckerResults(finalResult, singleResult, true)
			//fmt.Println("finalResult ", finalResult)
		}
	}
	return finalResult
}

func EventCheck(events TypeEvents, scheduleName string) (ScheduleCheckerResult, error) {
	results := CheckEventScheduleCollection(events, scheduleName)
	return results, nil
}
