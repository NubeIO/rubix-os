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
}

func (inst *Instance) syncAddMissingDevices(csDevices []csmodel.Device) {
	for _, csDev := range csDevices {
		currDev, _ := inst.db.GetOneDeviceByArgs(api.Args{AddressUUID: &csDev.DevEUI})
		if currDev == nil {
			newDev := model.Device{}
			inst.csConvertDevice(&newDev, &csDev)
			_, err := inst.db.CreateDevice(&newDev)
			if err != nil {
				log.Error("lorawan: Error adding new device during sync: ", err)
				continue
			}
			inst.euiAdd(csDev.DevEUI)
			log.Info("lorawan: Added device ", csDev.DevEUI)
		}
	}
}

func (inst *Instance) syncRemoveOldDevices(csDevices []csmodel.Device) {
	currDevices, _ := inst.db.GetDevices(api.Args{})
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
		inst.euiRemove(*currDev.CommonDevice.AddressUUID)
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

func (inst *Instance) euiIndex(eui string) int {
	for i, v := range inst.deviceEUIs {
		if v == eui {
			return i
		}
	}
	return -1
}

func (inst *Instance) euiExists(eui string) bool {
	return inst.euiIndex(eui) > -1
}

func (inst *Instance) euiAdd(eui string) {
	inst.deviceEUIs = append(inst.deviceEUIs, eui)
}

func (inst *Instance) euiRemove(eui string) {
	arr := inst.deviceEUIs
	i := inst.euiIndex(eui)
	if i == -1 {
		return
	}
	arr[i] = arr[len(arr)-1]
	arr[len(arr)-1] = ""
	inst.deviceEUIs = arr[:len(arr)-1]
}
