package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.edgeinfluxDebugMsg("EDGEINFLUX Plugin Enable()")
	inst.enabled = true
	inst.fault = false
	inst.setUUID()
	cron = gocron.NewScheduler(time.UTC)
	influxDetails := inst.initializeInfluxSettings()
	inst.influxDetails = influxDetails
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("SyncInflux").Do(inst.syncInflux, influxDetails)
	cron.StartAsync()
	inst.StartMQTTSubscribeCOV() // this runs in a go routine with cancel on mqttCancel()
	inst.running = true
	return nil
}

func (inst *Instance) Disable() error {
	inst.edgeinfluxDebugMsg("EDGEINFLUX Plugin Disable()")
	inst.enabled = false
	if inst.mqttCancel != nil {
		inst.mqttCancel()
	}
	cron.Clear()
	for _, influxConnectionInstance := range influxConnectionInstances {
		influxConnectionInstance.client.Close()
	}
	inst.running = false
	inst.fault = false
	return nil
}
