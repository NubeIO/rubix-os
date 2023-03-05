package main

import (
	"github.com/NubeIO/flow-framework/src/jobs"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) schedule() {
	j, ok := jobs.GetJobService()
	if ok {
		_, err := j.Every(30).Second().Do(inst.runSchedule)
		if err != nil {
			log.Debugf("system-plugin-schedule: error on create job %v\n", err)
		}
	}
}
