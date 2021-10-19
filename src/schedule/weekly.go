package schedule

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"log"
	"strconv"
	"strings"
	"time"
)

type SchedJSON struct {
	Weekly  WeeklyScheduleCollection
	Events  interface{}
	Holiday interface{}
}

type WeeklyScheduleCollection map[string]WeeklyScheduleEntry

type WeeklyScheduleEntry struct {
	Name     string
	Days     []string
	DaysNums []DaysOfTheWeek
	Start    string
	End      string
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

var DaysMap = map[string]DaysOfTheWeek{
	"sunday":    0,
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,
}

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

//CheckWeeklyScheduleEntryWithEntryTimezone checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period ignores entry.Timezone and uses Local timezone.
func CheckWeeklyScheduleEntryWithEntryTimezone(entry WeeklyScheduleEntry) WeeklyScheduleCheckerResult {
	return CheckWeeklyScheduleEntry(entry, entry.Timezone)
}

//CheckWeeklyScheduleEntry checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleEntry(entry WeeklyScheduleEntry, checkTimezone string) WeeklyScheduleCheckerResult {
	result := WeeklyScheduleCheckerResult{}
	result.Payload = entry.Value
	result.IsActive = false

	//get time.Location for entry timezone and check timezone
	location, err := time.LoadLocation(checkTimezone)
	if err != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Timezone")
		return result
	}

	now := time.Now().In(location)
	result.CheckTime = now.Unix()
	//get day of week and compare with entry.DaysNums
	//nowHour, nowMinute, nowSecond := now.Clock()
	nowYear, nowMonth, nowDate := now.Date()
	log.Println("nowYear: ", nowYear, "nowMonth: ", nowMonth, "nowDate: ", nowDate)
	nowDayOfWeek := DaysOfTheWeek(now.Weekday())
	//nowDayOfWeekString := now.String()

	//parse start and stop times
	var entryStartHour, entryStartMins, entryStopHour, entryStopMins int
	n, err1 := fmt.Sscanf(entry.Start, "%d:%d", &entryStartHour, &entryStartMins)
	m, err2 := fmt.Sscanf(entry.End, "%d:%d", &entryStopHour, &entryStopMins)
	if n != 2 || m != 2 || err1 != nil || err2 != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Start/Stop Time")
	}
	log.Println("entryStartHour: ", entryStartHour, "entryStartMins: ", entryStartMins, "entryStopHour: ", entryStopHour, "entryStopMins: ", entryStopMins)

	//parse start and end time into current day timestamps
	startTimestamp := time.Date(nowYear, nowMonth, nowDate, entryStartHour, entryStartMins, 0, 0, location)
	stopTimestamp := time.Date(nowYear, nowMonth, nowDate, entryStopHour, entryStopMins, 59, 0, location)

	//Check if the schedule is active today
	scheduleActiveToday := false
	for _, day := range entry.DaysNums {
		if day == nowDayOfWeek {
			scheduleActiveToday = true
			break
		}
	}
	log.Println("scheduleActiveToday: ", scheduleActiveToday)

	//find the next active schedule day
	log.Println("nowDayOfWeek: ", nowDayOfWeek, "entry.DaysNums: ", entry.DaysNums)
	nextDay, err := getNextScheduleDay(nowDayOfWeek, entry.DaysNums)
	if err != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Scheduled Days are Invalid")
	}
	nextDayDuration := getDurationTillNextScheduleDay(nowDayOfWeek, nextDay)
	log.Println("nextDay: ", nextDay, "nextDayDuration: ", nextDayDuration)

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
		nextNextDay, err := getNextScheduleDay(nextDay, entry.DaysNums)
		if err != nil {
			result.ErrorFlag = true
			result.ErrorStrings = append(result.ErrorStrings, "Critical: Scheduled Days are Invalid")
		}
		nextNextDayDuration := getDurationTillNextScheduleDay(nextDay, nextNextDay)
		result.NextStart = startTimestamp.Add(nextDayDuration).Add(nextNextDayDuration).Unix()
		result.NextStop = stopTimestamp.Add(nextDayDuration).Add(nextNextDayDuration).Unix()
	}
	log.Println("CheckWeeklyScheduleEntry RETURN: ", result)
	return result
}

//CombineWeeklyScheduleCheckerResults checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CombineWeeklyScheduleCheckerResults(current WeeklyScheduleCheckerResult, new WeeklyScheduleCheckerResult, nextStopStartRequired bool) WeeklyScheduleCheckerResult {
	log.Println("CombineWeeklyScheduleCheckerResults()")
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
	log.Println("overlap: ", overlap)

	//IsActive
	result.IsActive = current.IsActive || new.IsActive

	//CheckTime
	if current.CheckTime < new.CheckTime {
		result.CheckTime = current.CheckTime
	} else {
		result.CheckTime = new.CheckTime
	}

	//PeriodStart
	if current.PeriodStart <= new.PeriodStart && current.PeriodStart != 0 {
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
		if current.PeriodStart < new.PeriodStart && current.PeriodStart != 0 {
			result.PeriodStop = current.PeriodStop
		} else {
			result.PeriodStop = new.PeriodStop
		}
	}

	//Payload
	if current.PeriodStart < new.PeriodStart && current.PeriodStart != 0 {
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
		//NextStart and NextStop
		var currentNextStartTime int64
		var newNextStartTime int64
		var currentNextStopTime int64
		var newNextStopTime int64

		if current.IsActive {
			currentNextStartTime = current.NextStart
			currentNextStopTime = current.NextStop
		} else {
			currentNextStartTime = current.PeriodStart
			currentNextStopTime = current.PeriodStop
		}

		if new.IsActive {
			newNextStartTime = new.NextStart
			newNextStopTime = new.NextStop
		} else {
			newNextStartTime = new.PeriodStart
			newNextStopTime = new.PeriodStop
		}

		//select NextStart
		if currentNextStartTime <= newNextStartTime && current.PeriodStart != 0 {
			result.NextStart = currentNextStartTime
		} else {
			result.NextStart = newNextStartTime
		}

		//Check if next periods overlap
		nextOverlap := false
		if (newNextStartTime >= currentNextStartTime && newNextStartTime <= currentNextStopTime) || (newNextStopTime >= currentNextStartTime && newNextStopTime <= currentNextStopTime) {
			nextOverlap = true
		}
		log.Println("nextOverlap: ", nextOverlap)

		if nextOverlap {
			if currentNextStopTime >= newNextStopTime {
				result.NextStop = currentNextStopTime
			} else {
				result.NextStop = newNextStopTime
			}
		} else { // no overlap so find which period is first
			if currentNextStartTime < newNextStartTime && currentNextStartTime != 0 {
				result.NextStop = currentNextStopTime
			} else {
				result.NextStop = newNextStopTime
			}
		}
	}
	return result
}

