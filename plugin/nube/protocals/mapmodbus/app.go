package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/plugin/nube/protocals/mapmodbus/legacymodbusrest"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/boolean"
)

func (inst *Instance) mapModbus(forceAddToExistingFFNetwork bool) error {
	inst.mapmodbusDebugMsg("mapModbus()")

	// CREATE FF LORA NETWORK FROM LEGACY LORA NETWORK
	legacyModbusNets, err := inst.GetLegacyModbusNetworksAndDevices()
	if err != nil {
		inst.mapmodbusErrorMsg(err)
	}
	if legacyModbusNets == nil || len(*legacyModbusNets) <= 0 {
		inst.mapmodbusErrorMsg("no legacy modbus network found")
		return errors.New("no legacy modbus network found")
	}

	nets, err := inst.db.GetNetworksByPluginName("modbus", api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}

	for _, legacyModbusNet := range *legacyModbusNets {
		networkExistsForPort := false
		var existingFFNet *model.Network
		for _, ffNet := range nets {
			if ffNet.TransportType == "ip" && legacyModbusNet.Type == "TCP" {
				if ffNet.Host != nil && *ffNet.Host == legacyModbusNet.TCPIP {
					networkExistsForPort = true
					existingFFNet = ffNet
					break
				}
			} else if ffNet.TransportType == "serial" && legacyModbusNet.Type == "RTU" {
				if ffNet.SerialPort != nil && *ffNet.SerialPort == legacyModbusNet.RTUPort {
					networkExistsForPort = true
					existingFFNet = ffNet
					break
				}
			}
		}

		ffNetwork := existingFFNet
		if networkExistsForPort {
			if legacyModbusNet.Type == "TCP" {
				inst.mapmodbusErrorMsg("existing ff modbus network exists for ", legacyModbusNet.TCPIP)
				if !forceAddToExistingFFNetwork {
					inst.mapmodbusErrorMsg(fmt.Sprint("existing ff modbus network exists for ", legacyModbusNet.TCPIP))
					continue
				}
			} else if legacyModbusNet.Type == "RTU" {
				inst.mapmodbusErrorMsg("existing ff modbus network exists for ", legacyModbusNet.RTUPort)
				if !forceAddToExistingFFNetwork {
					inst.mapmodbusErrorMsg(fmt.Sprint("existing ff modbus network exists for ", legacyModbusNet.RTUPort))
					continue
				}
			}
		} else {
			ffNetwork, err = inst.addModbusNetworkFromLegacyNetwork(&legacyModbusNet)
			if err != nil || ffNetwork == nil {
				inst.mapmodbusErrorMsg("addModbusNetworkFromLegacyNetwork() err: ", err)
				continue
			}
		}

		// CREATE FF LORA DEVICES FROM LEGACY LORA DEVICES
		if legacyModbusNet.Devices == nil || len(legacyModbusNet.Devices) <= 0 {
			inst.mapmodbusErrorMsg("mapModbus() Creating FF Devices: legacy device array is nil")
			continue
		}
		for _, legacyModbusDev := range legacyModbusNet.Devices {
			legacyModbusDevWithPoints, err := inst.GetLegacyModbusDeviceAndPoints(legacyModbusDev.UUID)
			if err != nil {
				inst.mapmodbusErrorMsg(err)
			}
			if legacyModbusDevWithPoints == nil {
				inst.mapmodbusErrorMsg("no legacy modbus devices found")
				continue
			}

			inst.mapmodbusDebugMsg("mapModbus() addModbusDeviceFromLegacyDevice() legacyDevice: ", legacyModbusDev.Name)
			ffDevice, err := inst.addModbusDeviceFromLegacyDevice(ffNetwork, &legacyModbusDev)
			if err != nil || ffDevice == nil {
				inst.mapmodbusErrorMsg("couldn't add ff modbus device: ", err)
			}

			// CREATE FF LORA POINTS FROM LEGACY LORA POINTS
			if legacyModbusDevWithPoints.Points == nil || len(legacyModbusDevWithPoints.Points) <= 0 {
				inst.mapmodbusErrorMsg("mapModbus() Creating FF Points: legacy point array is nil")
				continue
			}
			for _, legacyPnt := range legacyModbusDevWithPoints.Points {
				inst.mapmodbusDebugMsg("mapModbus() addModbusPointFromLegacyPoint() legacyDevice: ", legacyPnt.Name)
				ffPoint, err := inst.addModbusPointFromLegacyPoint(ffDevice, &legacyPnt)
				if err != nil || ffPoint == nil {
					inst.mapmodbusErrorMsg("couldn't add ff modbus point: ", err)
				}
			}
		}
	}
	return nil
}

