package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/maplora/legacylorarest"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *Instance) mapLoRa(forceAddToExistingFFNetwork bool) error {
	inst.maploraDebugMsg("mapLoRa()")

	// CREATE FF LORA NETWORK FROM LEGACY LORA NETWORK
	legacyLoRaNet, err := inst.GetLegacyLoRaNetwork()
	if err != nil {
		inst.maploraErrorMsg(err)
	}
	if legacyLoRaNet == nil {
		inst.maploraErrorMsg("no legacy lora network found")
		return errors.New("no legacy lora network found")
	}

	nets, err := inst.db.GetNetworksByPluginName("lora", api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}
	networkExistsForPort := false
	var existingFFNet *model.Network
	for _, ffNet := range nets {
		if ffNet.SerialPort != nil && *ffNet.SerialPort == legacyLoRaNet.Port {
			networkExistsForPort = true
			existingFFNet = ffNet
			break
		}
	}

	ffNetwork := existingFFNet
	if networkExistsForPort {
		inst.maploraErrorMsg("existing ff lora network exists for port:", legacyLoRaNet.Port)
		if !forceAddToExistingFFNetwork {
			return errors.New(fmt.Sprint("existing ff lora network exists for port:", legacyLoRaNet.Port))
		}
	} else {
		ffNetwork, err = inst.addLoRaNetworkFromLegacyNetwork(legacyLoRaNet)
		if err != nil || ffNetwork == nil {
			return err
		}
	}

	// CREATE FF LORA DEVICES FROM LEGACY LORA DEVICES
	legacyLoRaDevArray, err := inst.GetLegacyLoRaDevices()
	if err != nil {
		inst.maploraErrorMsg(err)
	}
	if legacyLoRaDevArray == nil || len(*legacyLoRaDevArray) <= 0 {
		inst.maploraErrorMsg("no legacy lora devices found")
		return errors.New("no legacy lora devices found")
	}

	for _, legacyDev := range *legacyLoRaDevArray {
		ffDevice, err := inst.addLoRaDeviceFromLegacyDevice(ffNetwork, &legacyDev)
		if err != nil || ffDevice == nil {
			inst.maploraErrorMsg("couldn't add ff lora device: ", err)
		}
	}
	return nil
}

func (inst *Instance) addLoRaNetworkFromLegacyNetwork(legacyNet *legacylorarest.LoRaNet) (network *model.Network, err error) {
	if legacyNet != nil {
		inst.maploraDebugMsg("addLoRaNetworkFromLegacyNetwork(): ", legacyNet.Port)
	} else {
		return nil, errors.New("failed to create lora network")
	}
	network = &model.Network{}
	network.PluginPath = "lora"

	network.Enable = boolean.NewTrue()

	port := legacyNet.Port
	network.SerialPort = &port
	network.Name = port

	baud := uint(legacyNet.BaudRate)
	network.SerialBaudRate = &baud

	stopBits := uint(legacyNet.StopBits)
	network.SerialStopBits = &stopBits

	parity := legacyNet.Parity
	network.SerialParity = &parity

	byteSize := uint(legacyNet.ByteSize)
	network.SerialDataBits = &byteSize

	timeout := int(legacyNet.Timeout)
	network.SerialTimeout = &timeout

	network, err = inst.db.CreateNetwork(network, false)
	if err != nil {
		inst.maploraErrorMsg("addLoRaNetworkFromLegacyNetwork(): failed to create tmv network")
		return nil, errors.New("failed to create lora network")
	}

	return network, nil
}

func (inst *Instance) addLoRaDeviceFromLegacyDevice(ffNet *model.Network, legacyDev *legacylorarest.LoRaDev) (device *model.Device, err error) {
	if legacyDev != nil {
		inst.maploraDebugMsg("addLoRaDeviceFromLegacyDevice(): ", legacyDev.Name)
	} else {
		return nil, errors.New("failed to create lora device")
	}

	device = &model.Device{}
	device.Name = legacyDev.Name
	device.Description = legacyDev.Description
	device.Enable = boolean.NewTrue()

	device.NetworkUUID = ffNet.UUID

	id := legacyDev.ID
	device.AddressUUID = &id

	ffModelName, err := inst.ConvertLegacyLoRaModelToFFModel(legacyDev.DeviceModel)
	if err != nil {
		return nil, errors.New("failed to create lora device")
	}
	device.Model = ffModelName

	cli := client.NewLocalClient()
	device, err = cli.CreateDevicePlugin(device, "lora")
	if err != nil {
		inst.maploraErrorMsg("addLoRaDeviceFromLegacyDevice(): failed to create lora device")
		return nil, errors.New("failed to create lora device")
	}

	return device, nil
}
