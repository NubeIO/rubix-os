package main

import (
	"context"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) syncChirpstackDevicesLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug("lorawan: Stopping CS connection loop")
			return
		default:
			if inst.csConnected {
				inst.syncChirpstackDevices()
			}
			time.Sleep(time.Duration(inst.config.SyncPeriodMins) * time.Minute)
		}
	}
}

func (inst *Instance) syncChirpstackDevices() {
	log.Info("lorawan: Syncing Devices")
	devices, err := inst.REST.GetDevices()
	if err != nil {
		if csrest.IsCSConnectionError(err) {
			inst.setCSDisconnected(err)
		}
		log.Warn("lorawan: Failed to sync devices due to CS connection error")
		return
	}
	if len(devices.Result) == 0 {
		log.Warn("lorawan: No devices in CS to sync")
		return
	}
	inst.syncAddMissingDevices(devices.Result)
	inst.syncRemoveOldDevices(devices.Result)
	inst.syncUpdateDevices(devices.Result)
}

func (inst *Instance) syncAddMissingDevices(csDevices []csmodel.Device) {
	for _, csDev := range csDevices {
		currDev, _ := inst.db.GetDeviceByArgs(api.Args{AddressUUID: &csDev.DevEUI})
		if currDev == nil {
			inst.createDeviceFromCSDevice(&csDev)
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func (inst *Instance) syncRemoveOldDevices(csDevices []csmodel.Device) {
	currNetwork, err := inst.db.GetNetwork(inst.networkUUID, api.Args{WithDevices: true})
	if err != nil || currNetwork == nil {
		return
	}
	currDevices := currNetwork.Devices
	for _, currDev := range currDevices {
		found := false
		for _, csDev := range csDevices {
			if csDev.DevEUI == *currDev.AddressUUID {
				found = true
				break
			}
		}
		if found {
			continue
		}
		log.Warn("lorawan: Removing old device. EUI=", *currDev.AddressUUID)
		inst.db.DeleteDevice(currDev.UUID)
		time.Sleep(20 * time.Millisecond)
	}
}

func (inst *Instance) syncUpdateDevices(csDevices []csmodel.Device) {
	for _, csDev := range csDevices {
		currDev, err := inst.db.GetDeviceByArgs(api.Args{AddressUUID: &csDev.DevEUI})
		if err != nil || currDev == nil {
			log.Error("lorawan: GetDeviceByArgs() err: ", err)
			continue
		}
		if currDev.Name != csDev.Name &&
			currDev.Description != csDev.Description {
			currDev.Name = csDev.Name
			currDev.Description = csDev.Description
			_, err = inst.db.UpdateDevice(currDev.UUID, currDev)
			if err != nil {
				log.Error("lorawan: Error updating device during sync: ", err)
			} else {
				log.Debugf("lorawan: Updated device during sync: EUI=%s UUID=%s", csDev.DevEUI, currDev.UUID)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func (inst *Instance) connectToCS() error {
	port := inst.config.CSPort
	addr := inst.config.CSAddress
	token := inst.config.CSToken
	// port = 8080
	// addr = "192.168.15.15"
	// token = ""
	inst.REST = csrest.InitRest(addr, port)
	inst.REST.SetDeviceLimit(inst.config.DeviceLimit)
	if token == "" {
		err := inst.REST.Login(inst.config.CSUsername, inst.config.CSPassword)
		if err != nil {
			return err
		}
	} else {
		inst.REST.SetToken(token)
	}
	err := inst.REST.ConnectTest()
	if err == nil {
		inst.setCSConnected()
	} else if !csrest.IsCSConnectionError(err) {
		log.Error("lorawan: Failed to connect to Chirpstack and unable to retry. Error: ", err)
	}
	return err
}

func (inst *Instance) connectToCSLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Debug("lorawan: Stopping main loop")
			return ctx.Err()
		default:
			time.Sleep(10 * time.Second)
			err := inst.connectToCS()
			if !csrest.IsCSConnectionError(err) {
				return err
			}
		}
	}
}

func (inst *Instance) setCSConnected() {
	inst.csConnected = true
	log.Info("lorawan: Connected to Chirpstack")
	net := model.Network{
		CommonUUID: model.CommonUUID{
			UUID: inst.networkUUID,
		},
		CommonFault: model.CommonFault{
			InFault: false,
			Message: "",
		},
	}
	inst.db.UpdateNetworkErrors(inst.networkUUID, &net)
}

func (inst *Instance) setCSDisconnected(err error) {
	inst.csConnected = false
	log.Warn("lorawan: Lost connection to Chirpstack. Error: ", err)
	net := model.Network{
		CommonFault: model.CommonFault{
			InFault: true,
			Message: err.Error(),
		},
	}
	inst.db.UpdateNetwork(inst.networkUUID, &net)
	go inst.connectToCSLoop(inst.ctx)
}

func (inst *Instance) createNetwork() (*model.Network, error) {
	maxNetworks := new(int)
	*maxNetworks = maxAllowedNetworks
	enable := true
	net := model.Network{
		CommonNameUnique:          model.CommonNameUnique{Name: pluginName},
		CommonDescription:         model.CommonDescription{Description: "Chirpstack"},
		CommonEnable:              model.CommonEnable{Enable: &enable},
		PluginPath:                pluginPath,
		NumberOfNetworksPermitted: maxNetworks,
		TransportType:             transportType,
		IP:                        inst.config.CSAddress,
		Port:                      &inst.config.CSPort,
	}
	return inst.addNetwork(&net)
}
