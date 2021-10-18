package host

import (
	"fmt"
	"strconv"
	"strings"
)

func getMemInfo(debug bool) MemInfo {
	var m MemInfo
	sh := "free -m | grep Mem | sed 's/Mem://g' | awk '{print $1, $2, $5, $6}'"
	res, e := cmdRun(sh, debug)
	mem_info := strings.Fields(string(res))

	if e != nil || len(mem_info) < 4 {
		return MemInfo{}
	}

	total, _ := strconv.Atoi(mem_info[1])
	used, _ := strconv.Atoi(mem_info[0])

	m.MemUsage = fmt.Sprintf("%.2f", float64(total)/float64(used)*100)
	m.MemUsed = mem_info[1] + "M"
	m.MemCache = mem_info[2] + "M"
	m.MemFree = mem_info[3] + "M"
	return m
}

func getMemInfoDetail(debug bool) MemInfoDetail {
	var mem MemInfoDetail

	sh := "dmidecode | grep -A2 \"System Information\"|grep -v \"System Information\" | tr '\n' '|'"
	res, e := cmdRun(sh, debug)
	if e != nil {
		mem.ManuFacturer = "unknown"
		mem.Product = "unknown"
	} else {
		d := strings.Split(string(res), "|")
		mem.ManuFacturer = strings.Split(strings.Trim(d[0], "\n"), ":")[1]
		mem.Product = strings.Split(strings.Trim(d[1], "\n"), ":")[1]
	}

	sh = "dmidecode -t memory|grep  -E 'Size|Speed|Total Width'|head -n 3|sort|tr '\n' '|'"
	res, e = cmdRun(sh, debug)
	if e != nil {
		mem.Size = "unknown"
		mem.Speed = "unknown"
		mem.Width = "unknown"
	} else {
		d := strings.Split(string(res), "|")
		mem.Size = strings.Split(strings.Trim(d[0], "\n"), ":")[1]
		mem.Speed = strings.Split(strings.Trim(d[1], "\n"), ":")[1]
		mem.Width = strings.Split(strings.Trim(d[2], "\n"), ":")[1]
	}

	return mem
}
