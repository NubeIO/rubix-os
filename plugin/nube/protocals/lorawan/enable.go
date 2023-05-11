package main

import (
	"context"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) Enable() error {
	log.Info("lorawan: plugin enable")
	inst.enabled = true
	inst.setUUID()
	inst.ctx, inst.cancel = context.WithCancel(context.Background())
	err := inst.setupLorawanPlugin()
	if err != nil {
		return err
	}
	inst.busServ()
	return nil
}

func (inst *Instance) Disable() error {
	log.Info("lorawan: plugin disable")
	inst.enabled = false
	inst.csConnected = false
	if inst.cancel != nil {
		inst.cancel()
	}
	inst.busDisable()
	return nil
}

func (inst *Instance) setupLorawanPlugin() error {
	err := inst.getOrMakeLorawanNetwork()
	if err != nil {
		return err
	}
	inst.csConnected = false
	err = inst.connectToCS()
	if csrest.IsCSConnectionError(err) {
		go inst.connectToCSLoop(inst.ctx)
	}
	go inst.syncChirpstackDevicesLoop(inst.ctx)
	return nil
}

func (inst *Instance) getOrMakeLorawanNetwork() error {
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if err != nil {
		q, err = inst.createNetwork()
		if err != nil {
			log.Error("lorawan: cannot create network: ", err)
			return err
		}
		log.Debug("lorawan: created default network")
	}
	inst.networkUUID = q.UUID
	return nil
}
