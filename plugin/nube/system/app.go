package main

import (
	"github.com/NubeIO/rubix-os/src/jobs"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) schedule() {
	j, ok := jobs.GetJobService()
	if ok {
		_, err := j.Every(30).Second().Do(inst.runSchedule)
		if err != nil {
			log.Debugf("Schedule Checks: error on create job %v\n", err)
		}
	}
}
