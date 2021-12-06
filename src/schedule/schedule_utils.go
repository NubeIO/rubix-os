package schedule

import (
	"encoding/json"
	"errors"
	"gorm.io/datatypes"
	"reflect"
)

type SchTypes struct {
	Weekly     TypeWeekly
	Events     interface{}
	Exceptions interface{}
}

type TypeWeekly map[string]WeeklyScheduleEntry

type TypeEvents map[string]EventScheduleEntry

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

type EventScheduleEntry struct {
	Name   string
	Dates  []EventStartStopTimestamps
	Start  string
	End    string
	Value  float64
	Colour string
}

type EventStartStopTimestamps struct {
	Start string
	End   string
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

//ScheduleCheckerResult this type defines the return values of any type of schedule.
//The `Period` is the first scheduled period, and the `Next` is the following scheduled period; these 2 periods cannot overlap.
//If ScheduleCheckerResults are combined and their scheduled periods overlap, they will be combined and the unnecessary timestamps dropped.
type ScheduleCheckerResult struct {
	IsActive     bool     `json:"is_active"`
	Payload      float64  `json:"payload"`
	PeriodStart  int64    `json:"period_start"` //unix timestamp in seconds
	PeriodStop   int64    `json:"period_stop"`  //unix timestamp in seconds
	NextStart    int64    `json:"next_start"`   //unix timestamp in seconds.  Start time for the following scheduled period.
	NextStop     int64    `json:"next_stop"`    //unix timestamp in seconds   End time for the following scheduled period.
	CheckTime    int64    `json:"check_time"`   //unix timestamp in seconds
	ErrorFlag    bool     `json:"error_flag"`
	AlertFlag    bool     `json:"alert_flag"`
	ErrorStrings []string `json:"error_strings"`
}

func (existing ScheduleCheckerResult) CopyScheduleCheckerResult() ScheduleCheckerResult {
	result := existing
	result.ErrorStrings = make([]string, len(existing.ErrorStrings))
	return result
}

func (existing ScheduleCheckerResult) CheckIfEquals(other ScheduleCheckerResult) bool {
	return reflect.DeepEqual(existing, other)
}

func DecodeSchedule(schedules datatypes.JSON) (SchTypes, error) {
	var AllSchedules SchTypes
	err := json.Unmarshal(schedules, &AllSchedules)
	if err != nil {
		//log.Println("Unexpected error parsing json")
		return SchTypes{}, err
	}
	return AllSchedules, nil
}

// THE FOLLOWING FUNCTION NEEDS REVIEW AND TESTING

//CombineScheduleCheckerResults Combines 2 ScheduleCheckerResults into a single ScheduleCheckerResult, calculating PeriodStart, PeriodStop, NextStart, and NextStop times of the combined ScheduleCheckerResult.
func CombineScheduleCheckerResults(current ScheduleCheckerResult, new ScheduleCheckerResult, nextStopStartRequired bool) (ScheduleCheckerResult, error) {
	//log.Println("CombineScheduleCheckerResults()")
	result := ScheduleCheckerResult{}
	var err error = nil

	//Check for empty ScheduleCheckerResult periods
	currentEmptyPeriod := false
	newEmptyPeriod := false
	if current.PeriodStart == 0 || current.PeriodStop == 0 {
		currentEmptyPeriod = true
	}
	if new.PeriodStart == 0 || new.PeriodStop == 0 {
		currentEmptyPeriod = true
	}
	if currentEmptyPeriod && newEmptyPeriod {
		err = errors.New("no valid schedule periods found")
		return result, err
	} else if currentEmptyPeriod { //return the valid (new) ScheduleCheckerResult
		err = errors.New("no valid schedule periods found in `current` ScheduleCheckerResult")
		return new, err
	} else if newEmptyPeriod { //return the valid (current) ScheduleCheckerResult
		err = errors.New("no valid schedule periods found in `new` ScheduleCheckerResult")
		return current, err
	}

	//AlertFlag & ErrorFlag & ErrorStrings
	if new.ErrorFlag {
		err = errors.New("`new` ScheduleCheckerResult has an ErrorFlag, cannot combine")
		return current, err
	}
	if current.ErrorFlag {
		err = errors.New("`current` ScheduleCheckerResult has an ErrorFlag, cannot combine")
		return new, err
	}
	result.AlertFlag = current.AlertFlag || new.AlertFlag
	result.ErrorStrings = append(current.ErrorStrings, new.ErrorStrings...)

	//Find order of periods and next
	currentPeriod := 0
	currentNext := 0
	newPeriod := 0
	newNext := 0

	//Find which Period is first, then logic the rest of the order (only looking at Start Times)
	if current.PeriodStart <= new.PeriodStart && (new.NextStart == 0 || current.PeriodStart <= new.NextStart) {
		currentPeriod = 1 //current period is first
		if current.NextStart != 0 && current.NextStart <= new.PeriodStart {
			currentNext = 2
			newPeriod = 3
			if new.NextStart != 0 {
				newNext = 4
			}
		} else {
			newPeriod = 2
			if (current.NextStart != 0 && new.NextStart == 0) || (current.NextStart != 0 && current.NextStart <= new.NextStart) {
				currentNext = 3
				if new.NextStart != 0 {
					newNext = 4
				}
			} else if (new.NextStart != 0 && current.NextStart == 0) || (new.NextStart != 0 && new.NextStart <= current.NextStart) {
				newNext = 3
				if current.NextStart != 0 {
					currentNext = 4
				}
			}
		}
	} else if new.PeriodStart <= current.PeriodStart && new.PeriodStart <= current.NextStart {
		newPeriod = 1 //new period is first
		if new.NextStart != 0 && new.NextStart <= current.PeriodStart {
			newNext = 2
			currentPeriod = 3
			if current.NextStart != 0 {
				currentNext = 4
			}
		} else {
			currentPeriod = 2
			if (new.NextStart != 0 && current.NextStart == 0) || (new.NextStart != 0 && new.NextStart <= current.NextStart) {
				newNext = 3
				if current.NextStart != 0 {
					currentNext = 4
				}
			} else if (current.NextStart != 0 && new.NextStart == 0) || (current.NextStart != 0 && current.NextStart <= new.NextStart) {
				currentNext = 3
				if new.NextStart != 0 {
					newNext = 4
				}
			}
		}
	}

	//Check if first period overlaps with other scheduled periods. Note that a Period or Next interval from the same ScheduleCheckerResult cannot overlap.
	secondPeriodOverlap := false
	thirdPeriodOverlap := false
	if currentPeriod == 1 {
		if newPeriod == 2 && new.PeriodStart < current.PeriodStop {
			secondPeriodOverlap = true
			if newNext == 3 && new.NextStart < current.PeriodStop {
				thirdPeriodOverlap = true
			}
		}
	} else if newPeriod == 1 {
		if currentPeriod == 2 && current.PeriodStart < new.PeriodStop {
			secondPeriodOverlap = true
			if currentNext == 3 && current.NextStart < new.PeriodStop {
				thirdPeriodOverlap = true
			}
		}
	}

	//COMPLETE TO HERE!!!!!!!!!!!!!!

	//if (new.PeriodStart >= current.PeriodStart && new.PeriodStart <= current.PeriodStop) || (new.PeriodStop >= current.PeriodStart && new.PeriodStop <= current.PeriodStop) {
	if (new.PeriodStart >= current.PeriodStart && new.PeriodStart <= current.PeriodStop) || (new.PeriodStop >= current.PeriodStart && new.PeriodStop <= current.PeriodStop) || (current.PeriodStart >= new.PeriodStart && current.PeriodStart <= new.PeriodStop) {
		periodOverlap = true
	}
	//log.Println("periodOverlap: ", periodOverlap)

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
	if periodOverlap {
		if current.PeriodStop >= new.PeriodStop {
			result.PeriodStop = current.PeriodStop
		} else {
			result.PeriodStop = new.PeriodStop
		}
	} else { // no periodOverlap so find which period is first
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

		/*
			//For ScheduleCheckerResult that don't have Next times assigned
			if !periodOverlap {
				if new.PeriodStart >= current.PeriodStart {
					result.NextStart = new.PeriodStart
					result.NextStop = new.PeriodStop
				} else {
					result.NextStart = current.PeriodStart
					result.NextStop = current.PeriodStop
				}
			}

		*/

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
		if (newNextStartTime >= currentNextStartTime && newNextStartTime <= currentNextStopTime) || (newNextStopTime >= currentNextStartTime && newNextStopTime <= currentNextStopTime) || (currentNextStartTime >= newNextStartTime && currentNextStartTime <= newNextStopTime) {
			nextOverlap = true
		}
		//log.Println("nextOverlap: ", nextOverlap)

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
}
