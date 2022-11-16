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
	token, err := inst.GetChirpstackToken(inst.config.ChirpstackUsername, inst.config.ChirpstackPassword)
	if err != nil || token == nil && token.Token == "" {
		log.Error("createAndActivateChirpstackDevices() err: ", err)
		return nil
	}
	inst.REST = csrest.InitRest(inst.config.CSAddress, inst.config.CSPort, token.Token)
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
	if inst.cancel != nil {
		inst.cancel()
	}
	return nil
}
