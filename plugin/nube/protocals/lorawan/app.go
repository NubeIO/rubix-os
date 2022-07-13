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
		}
		time.Sleep(time.Duration(inst.config.SyncPeriodMins) * time.Minute)
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
	currNetwork, err := inst.db.GetNetwork(inst.networkUUID, api.Args{WithDevices: true})
	if err != nil {
		return
	}
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
	err := inst.REST.Connect()
	if err == nil {
		inst.setCSConnected()
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
			return err
		}
	}
}

func (inst *Instance) setCSConnected() {
	inst.csConnected = true
	log.Info("lorawan: Connected to Chirpstack")
	net := model.Network{
		CommonFault: model.CommonFault{
			InFault: false,
			Message: "",
		},
	}
	inst.db.UpdateNetwork(inst.networkUUID, &net, true)
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
	inst.db.UpdateNetwork(inst.networkUUID, &net, true)
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
