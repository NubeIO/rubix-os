package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.tmvDebugMsg("TMV Plugin Enable()")
	inst.enabled = true
	inst.pluginName = name
	inst.setUUID()
	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if nets != nil {
		inst.networks = nets
	} else if err != nil {
		inst.networks = nil
	}

	cron = gocron.NewScheduler(time.Local)
	if inst.config.Job.EnableConfigSteps {
		_, _ = cron.Every(inst.config.Job.Frequency).Tag("CreateAndActivateLoRaWANDevices").Do(inst.createAndActivateChirpstackDevices)
		_, _ = cron.Every(inst.config.Job.Frequency).Tag("SyncTMVPointNamesToLorawan").Do(inst.updatePointNames)
		_, _ = cron.Every(inst.config.Job.Frequency).Tag("CreateModbusNetworkDevicesPoints").Do(inst.createModbusNetworkDevicesAndPoints)
	}
	_, _ = cron.Every(1).Day().At("02:00").Tag("UpdateIOModuleRTC").Do(inst.updateIOModuleRTC)
	// _, _ = cron.Every(inst.config.Job.Frequency).Tag("UpdateIOModuleRTC").Do(inst.updateIOModuleRTC)
	cron.StartAsync()
	_, next := cron.NextRun()
	inst.tmvDebugMsg("Next CRON job @ ", next.String())
	return nil
}

func (inst *Instance) Disable() error {
	inst.tmvDebugMsg("TMV Plugin Disable()")
	inst.enabled = false
	cron.Clear()
	if inst.pollingEnabled && inst.pollingCancel != nil {
		inst.pollingEnabled = false
		inst.pollingCancel()
		inst.pollingCancel = nil
	}
	return nil
}