//CheckWeeklyScheduleCollection checks if there is a WeeklyScheduleEntry in the provided WeeklyScheduleCollection that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleCollection(scheduleMap WeeklyScheduleCollection, scheduleName string) WeeklyScheduleCheckerResult {
	finalResult := WeeklyScheduleCheckerResult{}
	var singleResult WeeklyScheduleCheckerResult
	count := 0
	for i, scheduleEntry := range scheduleMap {
		if scheduleEntry.Name == scheduleName {
			scheduleEntry = ConvertDaysStringsToInt(scheduleEntry)
			log.Println("WEEKLY SCHEDULE ", i, ": ", scheduleEntry)
			//singleResult = CheckWeeklyScheduleEntry(scheduleEntry, "Australia/Sydney")
			singleResult = CheckWeeklyScheduleEntryWithEntryTimezone(scheduleEntry)
			log.Println("finalResult ", finalResult, "singleResult: ", singleResult)
			if count == 0 {
				finalResult = singleResult
			} else {
				finalResult = CombineWeeklyScheduleCheckerResults(finalResult, singleResult, true)
			}
			log.Println("finalResult ", finalResult)
			count++
		}
	}
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
	//strconv.Itoa() converts int to string
	log.Println("daysTillNext: ", daysTillNext, "strconv.Itoa(daysTillNext): ", strconv.Itoa(daysTillNext))
	tillNextDayDuration, _ := time.ParseDuration(strconv.Itoa(daysTillNext*24) + "h")
	return tillNextDayDuration
}

//ConvertDaysStringsToInt converts strings of weekdays to integers
func ConvertDaysStringsToInt(weeklyScheduleEntry WeeklyScheduleEntry) WeeklyScheduleEntry {
	var lowerCaseStringDay string
	for _, v := range weeklyScheduleEntry.Days {
		lowerCaseStringDay = strings.ToLower(v)
		dayInt := DaysMap[lowerCaseStringDay]
		weeklyScheduleEntry.DaysNums = append(weeklyScheduleEntry.DaysNums, dayInt)
	}
	return weeklyScheduleEntry
}

//GetNextStartStop gets the next start and stop times from a single WeeklyScheduleCheckerResult
func GetNextStartStop(weeklyResultObj WeeklyScheduleCheckerResult) (nextStart int64, nextStop int64) {
	var start, stop int64
	if weeklyResultObj.IsActive {
		stop = weeklyResultObj.PeriodStop
		start = weeklyResultObj.NextStart
	} else {
		stop = weeklyResultObj.PeriodStop
		start = weeklyResultObj.PeriodStart
	}
	return start, stop
}

func WeeklyCheck(schedules datatypes.JSON, scheduleName string) (WeeklyScheduleCheckerResult, error) {
	var AllSchedules SchedJSON
	err := json.Unmarshal(schedules, &AllSchedules)
	if err != nil {
		log.Println("Unexpected error parsing json")
		return WeeklyScheduleCheckerResult{}, err
	}
	var AllWeeklySchedules = AllSchedules.Weekly
	for i, v := range AllWeeklySchedules {
		log.Println("WEEKLY SCHEDULE ", i, ": ", v)
	}
	results := CheckWeeklyScheduleCollection(AllWeeklySchedules, scheduleName)
	log.Println("RESULT: ", results)
	return results, nil
}

//TODO old way for reading json file
//func WeeklyCheck(file string, scheduleName string) (WeeklyScheduleCheckerResult, error) {
//	fileContentsInBytes, err := ioutil.ReadFile(file)
//	if err != nil {
//		log.Fatal(err)
//		return WeeklyScheduleCheckerResult{}, err
//	}
//	var AllSchedules SchedJSON
//	err = json.Unmarshal(fileContentsInBytes, &AllSchedules)
//	if err != nil {
//		log.Println("Unexpected error parsing json")
//		return WeeklyScheduleCheckerResult{}, err
//	}
//	var AllWeeklySchedules WeeklyScheduleCollection = AllSchedules.Weekly
//	for i, v := range AllWeeklySchedules {
//		log.Println("WEEKLY SCHEDULE ", i, ": ", v)
//	}
//	results := CheckWeeklyScheduleCollection(AllWeeklySchedules, scheduleName)
//	log.Println("RESULT: ", results)
//	return results, nil
//}

func main() {
	//var wg sync.WaitGroup
	log.Println("Starting Weekly Checks")
	//WeeklyCheck("./schedule/weekly_schedule.json", "TEST")
	/*
		go func() {
			wg.Add(1)
			WeeklyCheck("./schedule/weekly_schedule.json", "TEST")
			wg.Done()
		}()
	*/
}

func forever() {
	for {
		//fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}
