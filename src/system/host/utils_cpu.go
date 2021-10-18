package host

import (
	"strconv"
	"strings"
)

func equalZero(s string) string {
	if s == "0" || s == "0.0" {
		return "0"
	}
	return s
}

func getCPUInfo(debug bool) CPUInfo {
	var cpu CPUInfo
	// sh := "top -bn 1 -i -c | grep %Cpu | sed 's/%Cpu(s)://g'"
	// avg-cpu:  %user   %nice %system %iowait  %steal   %idle
	sh := "iostat -c|awk 'NR==4{print $1,$2,$3,$4,$6}'"
	res, e := CMDRun(sh, debug)
	cpu_info := strings.Split(string(res), " ")

	if e != nil || len(cpu_info) < 5 {
		cpu.CpuUsage, cpu.CpuUsageSys, cpu.CpuUsageUser = "0", "0", "0"
		cpu.CpuFree, cpu.CpuIOWait = "0", "0"
	} else {
		cpu_us := equalZero(cpu_info[0])
		cpu_sy := equalZero(cpu_info[2])
		cpu_free := equalZero(cpu_info[4])
		t, _ := strconv.Atoi(strings.Split(cpu_free, ".")[0])
		cpu_usage := strconv.Itoa(100 - t)
		cpu_io_wait := equalZero(cpu_info[3])

		cpu.CpuUsage = cpu_usage
		cpu.CpuUsageSys = cpu_sy
		cpu.CpuUsageUser = cpu_us
		cpu.CpuFree = cpu_free
		cpu.CpuIOWait = cpu_io_wait
	}

	// load
	sh = "cat /proc/loadavg | awk '{print $1,$2,$3}'"
	res, e = CMDRun(sh, debug)
	if e != nil {
		cpu.CpuLoad = "0 0 0"
	} else {
		cpu.CpuLoad = strings.Trim(string(res), "\n")
	}

	// cpu count
	sh = "cat /proc/cpuinfo | grep \"cpu cores\" | uniq | awk '{print $4}'"
	res, e = CMDRun(sh, debug)
	if e != nil {
		cpu.CpuCount = "0"
	} else {
		cpu.CpuCount = strings.Trim(string(res), "\n")
	}

	sh = "cat /proc/cpuinfo | grep \"processor\" | wc -l"
	res, e = CMDRun(sh, debug)
	if e != nil {
		cpu.CpuPhysical = "0"
	} else {
		cpu.CpuPhysical = strings.Trim(string(res), "\n")
	}

	// cpu run days
	sh = "uptime | awk '{print $3}'"
	res, e = CMDRun(sh, debug)
	if e != nil {
		cpu.CpuRun = "0"
	} else {
		cpu.CpuRun = strings.Trim(string(res), "\n")
	}

	return cpu
}

func getCpuInfoDetail(debug bool) CPUInfoDetail {
	var cpu CPUInfoDetail

	sh := "cat /proc/cpuinfo | grep name | cut -f2 -d: | uniq"
	res, e := CMDRun(sh, debug)
	if e != nil {
		cpu.Info = "unknown cpu"
	} else {
		cpu.Info = strings.Trim(strings.TrimSpace(string(res)), "\n")
	}

	sh = "cat /proc/cpuinfo | grep \"cpu MHz\" | uniq"
	res, e = CMDRun(sh, debug)
	if e != nil {
		cpu.Freq = "unknown cpu"
	} else {
		cpu.Freq = strings.Trim(strings.TrimSpace(string(res)), "\n")
	}

	sh = "cat /proc/cpuinfo | grep \"cache size\" | uniq"
	res, e = CMDRun(sh, debug)
	if e != nil {
		cpu.Cache = "unknown cpu"
	} else {
		cpu.Cache = strings.Trim(strings.TrimSpace(string(res)), "\n")
	}

	return cpu
}
