package main

import (
	"context"
	"errors"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if err != nil {
		q, err = inst.createNetwork()
		if err != nil {
			log.Error("lorawan: Cannot create network: ", err)
			return err
		}
		log.Info("lorawan: Created default network")
		err = nil
	}
	if q != nil {
		inst.networkUUID = q.UUID
	} else {
		return errors.New("lorawan: Error creating default network")
	}

	inst.csConnected = false
	inst.REST = csrest.InitRest(inst.config.CSAddress, inst.config.CSPort, inst.config.CSToken)
	inst.REST.SetDeviceLimit(inst.config.DeviceLimit)
	err = inst.connectToCS()
	inst.ctx, inst.cancel = context.WithCancel(context.Background())
	if csrest.IsCSConnectionError(err) {
		go inst.connectToCSLoop(inst.ctx)
	}
	go inst.syncChirpstackDevicesLoop(inst.ctx)
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	inst.cancel()
	return nil
}
