package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()
	cron = gocron.NewScheduler(time.UTC)
	postgresSetting := inst.initializePostgresSetting()
	postgresSetting, err := New(postgresSetting)
	if err != nil {
		return err
	}
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("SyncPostgres").Do(inst.syncPostgres, postgresSetting)
	cron.StartAsync()
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	cron.Clear()
	if postgresConnectionInstance != nil {
		conn, _ := postgresConnectionInstance.db.DB()
		_ = conn.Close()
	}
	return nil
}
