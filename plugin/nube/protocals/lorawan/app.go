package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/plugin/nube/protocals/lorawan/csrest"
	"github.com/NubeIO/rubix-os/schema/lorawanschema"
	"github.com/NubeIO/rubix-os/schema/schema"
	"github.com/NubeIO/rubix-os/utils/boolean"
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
			time.Sleep(time.Duration(inst.config.SyncPeriodMins * float32(time.Minute)))
		}
	}
}

func (inst *Instance) syncChirpstackDevices() {
	log.Info("lorawan: Syncing Devices")
	devices, err := inst.chirpStack.GetDevices()
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

func (inst *Instance) syncAddMissingDevices(csDevices []*csrest.DevicesResult) {
	for _, csDev := range csDevices {
		currDev, _ := inst.db.GetDeviceByArgs(api.Args{AddressUUID: &csDev.DevEUI})
		if currDev == nil {
			_, err := inst.addMissingDeviceResult(csDev)
			if err != nil {
				continue
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (inst *Instance) addMissingDeviceResult(csDev *csrest.DevicesResult) (*model.Device, error) {
	key, err := inst.chirpStack.DeviceOTAAKeyGet(csDev.DevEUI)
	if err != nil {
		log.Error("lorawan: error getting cs device keys on addMissingDevice: ", err)
		return nil, err
	}
	ffDev := csDeviceResultToFlowDeviceWithKey(csDev, inst.networkUUID, &key)
	err = inst.createDevice(&ffDev)
	if err != nil {
		log.Error("lorawan: error adding device on addMissingDevice: ", err)
	}
	return nil, err
}

func (inst *Instance) addMissingDeviceSingle(csDev *csrest.DeviceSingle) (*model.Device, error) {
	key, err := inst.chirpStack.DeviceOTAAKeyGet(csDev.Device.DevEUI)
	if err != nil {
		log.Error("lorawan: error getting cs device keys on addMissingDevice")
		return nil, err
	}
	ffDev := csDeviceSingleToFlowDeviceWithKey(csDev, inst.networkUUID, &key)
	err = inst.createDevice(&ffDev)
	if err != nil {
		log.Error("lorawan: error adding device on addMissingDevice")
	}
	return nil, err
}

func (inst *Instance) syncRemoveOldDevices(csDevices []*csrest.DevicesResult) {
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

func (inst *Instance) syncUpdateDevices(csDevices []*csrest.DevicesResult) {
	for _, csDev := range csDevices {
		currDev, err := inst.db.GetDeviceByArgs(api.Args{AddressUUID: &csDev.DevEUI})
		if err != nil || currDev == nil {
			log.Error("lorawan: get device: ", err)
			continue
		}
		csDevSingle, err := inst.chirpStack.GetDevice(csDev.DevEUI)
		if err != nil {
			log.Error("lorawan: Chirpstack get device: ", err)
			continue
		}
		key, err := inst.chirpStack.DeviceOTAAKeyGet(csDev.DevEUI)
		if err != nil {
			log.Error("lorawan: Chirpstack get device key: ", err)
			continue
		}
		ffDevNew := csDeviceSingleToFlowDeviceWithKey(csDevSingle, inst.networkUUID, &key)
		ffDevNew.UUID = currDev.UUID
		_, err = inst.db.UpdateDevice(ffDevNew.UUID, &ffDevNew)
		if err != nil {
			log.Error("lorawan: update device during sync: ", err)
		} else {
			log.Tracef("lorawan: updated device during sync: EUI=%s UUID=%s", csDev.DevEUI, currDev.UUID)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (inst *Instance) connectToCS() error {
	inst.chirpStack = csrest.InitRest(inst.config.CSAddress, inst.config.CSPort, inst.basePath)
	inst.chirpStack.SetDeviceLimit(inst.config.DeviceLimit)
	token := inst.getCSToken()
	inst.chirpStack.SetToken(token)
	err := inst.chirpStack.ConnectTest()
	if err == nil {
		inst.setCSConnected()
	} else {
		log.Error("lorawan: failed to connect to Chirpstack: ", err)
	}
	return err
}

func (inst *Instance) connectToCSLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Debug("lorawan: stopping Chirpstack connection loop")
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

func (inst *Instance) getCSToken() string {
	if inst.config.CSToken == "" {
		data, err := os.ReadFile(inst.config.CSTokenFilePath)
		if err != nil {
			log.Error("lorwan: no Chirpstack token provided in either file or config: ", err)
		}
		str := string(data)
		str = strings.TrimSpace(str)
		return str
	}
	return inst.config.CSToken
}

func (inst *Instance) createNetwork() (*model.Network, error) {
	maxNetworks := new(int)
	*maxNetworks = maxAllowedNetworks
	net := model.Network{
		Name:                      pluginName,
		CommonDescription:         model.CommonDescription{Description: "Chirpstack"},
		CommonEnable:              model.CommonEnable{Enable: boolean.NewTrue()},
		PluginPath:                pluginPath,
		NumberOfNetworksPermitted: maxNetworks,
		TransportType:             transportType,
		IP:                        inst.config.CSAddress,
		Port:                      &inst.config.CSPort,
	}
	return inst.addNetwork(&net)
}

func (inst *Instance) deleteNetwork() {
	inst.db.DeleteNetwork(inst.networkUUID)
	inst.Disable()
}

func flowDeviceToCsDevice(ffDev *model.Device) csrest.DeviceSingle {
	csDev := csrest.DeviceBody{}
	csDev.Name = ffDev.Name
	csDev.Description = ffDev.Description
	csDev.DevEUI = *ffDev.AddressUUID
	csDev.IsDisabled = !*ffDev.Enable
	csDev.ApplicationID = fmt.Sprintf("%d", ffDev.AddressId)
	csDev.DeviceProfileID = ffDev.Model
	csDev.SkipFCntCheck = *ffDev.ZeroMode
	return csrest.DeviceSingle{Device: &csDev}
}

func csDeviceResultToFlowDeviceWithKey(csDev *csrest.DevicesResult, netUUID string, key *string) model.Device {
	ffDev := model.Device{}
	keyStr := ""
	if key != nil {
		keyStr = *key
	}
	csConvertDevice(&ffDev,
		netUUID,
		csDev.Name,
		csDev.Description,
		csDev.DevEUI,
		true,
		csDev.ApplicationID,
		csDev.DeviceProfileID,
		keyStr,
		nil)
	return ffDev
}

func csDeviceResultToFlowDevice(csDev *csrest.DevicesResult, netUUID string) model.Device {
	return csDeviceResultToFlowDeviceWithKey(csDev, netUUID, nil)
}

func csDeviceSingleToFlowDeviceWithKey(csDev *csrest.DeviceSingle, netUUID string, key *string) model.Device {
	ffDev := model.Device{}
	keyStr := ""
	if key != nil {
		keyStr = *key
	}
	csConvertDevice(&ffDev,
		netUUID,
		csDev.Device.Name,
		csDev.Device.Description,
		csDev.Device.DevEUI,
		csDev.Device.IsDisabled,
		csDev.Device.ApplicationID,
		csDev.Device.DeviceProfileID,
		keyStr,
		&csDev.Device.SkipFCntCheck)
	return ffDev
}

func csDeviceSingleToFlowDevice(csDev *csrest.DeviceSingle, netUUID string) model.Device {
	return csDeviceSingleToFlowDeviceWithKey(csDev, netUUID, nil)
}

func csConvertDevice(ffDev *model.Device, netUUID string, name string, description string, eui string, disabled bool, appID string, devProfID string, key string, skipFcntCheck *bool) {
	ffDev.NetworkUUID = netUUID
	ffDev.Name = name
	ffDev.Description = description
	ffDev.AddressUUID = &eui
	disabled = !disabled
	ffDev.Enable = &disabled
	addrId, _ := strconv.Atoi(appID)
	ffDev.AddressId = addrId
	ffDev.Model = devProfID
	ffDev.Manufacture = key
	ffDev.ZeroMode = skipFcntCheck
}

func (inst *Instance) addDevice(device *model.Device) error {
	csDev := flowDeviceToCsDevice(device)
	err := inst.chirpStack.AddDevice(&csDev)
	if err != nil {
		return err
	}
	err = inst.chirpStack.DeviceOTAAKeyAdd(csDev.Device.DevEUI, device.Manufacture)
	if err != nil {
		inst.chirpStack.DeleteDevice(csDev.Device.DevEUI)
		return err
	}
	err = inst.createDevice(device)
	return err
}

func (inst *Instance) updateDevice(device *model.Device) error {
	csDev := flowDeviceToCsDevice(device)
	err := inst.chirpStack.UpdateDevice(&csDev)
	if err != nil {
		return err
	}
	err = inst.chirpStack.DeviceOTAAKeyUpdate(*device.AddressUUID, device.Manufacture)
	if err != nil {
		return err
	}
	_, err = inst.db.UpdateDevice(device.UUID, device)
	return err
}

func (inst *Instance) deleteDevice(device *model.Device) error {
	err := inst.chirpStack.DeleteDevice(*device.AddressUUID)
	if err != nil {
		return err
	}
	_, err = inst.db.DeleteDevice(device.UUID)
	return err
}

func (inst *Instance) fillDeviceProfilesSchema(devSchema *lorawanschema.DeviceSchema) {
	devProfiles, err := inst.chirpStack.GetDeviceProfiles()
	if err != nil {
		return
	}
	applications, err := inst.chirpStack.GetApplications()
	if err != nil {
		return
	}
	for i := range devProfiles.Result {
		option := schema.OptionOneOf{
			Const: devProfiles.Result[i].Id,
			Title: devProfiles.Result[i].Name,
		}
		devSchema.DeviceProfileID.Options = append(devSchema.DeviceProfileID.Options, option)
	}
	devSchema.DeviceProfileID.Default = devProfiles.Result[0].Id
	for i := range applications.Result {
		num, _ := strconv.Atoi(applications.Result[i].Id)
		option := schema.OptionOneOfInt{
			Const: num,
			Title: applications.Result[i].Name,
		}
		devSchema.ApplicationID.Options = append(devSchema.ApplicationID.Options, option)
	}
	num, _ := strconv.Atoi(applications.Result[0].Id)
	devSchema.ApplicationID.Default = num
	if len(applications.Result) == 1 {
		devSchema.ApplicationID.ReadOnly = true
	}
}
