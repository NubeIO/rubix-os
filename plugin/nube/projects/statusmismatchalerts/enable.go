package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.statusmismatchalertsDebugMsg("History StatusMismatch Alerts Plugin Enable()")
	inst.enabled = true
	inst.pluginName = name
	inst.setUUID()
	cron = gocron.NewScheduler(time.Local)
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("ProcessStatusMismatchAlerts").Do(inst.ProcessStatusMismatchAlerts)
	cron.StartAsync()
	_, next := cron.NextRun()
	inst.statusmismatchalertsDebugMsg("Next CRON job @ ", next.String())
	return nil
}

func (inst *Instance) Disable() error {
	inst.statusmismatchalertsDebugMsg("History StatusMismatch Alerts Plugin Disable()")
	inst.enabled = false
	cron.Clear()
	return nil
}
