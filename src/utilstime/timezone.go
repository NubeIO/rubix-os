package utilstime

import (
	"github.com/NubeDev/flow-framework/src/system/command"
	"strings"
)

func GetTimeZoneList() ([]string, error) {
	cmd := "timedatectl list-timezones"
	o, err := command.RunCMD(cmd, false)
	if err != nil {
		return nil, err
	}
	list := strings.Split(string(o), "\n")
	return list, nil
}
