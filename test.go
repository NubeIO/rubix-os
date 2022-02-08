package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/src/schedule"
	"github.com/NubeIO/flow-framework/src/utilstime"
)

func main() {

	tz, err := utilstime.GetHardwareTZ()
	if err != nil {
		return
	}
	fmt.Println(tz, err)

	schedule.ScheduleTest()

	fmt.Println(schedule.ConvertToHumanDatetime(int64(1644359400)))
	//schedule.ModbusScheduleTest()
}
