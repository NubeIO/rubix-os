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
	influxDetails := inst.initializeInfluxSettings()
	inst.influxDetails = influxDetails
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("SyncInflux").Do(inst.syncInflux, influxDetails)
	cron.StartAsync()
	inst.StartMQTTSubscribeCOV() // this runs in a go routine with cancel on mqttCancel()

	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	inst.mqttCancel()
	cron.Clear()
	for _, influxConnectionInstance := range influxConnectionInstances {
		influxConnectionInstance.client.Close()
	}
	return nil
}
