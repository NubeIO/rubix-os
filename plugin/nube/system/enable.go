package main

import (
	"fmt"
	"github.com/NubeIO/rubix-os/api"
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
	var frequency = "60s"
	if inst.config.Schedule.Frequency != "" {
		frequency = inst.config.Schedule.Frequency
	}
	_, _ = cron.Every(frequency).Tag("ScheduleCheck").Do(inst.runSchedule)
	cron.StartAsync()

	// TODO: this is added just to update the EnableWriteable property on legacy points.  It should be removed in future versions
	networks, err := inst.db.GetNetworks(api.Args{WithDevices: true, WithPoints: true, WithPriority: true})
	if err != nil {
		log.Error("SYSTEM Enable() GetNetworks error:", err)
		return nil
	}
	fmt.Println("SYSTEM PLUGIN ENABLE: TEMPORARY CODE TO SET EnableWriteable PROPERTY ON EVERY POINT.  REMOVE IN FUTURE RELEASES")
	for _, network := range networks {
		for _, device := range network.Devices {
			for _, point := range device.Points {
				fmt.Println("SYSTEM PLUGIN ENABLE point: ", point.Name)
				inst.db.UpdatePointPlugin(point.UUID, point)
			}
		}
	}
	return nil
}

func (inst *Instance) Disable() error {
	log.Info("SYSTEM Plugin Disable()")
	inst.enabled = false
	cron.Clear()
	return nil
}
