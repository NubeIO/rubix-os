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

	cron = gocron.NewScheduler(time.UTC)
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("SyncTMVPointNamesToLorawan").Do(inst.updatePointNames)
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("CreateModbusNetworkDevicesPoints").Do(inst.createModbusNetworkDevicesAndPoints)
	cron.StartAsync()
	token, err := inst.GetChirpstackToken()
	if err != nil {
		inst.tmvDebugMsg("GetChirpstackToken() err: ", err)
		return nil
	}
	inst.tmvDebugMsg("GetChirpstackToken() token: ", token)
	profileUUID, err := inst.GetChirpstackDeviceProfileUUID(token.Token)
	inst.tmvDebugMsg("GetChirpstackDeviceProfileUUID() profileUUID: ", profileUUID)
	// response, _ := inst.AddChirpstackDevice(1, 666, "Test21", "4E75626549103FFF", "7ba56021-ac01-4957-b6a0-d320df36f5f0")
	// fmt.Println(response)
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
