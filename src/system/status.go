package system

import (
	"github.com/NubeIO/flow-framework/src/system/command"
	"runtime"
	"strings"
	"time"
)

type Details struct {
	HostName            string `json:"host_name"`
	User                string `json:"user"`
	HostUptime          string `json:"host_uptime"`
	FlowFrameworkUptime string `json:"flow_framework_uptime"`
}

func Info() Details {
	var s Details
	hostname, err := Hostname()
	if err != nil {
		return Details{}
	}
	name, err := Uname()
	if err != nil {
		return Details{}
	}
	up, err := HostUptime()
	if err != nil {
		return Details{}
	}

	pUptime := ProgramUptime()
	if err != nil {
		return Details{}
	}

	s.HostName = hostname
	s.User = name
	s.HostUptime = up
	s.FlowFrameworkUptime = pUptime
	return s
}

var startTime time.Time

func init() {
	startTime = time.Now()
}

func FormatDuration(d time.Duration) string {
	scale := 100 * time.Second
	for scale > d {
		scale = scale / 10
	}
	return d.Round(scale / 100).String()
}

// ProgramUptime fetches hostname
func ProgramUptime() string {
	out := time.Since(startTime)
	return FormatDuration(out)
}

// Hostname go lang program uptime
func Hostname() (result string, err error) {
	return command.Run("hostname")
}

// Uname fetches uname with '-a' parameter
// (`uname -a`)
func Uname() (result string, err error) {
	return command.Run("uname", "-a")
}

// HostUptime fetches system uptime
// (`uptime`)
func HostUptime() (result string, err error) {
	up, e := command.Run("uptime")
	if e != nil {
		return "", e
	}
	s := strings.Split(up, ",")
	if len(s) >= 1 {
		return s[0], nil
	}
	return "", nil
}

// FreeSpaces fetches disk usages
// (`df -h`)
func FreeSpaces() (result string, err error) {
	return command.Run("df", "-h")
}

// FreeMemory fetches free memory
// (`free -o -h`)
func FreeMemory() (result string, err error) {
	return command.Run("free", "-h")
}

// MemoryUsage fetches system & heap allocated memory usage
func MemoryUsage() (sys, heap uint64) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return m.Sys, m.HeapAlloc
}

// CpuTemperature fetches CPU temperature
// (`vcgencmd measure_temp`)
func CpuTemperature() (result string, err error) {
	result, err = command.Run("vcgencmd", "measure_temp")
	if err == nil {
		comps := strings.Split(result, "=") // eg: "temp=68.0'C"
		if len(comps) == 2 {
			return comps[1], nil
		}
	}
	return result, err
}
