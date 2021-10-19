package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/src/jobs"
	"github.com/NubeDev/flow-framework/src/schedule"
)

func (i *Instance) task() {
	fmt.Println(1111111111)
	schedules, err := i.db.GetSchedules()
	if err != nil {
		return
	}

	for _, sch := range schedules {
		check, err := schedule.WeeklyCheck(sch.Schedules, "Branch")
		if err != nil {
			return
		}
		fmt.Println(check.IsActive)
	}

}

func (i *Instance) schedule() {

	j, ok := jobs.GetJobService()
	if ok {
		_, err := j.Every(10).Second().Tag("foo").Do(i.task)
		if err != nil {
			//return
		}
	}

}