func (inst *Instance) addModbusNetworkFromLegacyNetwork(legacyNet *legacymodbusrest.ModbusNet) (network *model.Network, err error) {
	if legacyNet != nil {
		inst.mapmodbusDebugMsg("addModbusNetworkFromLegacyNetwork(): ", legacyNet.Name)
	} else {
		return nil, errors.New("failed to create modbus network")
	}

	network = &model.Network{}
	network.PluginPath = "modbus"
	network.Enable = boolean.NewTrue()
	network.Description = legacyNet.UUID
	network.Name = legacyNet.Name

	if legacyNet.Type == "RTU" {
		port := legacyNet.RTUPort
		network.SerialPort = &port

		baud := uint(legacyNet.RTUSpeed)
		network.SerialBaudRate = &baud

		stopBits := uint(legacyNet.RTUStopBits)
		network.SerialStopBits = &stopBits

		parity := legacyNet.RTUParity
		network.SerialParity = &parity

		byteSize := uint(legacyNet.RTUByteSize)
		network.SerialDataBits = &byteSize
	}

	if legacyNet.Type == "TCP" {
		host := legacyNet.TCPIP
		network.Host = &host

		port := legacyNet.TCPPort
		network.Port = &port
	}

	timeout := int(legacyNet.Timeout)
	network.SerialTimeout = &timeout

	cli := client.NewLocalClient()
	network, err = cli.CreateNetworkPlugin(network, "modbus")
	if err != nil {
		inst.mapmodbusErrorMsg("addModbusNetworkFromLegacyNetwork(): failed to create modbus network")
		return nil, errors.New("failed to create modbus network")
	}
	return network, nil
}

func (inst *Instance) addModbusDeviceFromLegacyDevice(ffNet *model.Network, legacyDev *legacymodbusrest.ModbusDev) (device *model.Device, err error) {
	if legacyDev != nil {
		inst.mapmodbusDebugMsg("addModbusDeviceFromLegacyDevice(): ", legacyDev.Name)
	} else {
		return nil, errors.New("failed to create modbus device")
	}

	device = &model.Device{}
	device.Name = legacyDev.Name
	device.Description = legacyDev.NetworkUUID
	device.Enable = boolean.NewTrue()

	device.NetworkUUID = ffNet.UUID

	zero := legacyDev.ZeroBased
	device.ZeroMode = &zero

	id := legacyDev.Address
	device.AddressId = id

	cli := client.NewLocalClient()
	device, err = cli.CreateDevicePlugin(device, "modbus")
	if err != nil {
		inst.mapmodbusErrorMsg("addModbusDeviceFromLegacyDevice(): failed to create modbus device")
		return nil, errors.New("failed to create modbus device")
	}

	return device, nil
}

func (inst *Instance) addModbusPointFromLegacyPoint(ffDev *model.Device, legacyPnt *legacymodbusrest.ModbusPnt) (point *model.Point, err error) {
	if legacyPnt != nil {
		inst.mapmodbusDebugMsg("addModbusPointFromLegacyPoint(): ", legacyPnt.Name)
	} else {
		return nil, errors.New("failed to create modbus point")
	}

	point = &model.Point{}
	point.Name = legacyPnt.Name
	point.Description = legacyPnt.DeviceUUID
	point.Enable = boolean.NewTrue()

	point.DeviceUUID = ffDev.UUID

	register := legacyPnt.Register
	point.AddressID = &register

	fallback := legacyPnt.FallbackValue
	point.Fallback = &fallback

	objectType, dataType, writeMode, dataLength, _ := inst.ConvertLegacyModbusPropsToFFModbusPoints(legacyPnt.FunctionCode, legacyPnt.DataType)

	point.ObjectType = objectType
	point.DataType = dataType
	point.WriteMode = writeMode
	point.AddressLength = &dataLength

	cli := client.NewLocalClient()
	point, err = cli.CreatePointPlugin(point, "modbus")
	if err != nil {
		inst.mapmodbusErrorMsg("addModbusPointFromLegacyPoint(): failed to create modbus point. err: ", err)
		return nil, errors.New("failed to create modbus point")
	}

	return point, nil
}
