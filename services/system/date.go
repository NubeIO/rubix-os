package system

import (
	"fmt"
	"github.com/NubeIO/lib-date/datectl"
	"github.com/NubeIO/lib-date/datelib"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

type DateBody struct {
	DateTime string `json:"date_time"`
	TimeZone string `json:"time_zone"`
}

func (inst *System) SystemTime() *datelib.Time {
	return datelib.New(&datelib.Date{}).SystemTime()
}

func (inst *System) GenerateTimeSyncConfig(body *datectl.TimeSyncConfig) string {
	return inst.datectl.GenerateTimeSyncConfig(body)
}

func (inst *System) GetHardwareTZ() (string, error) {
	return inst.datectl.GetHardwareTZ()
}

func (inst *System) GetHardwareClock() (*datectl.HardwareClock, error) {
	return inst.datectl.GetHardwareClock()
}

func (inst *System) GetTimeZoneList() ([]string, error) {
	return inst.datectl.GetTimeZoneList()
}

func (inst *System) UpdateTimezone(body DateBody) (*Message, error) {
	err := inst.datectl.UpdateTimezone(body.TimeZone)
	if err != nil {
		return nil, err
	}
	return &Message{
		Message: fmt.Sprintf("updated to %s", body.TimeZone),
	}, nil
}

func (inst *System) SetSystemTime(dateTime DateBody) (*datelib.Time, error) {
	layout := "2006-01-02 15:04:05"
	// parse time
	t, err := time.Parse(layout, dateTime.DateTime)
	if err != nil {
		return nil, fmt.Errorf("could not parse date try 2006-01-02 15:04:05 %s", err)
	}
	log.Infof("set time to %s", t.String())
	timeString := fmt.Sprintf("%s", dateTime.DateTime)
	cmd := exec.Command("date", "-s", timeString)
	output, err := cmd.Output()
	cleanCommand(string(output), cmd, err, debug)
	if err != nil {
		return nil, err
	}
	return datelib.New(&datelib.Date{}).SystemTime(), nil
}

func (inst *System) NTPEnable() (*Message, error) {
	msg, err := inst.datectl.NTPEnable()
	if err != nil {
		return nil, err
	}
	return &Message{
		Message: msg.Message,
	}, nil
}

func (inst *System) NTPDisable() (*Message, error) {
	msg, err := inst.datectl.NTPDisable()
	if err != nil {
		return nil, err
	}
	return &Message{
		Message: msg.Message,
	}, nil
}
