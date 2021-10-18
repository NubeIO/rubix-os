package host

import (
	"fmt"
)

func getLogData(log string, debug bool) string {
	if log == "" {
		return ""
	}
	sh := fmt.Sprintf("tail -n 100 %s", log)
	res, e := CMDRun(sh, debug)
	if e != nil {
		return ""
	}

	return string(res)
}

func delLog(log string, debug bool) string {
	if log == "" {
		return "ok"
	}
	sh := fmt.Sprintf(":> %s", log)
	_, e := CMDRun(sh, debug)
	if e != nil {
		return "fail"
	}

	return "ok"
}
