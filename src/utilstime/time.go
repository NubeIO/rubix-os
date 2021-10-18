package utilstime

import "time"

type Time struct {
	DateStamp       time.Time     `json:"date_stamp"`
	TimeLocal       string        `json:"time_local"`
	TimeUTC         string        `json:"time_utc"`
	CurrentDay      string        `json:"current_day"`
	CurrentDayUTC   string        `json:"current_day_utc"`
	DateFormatLocal string        `json:"date_format_local"`
	DateFormatUTC   string        `json:"date_format_utc"`
	SystemTimeZone  string        `json:"system_time_zone"`
	HardwareClock   HardwareClock `json:"hardware_clock"`
}

func SystemTime() *Time {
	t := new(Time)
	t.DateStamp = time.Now()
	timeUTC := t.DateStamp.UTC()
	t.TimeLocal = t.DateStamp.Format("15:04:05")
	t.TimeUTC = timeUTC.Format("15:04:05")
	t.CurrentDay = t.DateStamp.Format("Monday")
	t.CurrentDayUTC = timeUTC.Format("Monday")
	t.DateFormatLocal = t.DateStamp.Format("01-02-2006 15:04:05")
	t.DateFormatUTC = timeUTC.Format("01-02-2006 15:04:05")
	zone, _ := t.DateStamp.Zone()
	t.SystemTimeZone = zone
	tz, _ := GetHardwareClock()
	t.HardwareClock = tz
	return t
}
