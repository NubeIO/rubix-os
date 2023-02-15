package main

import (
	"context"
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	log.Info("LORAWAN Plugin Enable()")
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()

	cron = gocron.NewScheduler(time.Local)
	_, _ = cron.Every("30s").Tag("SetupLorawanPlugin").Do(inst.SetupLorawanPlugin)
	cron.StartAsync()
	return nil
}

func (inst *Instance) Disable() error {
	log.Info("LORAWAN Plugin Disable()")
	inst.enabled = false
	cron.Clear()
	if inst.cancel != nil {
		inst.cancel()
	}
	inst.running = false
	inst.fault = false
	return nil
}

func (inst *Instance) SetupLorawanPlugin() error {
	inst.lorawanDebugMsg("SetupLorawanPlugin()")
	if !inst.running {
		err := inst.GetOrMakeLorawanNetwork()
		if err == nil {
			inst.lorawanDebugMsg("Plugin Enable Success!")
			inst.running = true
			inst.fault = false
			inst.csConnected = false
			err = inst.connectToCS()
			if err != nil {
				return err
			}
			inst.ctx, inst.cancel = context.WithCancel(context.Background())
			if csrest.IsCSConnectionError(err) {
				go inst.connectToCSLoop(inst.ctx)
			}
			go inst.syncChirpstackDevicesLoop(inst.ctx)
		} else {
			inst.fault = true
			inst.lorawanErrorMsg("Couldn't start lorawan plugin, problem getting/creating lorawan network")
		}
	}
	return nil
}

func (inst *Instance) GetOrMakeLorawanNetwork() error {
	inst.lorawanDebugMsg("GetOrMakeLorawanNetwork()")
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if err != nil {
		q, err = inst.createNetwork()
		if err != nil {
			inst.lorawanErrorMsg("lorawan: Cannot create network: ", err)
			return err
		}
		inst.lorawanDebugMsg("lorawan: Created default network")
		err = nil
	}
	if q != nil {
		inst.networkUUID = q.UUID
		return nil
	} else {
		return errors.New("lorawan: Error creating default network")
	}
	return errors.New("123")
}
