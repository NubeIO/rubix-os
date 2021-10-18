package host

import (
	"errors"
	"fmt"
	"strings"
)

func getProgressData(debug bool) ProgressInfo {
	var pro ProgressInfo

	sh := "ps ax | wc -l;ps ax | awk '{print $3}' | grep R | wc -l;ps ax | awk '{print $3}' | grep Z | wc -l;ps ax | awk '{print $3}' | grep S | wc -l"
	res, e := cmdRun(sh, debug)
	if e != nil {
		pro.ProgressAll = "0"
		pro.ProgressRun = "0"
		pro.ProgressDead = "0"
		pro.ProgressSleep = "0"
	} else {
		d := strings.Split(string(res), "\n")
		if len(d) < 4 {
			pro.ProgressAll = "0"
			pro.ProgressRun = "0"
			pro.ProgressDead = "0"
			pro.ProgressSleep = "0"
		} else {
			pro.ProgressAll = strings.Trim(d[0], "\n")
			pro.ProgressRun = strings.Trim(d[1], "\n")
			pro.ProgressDead = strings.Trim(d[2], "\n")
			pro.ProgressSleep = strings.Trim(d[3], "\n")
		}
	}

	return pro
}

func fmtData(s string) (ProgressListInfo, error) {
	data := strings.Fields(s)
	if len(data) < 4 {
		return ProgressListInfo{}, errors.New("bad progress")
	}
	return ProgressListInfo{
		PID: data[0],
		CPU: data[1],
		Mem: data[2],
		Cmd: data[3],
	}, nil
}

func getProgressListData(num string, debug bool) []ProgressListInfo {
	var pro []ProgressListInfo

	if num == "" {
		num = "10"
	}

	sh := fmt.Sprintf("ps aux | grep -v PID | awk '{print $2, $3, $4, $11}' | sort -rn -k +3 | head -n %s | tr '\n' ','", num)
	res, e := cmdRun(sh, debug)
	if e != nil {
		return []ProgressListInfo{}
	}

	list := strings.Split(strings.Trim(string(res), "\n"), ",")
	for _, l := range list {

		if p, e := fmtData(l); e == nil {
			pro = append(pro, p)
		}
	}

	return pro
}
