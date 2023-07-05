package main

import (
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.tmvDebugMsg("TMV Plugin Enable()")
	inst.enabled = true
	inst.fault = false
	inst.running = false
	inst.pluginName = name
	inst.setUUID()
	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, argspkg.Args{})
	if nets != nil {
		inst.networks = nets
	} else if err != nil {
		inst.fault = true
		inst.networks = nil
	}

	inst.startTime = time.Now().Unix()
	cron = gocron.NewScheduler(time.Local)
	if inst.config.Job.EnableConfigSteps {
		_, _ = cron.Every(inst.config.Job.Frequency).Tag("runSetupSteps").Do(inst.runSetupSteps)
	}
	if inst.config.Job.EnableCommissioning {
		_, _ = cron.Every(inst.config.Job.Frequency).Tag("checkComissioningPoints").Do(inst.checkComissioningPoints)
	}

	_, _ = cron.Every(1).Day().At("02:00").Tag("UpdateIOModuleRTC").Do(inst.updateIOModuleRTC)
	_, _ = cron.Every(1).Day().At("02:00").Tag("DisableCommissioningPoints").Do(inst.DisableCommissioningPoints)
	// _, _ = cron.Every(inst.config.Job.Frequency).Tag("UpdateIOModuleRTC").Do(inst.updateIOModuleRTC)
	cron.StartAsync()
	_, next := cron.NextRun()
	inst.tmvDebugMsg("Next CRON job @ ", next.String())
	inst.running = true
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
	inst.fault = false
	inst.running = false
	return nil
}
