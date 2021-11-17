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
	job, err := i.db.GetJobByPluginConfId(i.pluginUUID)
	if err != nil {
		return nil, err
	}
	if job.UUID == "" {
		m := new(model.Job)
		m.StartDate = time.Date(2000, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
		m.EndDate = time.Date(2055, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
		m.Frequency = "15s"
		m.PluginConfId = i.pluginUUID
		job, err = i.db.CreateJob(m)
		if err != nil {
			return nil, err
		}
	}
	return job, err
}
