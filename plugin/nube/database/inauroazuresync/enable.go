package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.inauroazuresyncDebugMsg("INAURO AZURE SYNC Plugin Enable()")
	inst.enabled = true
	inst.fault = false
	inst.running = false
	inst.setUUID()
	// start periodic functions
	cron = gocron.NewScheduler(time.UTC)
	// cron.SetMaxConcurrentJobs(2, gocron.RescheduleMode)
	cron.SetMaxConcurrentJobs(1, gocron.WaitMode)
	// _, _ = cron.Every(inst.config.Job.GatewayPayloadSyncFrequency).Tag("GatewayPayloadSync").Do(inst.syncAzureGatewayPayloads)
	_, _ = cron.Every(inst.config.Job.SensorHistorySyncFrequency).Tag("SensorHistorySync").Do(inst.syncAzureSensorHistories)
	cron.StartAsync()
	inst.running = true
	log.Info(fmt.Sprintf("%s enabled", name))
	return nil
}

func (inst *Instance) Disable() error {
	log.Info("INAURO AZURE SYNC Plugin Disable()")
	inst.enabled = false
	if cron != nil {
		cron.Clear()
	}
	inst.running = false
	inst.fault = false
	return nil
}
