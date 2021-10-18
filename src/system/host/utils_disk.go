package host

import (
	"strings"
)

func getDiskInfo(d string, debug bool) DiskInfo {
	var disk DiskInfo
	// Size  Used Use%
	sh := "df -hl / | awk 'NR==2{print $1, $2, $3, $5}'"
	res, e := cmdRun(sh, debug)
	disk_info := strings.Fields(string(res))
	if e != nil || len(disk_info) < 4 {
		disk.DiskMount = "/"
		disk.DiskUsage = "0"
		disk.DiskUsed = "0"
		disk.DiskAll = "0"
	} else {
		disk.DiskMount = disk_info[0]
		disk.DiskAll = strings.Trim(disk_info[1], "\n")
		disk.DiskUsed = strings.Trim(disk_info[2], "\n")
		disk.DiskUsage = strings.Trim(disk_info[3], "%")
	}
	return disk
}

func getDiskInfoDetail(debug bool) DiskInfoDetail {
	var disk DiskInfoDetail

	sh := "fdisk -l|grep Disklabel"
	res, e := cmdRun(sh, debug)
	if e != nil {
		disk.Label = "unknown"
	} else {
		disk.Label = strings.Split(string(res), ":")[1]
	}

	sh = "fdisk -l|grep 'model'"
	res, e = cmdRun(sh, debug)
	if e != nil {
		disk.Model = "unknown"
	} else {
		disk.Model = strings.Split(string(res), ":")[1]
	}

	sh = "df -h|grep -v Filesystem|tr '\n' ','"
	res, e = cmdRun(sh, debug)
	if e != nil {
		disk.List = []map[string]string{}
	} else {
		var d []map[string]string
		data := strings.Split(string(res), ",")
		for _, di := range data {
			dis := strings.Fields(di)
			if len(dis) >= 6 {
				d = append(d, map[string]string{
					"system": dis[0],
					"size":   dis[1],
					"used":   dis[2],
					"avail":  dis[3],
					"use":    dis[4],
					"mount":  dis[5],
				})
			}
		}
		disk.List = d
	}
	return disk
}
