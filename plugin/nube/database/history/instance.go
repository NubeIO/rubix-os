package main

import (
	"github.com/NubeIO/flow-framework/model"
	"time"
)

func (i *Instance) setUUID() {
	q, err := i.db.GetPluginByPath(path)
	if err != nil {
		return
	}
	i.pluginUUID = q.UUID
}

func (i *Instance) createJob() (*model.Job, error) {
	i.REST = nil
	jobs, err := i.db.GetJobsByPluginConfigId(i.pluginUUID)
	if err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		m := new(model.Job)
		m.StartDate = time.Date(2000, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
		m.EndDate = time.Date(2055, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
		m.Frequency = "1m"
		m.PluginConfId = i.pluginUUID
		return i.db.CreateJob(m)
	}
	return jobs[0], nil
}
