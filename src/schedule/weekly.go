package schedule

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

//CheckWeeklyScheduleEntry checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleEntry(entry model.Weekly, timezone string) ScheduleCheckerResult {
	result := ScheduleCheckerResult{}
	result.Payload = entry.Value
	result.IsActive = false
	result.IsException = false

	//get time.Location for entry timezone and check timezone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Timezone")
		return result
	}

	now := time.Now().In(location)
	result.CheckTime = now.Unix()
	//get day of week and compare with entry.DaysNums
	nowYear, nowMonth, nowDate := now.Date()
	nowDayOfWeek := DaysOfTheWeek(now.Weekday())

	//parse start and stop times
	var entryStartHour, entryStartMinute, entryStopHour, entryStopMinutes int
	n, err1 := fmt.Sscanf(entry.Start, "%d:%d", &entryStartHour, &entryStartMinute)
	m, err2 := fmt.Sscanf(entry.End, "%d:%d", &entryStopHour, &entryStopMinutes)
	if n != 2 || m != 2 || err1 != nil || err2 != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Start/Stop Time")
		return result
	}

	//parse start and end time into current day timestamps
	startTimestamp := time.Date(nowYear, nowMonth, nowDate, entryStartHour, entryStartMinute, 0, 0, location)
	stopTimestamp := time.Date(nowYear, nowMonth, nowDate, entryStopHour, entryStopMinutes, 59, 0, location)

	//Check if the schedule is active today
	scheduleActiveToday := false
	dayStringsToIntegers := GetDaysStringsToIntegers(entry)
	for _, day := range dayStringsToIntegers {
		if day == nowDayOfWeek {
			scheduleActiveToday = true
			break
		}
	}

	//find the next active schedule day
	nextDay, err := getNextScheduleDay(nowDayOfWeek, dayStringsToIntegers)
	if err != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Scheduled Days are Invalid")
		return result
	}

	nextDayDuration := getDurationTillNextScheduleDay(nowDayOfWeek, nextDay)

	if scheduleActiveToday && now.Before(stopTimestamp) { //scheduled today and hasn't finished yet
		//check if today's schedule is currently active
		if now.After(startTimestamp) {
			result.IsActive = true
		}
		//PeriodStartStop
		result.PeriodStart = startTimestamp.Unix()
		result.PeriodStop = stopTimestamp.Unix()
		//NextStartStop
		result.NextStart = startTimestamp.Add(nextDayDuration).Unix()
		result.NextStop = stopTimestamp.Add(nextDayDuration).Unix()
	} else { // not scheduled today, OR scheduled today, but period already passed
		//PeriodStartStop
		result.PeriodStart = startTimestamp.Add(nextDayDuration).Unix()
		result.PeriodStop = stopTimestamp.Add(nextDayDuration).Unix()
		//NextStartStop
		//next scheduled day is being used for PeriodStart and PeriodStop, so get the NEXT next scheduled day for NextStart and NextStop
		nextNextDay, err := getNextScheduleDay(nextDay, dayStringsToIntegers)
		if err != nil {
			result.ErrorFlag = true
			result.ErrorStrings = append(result.ErrorStrings, "Critical: Scheduled Days are Invalid")
			return result
		}
		nextNextDayDuration := getDurationTillNextScheduleDay(nextDay, nextNextDay)
		result.NextStart = startTimestamp.Add(nextDayDuration).Add(nextNextDayDuration).Unix()
		result.NextStop = stopTimestamp.Add(nextDayDuration).Add(nextNextDayDuration).Unix()
	}
	fmt.Println("CheckWeeklyScheduleEntry RETURN: ", result)
	AddHumanReadableDatetimes(&result)
	return result
}

//CheckWeeklyScheduleCollection checks if there is a WeeklyScheduleEntry in the provided WeeklyScheduleCollection that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleCollection(weeklySchedulesMap model.WeeklyMap, scheduleName, timezone string) ScheduleCheckerResult {
	finalResult := ScheduleCheckerResult{}
	var singleResult ScheduleCheckerResult
	count := 0
	var err error
	for _, weeklySchedule := range weeklySchedulesMap {
		if scheduleName == "ANY" || scheduleName == "ALL" || weeklySchedule.Name == scheduleName {
			singleResult = CheckWeeklyScheduleEntry(weeklySchedule, timezone)
			singleResult.Name = scheduleName
			if count == 0 {
				finalResult = singleResult
			} else {
				finalResult, err = CombineScheduleCheckerResults(finalResult, singleResult)
				if err != nil {
					log.Errorf("CheckEventScheduleEntry %v\n", err)
				}
			}
			count++
		}
	}
	AddHumanReadableDatetimes(&finalResult)
	return finalResult
}

//getNextScheduleDay() returns the next scheduled day.
func getNextScheduleDay(today DaysOfTheWeek, scheduleDays []DaysOfTheWeek) (DaysOfTheWeek, error) {
	if len(scheduleDays) == 0 {
		return 0, errors.New("NO DAYS SCHEDULED")
	}
	i := int(today)
	j := 0
	var nextDay DaysOfTheWeek = 0
	nextFound := false
	//check until the next scheduled day is found in scheduleDays
	for j < 7 {
		//check the next day (rollover at end of week [Saturday(6)]
		if i == 6 {
			i = 0
		} else {
			i++
		}
		//check each of the days in scheduleDays
		for x := 0; x < len(scheduleDays); x++ {
			if scheduleDays[x] == DaysOfTheWeek(i) {
				nextDay = DaysOfTheWeek(i)
				nextFound = true
				break
			}
		}
		if nextFound {
			break
		}
		j++
	}
	if !nextFound {
		return 0, errors.New("NO VALID DAYS IN SCHEDULE")
	} else {
		return nextDay, nil
	}
}

//getDurationTillNextScheduleDay() returns the time.Duration to the next scheduled day.
func getDurationTillNextScheduleDay(today DaysOfTheWeek, nextDay DaysOfTheWeek) time.Duration {
	//make Duration to next scheduled day
	daysTillNext := 0
	if nextDay > today {
		daysTillNext = int(nextDay) - int(today)
	} else {
		//rest of the week, plus the next day as integer
		daysTillNext = (7 - int(today)) + int(nextDay)
	}
	tillNextDayDuration, _ := time.ParseDuration(strconv.Itoa(daysTillNext*24) + "h")
	return tillNextDayDuration
}

//GetDaysStringsToIntegers returns strings of weekdays to integers
func GetDaysStringsToIntegers(weekly model.Weekly) []DaysOfTheWeek {
	var result []DaysOfTheWeek
	var lowerCaseStringDay string
	for _, v := range weekly.Days {
		lowerCaseStringDay = strings.ToLower(v)
		dayInt := DaysMap[lowerCaseStringDay]
		result = append(result, dayInt)
	}
	return result
}

//WeeklyCheck checks all Weekly Schedules in the payload for active periods. It returns a combined ScheduleCheckerResult of all Weekly Schedules.
func WeeklyCheck(weekly model.WeeklyMap, scheduleName, timezone string) (ScheduleCheckerResult, error) {
	results := CheckWeeklyScheduleCollection(weekly, scheduleName, timezone)
	return results, nil
}
