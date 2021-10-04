package schedule

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"io/ioutil"
	"log"
	"math"
	"time"
)

type WeeklyScheduleCollection map[string]WeeklyScheduleEntry

type WeeklyScheduleEntry struct {
	Name     string
	Days     []DaysOfTheWeek
	Start    string
	Stop      string
	Timezone string
	Value    float64
	Colour   string
}

type DaysOfTheWeek int
const (
	sunday    DaysOfTheWeek = iota // 0
	monday                         // 1
	tuesday                        // 2
	wednesday                      // 3
	thursday                       // 4
	friday                         // 5
	saturday                       // 6
)

type WeeklyScheduleCheckerResult struct {
	IsActive     bool
	Payload      float64
	PeriodStart  int64 //unix timestamp in seconds
	PeriodStop   int64 //unix timestamp in seconds
	NextStart    int64 //unix timestamp in seconds.  Start time for the following scheduled period.
	NextStop     int64 //unix timestamp in seconds   End time for the following scheduled period.
	CheckTime    int64 //unix timestamp in seconds
	ErrorFlag    bool
	AlertFlag    bool
	ErrorStrings []string
}

//CheckWeeklyScheduleEntry checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleEntry(entry WeeklyScheduleEntry, timezone string) WeeklyScheduleCheckerResult {
	result := WeeklyScheduleCheckerResult{}
	result.payload = entry.Value

	//get local time parts in locale of entry.Timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Timezone")
		return result
	}
	now := time.Now().In(loc)

	//get day of week and compare with entry.Days
	//nowHour, nowMinute, nowSecond := now.Clock()
	nowYear, nowMonth, nowDate := now.Date()
	nowDayOfWeek := DaysOfTheWeek(now.Weekday())
	nowDayOfWeekString := now.String()

	scheduleActiveToday := false
	for _, day := range entry.Days {
		if day == nowDayOfWeek {
			scheduleActiveToday = true
			break
		}
	}

	//find the next active schedule day

	//if only 1 day in schedule
	/*
	if len(entry.Days) == 1 {
		oneWeek := time.ParseDuration("168h")
		result.NextStart = startTimestamp.Add(oneWeek).Unix()
		result.NextStart = stopTimestamp.Add(oneWeek).Unix()
	}
	*/

	//TODO: Turn this into a function (returns the int of the next scheduled day)
	i := int(nowDayOfWeek)
	j := 0
	var nextDay DaysOfTheWeek = 0
	nextFound := false
	//check until the next scheduled day is found in entry.Day
	for j < 7 {
		//check the next day (rollover at end of week (Saturday/6)
		if i == 6{
			i = 0
		} else {
			i++
		}
		//check each of the scheduled days in entry
		for x := 0; x <= len(entry.Days); x++ {
			if entry.Days[x] == DaysOfTheWeek(i) {
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

	//TODO: Turn this into a function (returns the Duration to the next scheduled day) takes (now, nextDayInt)
	//make Duration to next scheduled day
	daysTillNext := 0
	if nextDay > nowDayOfWeek {
		daysTillNext = int(nextDay) - int(nowDayOfWeek)
	} else {
		//rest of the week, plus the next day as integer
		daysTillNext = (7 - int(nowDayOfWeek)) + int(nextDay)
	}
	//strconv.Itoa() converts int to string
	tillNextDayDuration, _ := time.ParseDuration(strconv.Itoa(daysTillNext)+"h")

	//parse start and stop times
	var entryStartHour, entryStartMins, entryStopHour, entryStopMins int
	n, err1 := fmt.Sscanf(entry.Start, "%d:%d", &entryStartHour, &entryStartMins)
	m, err2 := fmt.Sscanf(entry.Stop, "%d:%d", &entryStopHour, &entryStopMins)
	if n != 2 || m != 2 || err1 != nil || err2 != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Start/Stop Time")
	}

	if scheduleActiveToday {
		//parse start and end time into current day timestamps
		startTimestamp := time.Date(nowYear, nowMonth, nowDate, entryStartHour, entryStartMins, 0, 0, loc)
		stopTimestamp := time.Date(nowYear, nowMonth, nowDate, entryStopHour, entryStopMins, 59, 0, loc)

		//check if today's schedule is currently active
		activeNow := false
		if now.After(startTimestamp) && now.Before(stopTimestamp) {
			activeNow = true
		}
		result.IsActive = activeNow

		if now.After(stopTimestamp) {
			//PeriodStartStop
			//NextStartStop
		}

		if activeNow  {
			//PeriodStartStop
			//NextStartStop
		} else if now.After(stopTimestamp)


		result.PeriodStart = startTimestamp.Unix()
		result.PeriodStop = stopTimestamp.Unix()
		result.CheckTime = now.Unix()



		//get timestamps for the next scheduled day
		if activeNow  {
			result.NextStart = startTimestamp.Add(tillNextDayDuration).Unix()
			result.NextStop = stopTimestamp.Unix()
		} else if now.After(stopTimestamp)

		tillNextDayDuration




	}


	//
}

//CombineWeeklyScheduleCheckerResults checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CombineWeeklyScheduleCheckerResults(current WeeklyScheduleCheckerResult, new WeeklyScheduleCheckerResult, nextStopStartRequired bool) WeeklyScheduleCheckerResult {
	result := WeeklyScheduleCheckerResult{}

	//AlertFlag & ErrorFlag & ErrorStrings
	if new.ErrorFlag {
		return current
	}
	if current.ErrorFlag {
		return new
	}
	result.AlertFlag = current.AlertFlag || new.AlertFlag
	result.ErrorStrings = append(current.ErrorStrings, new.ErrorStrings...)

	//Check if schedule periods overlap
	overlap := false
	if (new.PeriodStart >= current.PeriodStart && new.PeriodStart <= current.PeriodStop) || (new.PeriodStop >= current.PeriodStart && new.PeriodStop <= current.PeriodStop) {
		overlap = true
	}

	//IsActive
	result.IsActive = current.IsActive || new.IsActive

	//PeriodStart
	if current.PeriodStart <= new.PeriodStart {
		result.PeriodStart = current.PeriodStart
	} else {
		result.PeriodStart = new.PeriodStart
	}

	//PeriodStop
	if overlap {
		if current.PeriodStop >= new.PeriodStop {
			result.PeriodStop = current.PeriodStop
		} else {
			result.PeriodStop = new.PeriodStop
		}
	} else { // no overlap so find which period is first
		if current.PeriodStart < new.PeriodStart {
			result.PeriodStop = current.PeriodStart
		} else {
			result.PeriodStop = new.PeriodStart
		}
	}

	//Payload
	if current.PeriodStart < new.PeriodStart {
		result.Payload = current.Payload
	} else {
		result.Payload = new.Payload
	}
	if (current.IsActive && new.IsActive) && (current.Payload != 0 && new.Payload != 0) {
		result.AlertFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Multiple Payload Values For Active Schedule Period")
	}

	if !nextStopStartRequired {
		return result
	} else {
		//NextStart
		var currentNextStartTime int64
		var newNextStartTime int64

		if current.IsActive {
			currentNextStartTime = current.NextStart
		} else {
			currentNextStartTime = current.PeriodStart
		}

		if new.IsActive {
			newNextStartTime = new.NextStart
		} else {
			newNextStartTime = new.PeriodStart
		}

		if currentNextStartTime <= newNextStartTime {
			result.NextStart = currentNextStartTime
		} else {
			result.NextStart = newNextStartTime
		}

		//NextStop
		var currentNextStopTime int64
		var newNextStopTime int64

		if current.IsActive {
			currentNextStopTime = current.PeriodStop
		} else {
			currentNextStopTime = current.NextStop
		}

		if new.IsActive {
			newNextStopTime = new.PeriodStop
		} else {
			newNextStopTime = new.NextStop
		}

		if currentNextStopTime <= newNextStopTime {
			result.NextStop = currentNextStopTime
		} else {
			result.NextStop = newNextStopTime
		}
	}
	return result
}

