package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()
	inst.initializePostgresSetting()
	cron = gocron.NewScheduler(time.UTC)
	cron.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("SyncPostgres").Do(inst.syncPostgres)
	cron.StartAsync()
	log.Info(fmt.Sprintf("%s enabled", name))
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	cron.Clear()
	if postgresConnectionInstance != nil {
		conn, _ := postgresConnectionInstance.db.DB()
		_ = conn.Close()
	}
	log.Info(fmt.Sprintf("%s disabled", name))
	return nil
}
