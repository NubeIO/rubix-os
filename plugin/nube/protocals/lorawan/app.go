package main

import (
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) syncChirpstackDevicesLoop() {
	for {
		if inst.csConnected {
			inst.syncChirpstackDevices()
		} else {
			log.Warn("lorawan: Failed to sync devices due to CS connection error")
		}
		time.Sleep(time.Duration(inst.config.SyncPeriodMins) * time.Minute)
	}
}

func (inst *Instance) syncChirpstackDevices() {
	log.Info("lorawan: Syncing Devices")
	devices, err := inst.REST.GetDevices()
	if csrest.IsCSConnectionError(err) {
		inst.setCSDisconnected(err)
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
		}
	}
}

func (inst *Instance) syncRemoveOldDevices(csDevices []csmodel.Device) {
	currNetwork, _ := inst.db.GetNetwork(inst.networkUUID, api.Args{WithDevices: true})
	currDevices := currNetwork.Devices
	for _, currDev := range currDevices {
		found := false
		for _, csDev := range csDevices {
			if csDev.DevEUI == *currDev.CommonDevice.AddressUUID {
				found = true
				break
			}
		}
		if found {
			continue
		}
		log.Warn("lorawan: Removing old device. EUI=", *currDev.CommonDevice.AddressUUID)
		inst.db.DeleteDevice(currDev.UUID)
	}
}

func (inst *Instance) syncUpdateDevices(csDevices []csmodel.Device) {
	for _, csDev := range csDevices {
		currDev, _ := inst.db.GetDeviceByArgs(api.Args{AddressUUID: &csDev.DevEUI})
		if currDev.CommonName.Name != csDev.Name &&
			currDev.CommonDescription.Description != csDev.Description {
			currDev.CommonName.Name = csDev.Name
			currDev.CommonDescription.Description = csDev.Description
			_, err := inst.db.UpdateDevice(currDev.UUID, currDev, true)
			if err != nil {
				log.Error("lorawan: Error updating device during sync: ", err)
			} else {
				log.Debugf("lorawan: Updated device during sync: EUI=%s UUID=%s", csDev.DevEUI, currDev.UUID)
			}
		}
	}
}

func (inst *Instance) connectToCS() error {
	rest, err := csrest.CSLogin(inst.config.CSAddress, inst.config.CSPort,
		inst.config.CSUsername, inst.config.CSPassword)
	if err == nil {
		log.Info("lorawan: Connected to Chirpstack")
		inst.REST = rest
		inst.csConnected = true
	} else if !csrest.IsCSConnectionError(err) {
		log.Error("lorawan: Failed to connect to Chirpstack and unable to retry. Error: ", err)
	}
	return err
}

func (inst *Instance) connectToCSLoop() error {
	for {
		time.Sleep(5 * time.Second)
		err := inst.connectToCS()
		if !csrest.IsCSConnectionError(err) {
			// TODO: disable self
			return err
		}
	}
}

func (inst *Instance) setCSDisconnected(err error) {
	inst.csConnected = false
	log.Warn("lorawan: Lost connection to Chirpstack. Cause: ", err)
	go inst.connectToCSLoop()
}

func (inst *Instance) createNetwork() (*model.Network, error) {
	maxNetworks := new(int)
	*maxNetworks = maxAllowedNetworks
	net := model.Network{
		CommonNameUnique:          model.CommonNameUnique{Name: pluginName},
		CommonDescription:         model.CommonDescription{Description: "Chirpstack"},
		PluginPath:                pluginPath,
		NumberOfNetworksPermitted: maxNetworks,
		TransportType:             transportType,
		IP:                        inst.config.CSAddress,
		Port:                      &inst.config.CSPort,
	}
	return inst.addNetwork(&net)
}
