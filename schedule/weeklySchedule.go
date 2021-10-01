package schedule

import (
	"encoding/json"
	"fmt"
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

type WeeklyScheduleCheckerResult struct {
	IsActive     bool
	Payload      float64
	PeriodStart  uint64 //unix timestamp as from Date()
	PeriodStop   uint64 //unix timestamp as from Date()
	NextStart    uint64 //unix timestamp as from Date().  Start time for the following scheduled period.
	NextStop     uint64 //unix timestamp as from Date()   End time for the following scheduled period.
	CheckTime    uint64 //unix timestamp as from Date()
	ErrorFlag    bool
	AlertFlag    bool
	ErrorStrings []string
}

//CheckWeeklyScheduleEntry checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CheckWeeklyScheduleEntry(entry WeeklyScheduleEntry, timezone string) WeeklyScheduleCheckerResult {
	result := WeeklyScheduleCheckerResult{}
	//get local time parts in locale of entry.Timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		result.ErrorFlag = true
		result.ErrorStrings = append(result.ErrorStrings, "Critical: Invalid Timezone")
		return result
	}
	now := time.Now().In(loc)

	//get day of week and compare with entry.Days
	nowHour, nowMinute, nowSecond := now.Clock()
	nowDayOfWeek := now.Day()
	nowDayOfWeekString := now.String()
	found := false
	for day, _ = range entry.Days {
		if day == nowDayOfWeek {
			found = true
			break
		}
	}
	if found {
		//parse start and end time into current day timestamps
		//choose format with only date, and sub in hours and minutes from event
		https://golang.org/src/time/format.go
		func (t Time) Format(layout string) string
		https://pkg.go.dev/time#ParseInLocation
		func ParseInLocation(layout, value string, loc *Location) (Time, error)

		//check if now() is between the start and end timestamps



	}


	//
}

//CombineWeeklyScheduleCheckerResults checks if there is a WeeklyScheduleEntry that matches the specified schedule Name and is currently within the scheduled period.
func CombineWeeklyScheduleCheckerResults(current WeeklyScheduleCheckerResult, new WeeklyScheduleCheckerResult, nextStopStartRequired bool) WeeklyScheduleCheckerResult {
	result := WeeklyScheduleCheckerResult{}

	//AlertFlag & ErrorFlag & ErrorStrings
	if new.ErrorFlag {
		return current.ErrorFlag
	}
	if current.ErrorFlag {
		return new.ErrorFlag
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
		var currentNextStartTime uint64
		var newNextStartTime uint64

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
		var currentNextStopTime uint64
		var newNextStopTime uint64

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
