package main

import (
	"github.com/NubeIO/flow-framework/src/jobs"
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	j := new(jobs.Jobs)
	j.InitCron()
	cron = gocron.NewScheduler(time.UTC)
	cron.StartAsync()
	_, err := cron.Every(i.config.Job.Frequency).Tag("SyncInflux").Do(i.syncInflux)
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Disable() error {
	i.enabled = false
	cron.Clear()
	return nil
}
