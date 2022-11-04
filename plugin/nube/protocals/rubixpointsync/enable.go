package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	cron = gocron.NewScheduler(time.UTC)
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("SyncRubixToFF").Do(inst.SyncRubixToFF)
	cron.StartAsync()
	inst.StartMQTTSubscribeCOV() // this runs in a go routine with cancel on mqttCancel()
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	inst.mqttCancel()
	cron.Clear()
	return nil
}
