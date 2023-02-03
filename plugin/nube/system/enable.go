package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	log.Info("SYSTEM Plugin Enable()")
	inst.enabled = true
	inst.setUUID()
	var arg api.Args
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, arg)
	if q != nil {
		inst.networkUUID = q.UUID
	} else {
		inst.networkUUID = "NA"
	}
	if err != nil {
		log.Error("system-plugin: error on enable system-plugin")
	}
	cron = gocron.NewScheduler(time.UTC)
	_, _ = cron.Every(inst.config.Schedule.Frequency).Tag("ScheduleCheck").Do(inst.runSchedule)
	cron.StartAsync()
	return nil
}

func (inst *Instance) Disable() error {
	log.Info("SYSTEM Plugin Disable()")
	inst.enabled = false
	cron.Clear()
	return nil
}
