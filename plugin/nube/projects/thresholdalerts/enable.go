package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.thresholdalertsDebugMsg("History Threshold Alerts Plugin Enable()")
	inst.enabled = true
	inst.pluginName = name
	inst.setUUID()
	cron = gocron.NewScheduler(time.Local)
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("ProcessThresholdAlerts").Do(inst.ProcessThresholdAlerts)
	cron.StartAsync()
	_, next := cron.NextRun()
	inst.thresholdalertsDebugMsg("Next CRON job @ ", next.String())
	return nil
}

func (inst *Instance) Disable() error {
	inst.thresholdalertsDebugMsg("History Threshold Alerts Plugin Disable()")
	inst.enabled = false
	cron.Clear()
	return nil
}
