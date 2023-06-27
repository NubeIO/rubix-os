package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.cpsDebugMsg("CPS Plugin Enable()")
	inst.enabled = true
	inst.fault = false
	inst.running = false
	inst.pluginName = name
	inst.setUUID()
	inst.initializePostgresSetting()
	_, err := inst.initializePostgresDBConnection()
	if err != nil {
		inst.cpsErrorMsg("Enable() initializePostgresDBConnection() error: ", err)
	}

	cron = gocron.NewScheduler(time.Local)
	// cron.SetMaxConcurrentJobs(2, gocron.RescheduleMode)
	cron.SetMaxConcurrentJobs(1, gocron.WaitMode)
	_, _ = cron.Every("30m").Tag("initializePostgresDBConnection").Do(inst.initializePostgresDBConnection)
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("cpsProcessing").Do(inst.CPSProcessing)
	cron.StartAsync()
	_, next := cron.NextRun()
	inst.cpsDebugMsg("Next CRON job @ ", next.String())
	inst.running = true
	return nil
}

func (inst *Instance) Disable() error {
	inst.cpsDebugMsg("CPS Plugin Disable()")
	inst.enabled = false
	cron.Clear()
	inst.fault = false
	inst.running = false
	return nil
}