//CheckWeeklyScheduleCollection checks if there is a WeeklyScheduleEntry in the provided WeeklyScheduleCollection that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleCollection(scheduleMap WeeklyScheduleCollection, scheduleName string) WeeklyScheduleCheckerResult {
	finalResult := WeeklyScheduleCheckerResult{}
	var singleResult WeeklyScheduleCheckerResult
	for uuid, scheduleEntry := range scheduleMap {
		if scheduleEntry.Name == scheduleName {
			singleResult = CheckWeeklyScheduleEntry(scheduleEntry)
			finalResult = CombineWeeklyScheduleCheckerResults(finalResult, singleResult, true)
		}
	}
	return finalResult
}

func EventCheckLoop(file string) {
	json, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	sd, err := schedule.New(string(json))
	if err != nil {
		log.Println("Unexpected error parsing json")
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go func() {
		for dt := range ticker.C {
			if sd.Within(dt) {
				log.Println(sd.GetDescription(), dt, " is TRUE schedule.")
			} else {
				log.Println(sd.GetDescription(), dt, " is FALSE schedule.")
			}
		}
	}()
	time.Sleep(time.Hour * 1)
	ticker.Stop()
}

func main() {
	o, oo := TestOverlappingIntervals(time.Monday)
	fmt.Println(oo)
	fmt.Println(o, "is an overlap if its true")
	go EventCheckLoop("./schedule/test.json")
	go EventCheckLoop("./schedule/test2.json")
	go forever()
	select {} // block forever

}
func forever() {
	for {
		//fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}
