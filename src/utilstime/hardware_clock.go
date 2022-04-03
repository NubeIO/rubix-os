package utilstime

import (
	"github.com/NubeIO/flow-framework/src/system/command"
	"strings"
)

type HardwareClock struct {
	Localtime               string `json:"localtime"`
	UniversalTime           string `json:"utc_time"`
	RTCtime                 string `json:"rtc_time"`
	Timezone                string `json:"timezone"`
	SystemClockSynchronized string `json:"system_clock_synchronized"`
	NTPService              string `json:"ntp_service"`
	RTCInLocalTZ            string `json:"rtc_in_local_tz"`
}

func GetHardwareTZ() (string, error) {
	cmd := "cat /etc/timezone"
	o, err := command.RunCMD(cmd, false)
	if err != nil {
		return "", err
	}
	out := strings.Split(string(o), "\n")
	if len(out) >= 0 {
		return out[0], err
	} else {
		return "", err
	}

}

func GetHardwareClock() (HardwareClock, error) {
	cmd := "timedatectl status"
	o, err := command.RunCMD(cmd, false)
	var hc HardwareClock
	if err != nil {
		return hc, err
	}
	var items []string
	list := strings.Split(string(o), "\n")
	for _, s := range list {

		items = append(items, clean(s))
	}
	if len(items) >= 6 {
		hc.Localtime = items[0]
		hc.UniversalTime = items[1]
		hc.RTCtime = items[2]
		hc.Timezone = items[3]
		hc.SystemClockSynchronized = items[4]
		hc.NTPService = items[5]
		hc.RTCInLocalTZ = items[6]
	}
	return hc, nil
}
