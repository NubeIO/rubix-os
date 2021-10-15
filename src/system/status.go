package system

// Tools for retrieving various statuses of Raspberry Pi

import (
	"fmt"
	"github.com/NubeDev/flow-framework/src/system/command"
	"runtime"
	"strconv"
	"strings"
)

// Hostname fetches hostname
// (`hostname`)
func Hostname() (result string, err error) {
	return command.Run("hostname")
}

// Uname fetches uname with '-a' parameter
// (`uname -a`)
func Uname() (result string, err error) {
	return command.Run("uname", "-a")
}

// Uptime fetches system uptime
// (`uptime`)
func Uptime() (result string, err error) {
	return command.Run("uptime")
}

// FreeSpaces fetches disk usages
// (`df -h`)
func FreeSpaces() (result string, err error) {
	return command.Run("df", "-h")
}

// MemorySplit fetches memory split: arm and gpu
// (`vcgencmd get_mem arm; vcgencmd get_mem gpu`)
func MemorySplit() (result []string, err error) {
	var output string
	// arm memory
	output, err = command.Run("vcgencmd", "get_mem", "arm")
	result = append(result, output)
	if err == nil {
		// gpu memory
		output, err = command.Run("vcgencmd", "get_mem", "gpu")
		result = append(result, output)
	}
	return
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

// CpuFrequency fetches frequency of arm clock
// (`vcgencmd measure_clock arm`)
func CpuFrequency() (result string, err error) {
	result, err = command.Run("vcgencmd", "measure_clock", "arm")
	if err == nil {
		comps := strings.Split(result, "=") // eg: "frequency(48)=600169920"
		if len(comps) == 2 {
			num, _ := strconv.ParseFloat(strings.TrimSpace(comps[1]), 64)
			return fmt.Sprintf("%.1f MHz", num/1000.0/1000.0), nil
		}
	}
	return result, err
}

// CpuThrottled returns whether the system is throttled or not
func CpuThrottled() (result string, err error) {
	result, err = command.Run("vcgencmd", "get_throttled")
	if err == nil {
		comps := strings.Split(result, "=") // eg: throttled=0x50000
		if len(comps) == 2 {
			num, _ := strconv.ParseInt(strings.Replace(strings.TrimSpace(comps[1]), "0x", "", -1), 16, 64)

			results := []string{}

			// https://www.raspberrypi.org/forums/viewtopic.php?f=63&t=147781&start=50#p972790
			if num&1 > 0 {
				// under-voltage
				results = append(results, "under-voltage now")
			}
			if num&(1<<1) > 0 {
				// arm frequency capped
				results = append(results, "arm freq capped now")
			}
			if num&(1<<2) > 0 {
				// currently throttled
				results = append(results, "throttled now")
			}
			if num&(1<<16) > 0 && num&1 <= 0 {
				// under-voltage has occurred
				results = append(results, "under-voltage before")
			}
			if num&(1<<17) > 0 && num&(1<<1) <= 0 {
				// arm frequency capped has occurred
				results = append(results, "arm freq capped before")
			}
			if num&(1<<18) > 0 && num&(1<<2) <= 0 {
				// throttling has occurred
				results = append(results, "throttled before")
			}

			if len(results) <= 0 {
				result = "ok"
			} else {
				result = strings.Join(results, ", ")
			}

			return result, nil
		}
	}
	return result, err
}

// CpuInfo fetches CPU information
// (`cat /proc/cpuinfo`)
func CpuInfo() (result string, err error) {
	return command.Run("cat", "/proc/cpuinfo")
}
