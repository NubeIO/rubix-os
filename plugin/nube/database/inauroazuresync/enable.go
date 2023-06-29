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
	// inst.clearPluginConfStorage() // this may cause the system to re-send previously sent histories
	// start periodic functions
	cron = gocron.NewScheduler(time.UTC)
	// cron.SetMaxConcurrentJobs(2, gocron.RescheduleMode)
	cron.SetMaxConcurrentJobs(1, gocron.WaitMode)

	sensorFreqDuration, err := time.ParseDuration(inst.config.Job.SensorHistorySyncFrequency)
	if err == nil {
		_, _ = cron.Every(sensorFreqDuration).Tag("SensorHistorySync").Do(inst.syncAzureSensorHistories)
	} else if inst.config.Job.SensorHistorySyncFrequency == "" {
		inst.inauroazuresyncErrorMsg(`SENSOR PAYLOAD DISABLED.  Plugin/Module config 'sensor_history_sync_frequency' = ""`)
	} else {
		inst.inauroazuresyncErrorMsg(fmt.Sprintf("Invalid `sensor_history_sync_frequency` in plugin/module config.  Default frequency is used '%v'.", defaultSensorHistorySyncFrequency))
		_, _ = cron.Every(defaultSensorHistorySyncFrequency).Tag("SensorHistorySync").Do(inst.syncAzureSensorHistories)
	}

	gatewayFreqDuration, err := time.ParseDuration(inst.config.Job.GatewayPayloadSyncFrequency)
	if err == nil {
		_, _ = cron.Every(gatewayFreqDuration).Tag("GatewayPayloadSync").Do(inst.syncAzureGatewayPayloads)
	} else if inst.config.Job.GatewayPayloadSyncFrequency == "" {
		inst.inauroazuresyncErrorMsg(`GATEWAY PAYLOAD DISABLED.  Plugin/Module config 'gateway_payload_sync_frequency' = ""`)
	} else {
		inst.inauroazuresyncErrorMsg(fmt.Sprintf("Invalid `gateway_payload_sync_frequency` in plugin/module config.  Default frequency is used '%v'.", defaultGatewayPayloadSyncFrequency))
		_, _ = cron.Every(defaultGatewayPayloadSyncFrequency).Tag("GatewayPayloadSync").Do(inst.syncAzureGatewayPayloads)
	}

	cron.RunAll()
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
