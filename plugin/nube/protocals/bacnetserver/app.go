package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/writemode"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body.NetworkInterface == "" {
		interfaces, err := nets.GetInterfacesNames()
		if err != nil {
			return nil, err
		}
		for _, name := range interfaces.Names {
			if name != "lo" {
				iface, _ := nets.GetNetworkByIface(name)
				if iface.IP != "" {
					body.NetworkInterface = name
				}
			}
		}
		if body.NetworkInterface == "" {
			return nil, errors.New("network interface can not be empty try, eth0")
		}
	}
	body.NumberOfNetworksPermitted = integer.New(1)
	nets, err := inst.db.GetNetworksByPluginName(body.PluginPath, api.Args{})
	if err != nil {
		return nil, err
	}
	for _, net := range nets {
		if net != nil {
			errMsg := fmt.Sprintf("bacnet-server: only max one network is allowed with bacnet")
			inst.bacnetErrorMsg(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	body.Port = integer.New(defaultPort)
	inst.bacnetDebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	err = inst.bacnetStoreNetwork(network)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("err:%s", err.Error()))
	}
	device := &model.Device{
		Name:        network.Name,
		NetworkUUID: network.UUID,
		CommonEnable: model.CommonEnable{
			Enable: boolean.NewTrue(),
		},
	}
	device, err = inst.addDevice(device)
	if err != nil {
		return nil, err
	}

	if boolean.IsFalse(network.Enable) {
		err = inst.networkUpdateErr(network, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError)
		err = inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError, true)
	}

	return network, nil
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.bacnetDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	network, err := inst.db.GetNetwork(body.NetworkUUID, api.Args{WithDevices: true})
	if err != nil {
		return nil, err
	}
	if len(network.Devices) >= 1 {
		errMsg := fmt.Sprintf("only max one device is allowed")
		inst.bacnetErrorMsg(errMsg)
		return nil, errors.New(errMsg)
	}

	body.NumberOfDevicesPermitted = integer.New(1)
	body.CommonIP.Host = inst.getIp(network.NetworkInterface)
	if integer.IsNil(body.DeviceObjectId) {
		body.DeviceObjectId = integer.New(2508)
	}
	body.Port = 47808
	inst.bacnetDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}

	err = inst.bacnetStoreDevice(device)
	if err != nil {
		return nil, errors.New("issue on add bacnet-device to store")
	}

	if boolean.IsFalse(device.Enable) {
		err = inst.deviceUpdateErr(device, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
		inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	}

	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.bacnetDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	if body.ObjectType == "" {
		errMsg := fmt.Sprintf("point object type can not be empty")
		inst.bacnetErrorMsg(errMsg)
		return nil, errors.New(errMsg)
	}
	inst.bacnetDebugMsg("addPoint(): ", body.Name)
	point, err = inst.db.CreatePoint(body, true, true)
	if point == nil || err != nil {
		inst.bacnetDebugMsg("addPoint(): failed to create bacnet point: ", body.Name)
		return nil, errors.New("failed to create bacnet point")
	}
	inst.bacnetDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))
	if boolean.IsFalse(point.Enable) {
		err = inst.pointUpdateErr(point, "point disabled", model.MessageLevel.Warning, model.CommonFaultCode.PointError)
	}
	return point, nil
}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.bacnetDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("updateNetwork():  nil network object")
		return
	}
	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "network disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	} else {
		body.CommonFault.InFault = false
		body.CommonFault.MessageLevel = model.MessageLevel.Info
		body.CommonFault.MessageCode = model.CommonFaultCode.Ok
		body.CommonFault.Message = ""
		body.CommonFault.LastOk = time.Now().UTC()
	}

	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil || network == nil {
		return nil, err
	}

	if boolean.IsFalse(network.Enable) {
		// DO POLLING DISABLE ACTIONS
		inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError, true)
	} else {
		inst.db.ClearErrorsForAllDevicesOnNetwork(network.UUID, true)
	}

	network, err = inst.db.UpdateNetwork(body.UUID, network, true)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.bacnetDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("updateDevice(): nil device object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Warning
		body.CommonFault.MessageCode = model.CommonFaultCode.DeviceError
		body.CommonFault.Message = "device disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	} else {
		body.CommonFault.InFault = false
		body.CommonFault.MessageLevel = model.MessageLevel.Info
		body.CommonFault.MessageCode = model.CommonFaultCode.Ok
		body.CommonFault.Message = ""
		body.CommonFault.LastOk = time.Now().UTC()
	}

	device, err = inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil {
		return nil, err
	}

	err = inst.bacnetStoreDevice(device)
	if err != nil {
		return nil, err
	}

	if boolean.IsFalse(device.Enable) {
		inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)

	} else {
		// DO POLLING ENABLE ACTIONS FOR DEVICE
		err = inst.db.ClearErrorsForAllPointsOnDevice(device.UUID)
		if err != nil {
			inst.bacnetDebugMsg("updateDevice(): error on ClearErrorsForAllPointsOnDevice(): ", err)
		}
	}

	device, err = inst.db.UpdateDevice(device.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) updatePoint(body *model.Point) (*model.Point, error) {
	inst.bacnetDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("updatePoint(): nil point object")
		return nil, errors.New("nil point object")
	}

	if !isWriteable(body.WriteMode, body.ObjectType) { // clear writeable point properties if point is not writeable
		body = writemode.ResetWriteableProperties(body)
	}

	inst.bacnetDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.bacnetDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "point disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	}

	point, err := inst.db.UpdatePoint(body.UUID, body, true, false)
	if err != nil {
		inst.bacnetDebugMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}
	err = inst.updatePointName(body)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.

	point = nil
	inst.bacnetDebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.bacnetDebugMsg("writePoint(): nil point object")
		return
	}

	inst.bacnetDebugMsg(fmt.Sprintf("writePoint() body: %+v", body))
	inst.bacnetDebugMsg(fmt.Sprintf("writePoint() priority: %+v", body.Priority))

	/*
		point, err = inst.db.GetPoint(pntUUID, api.Args{})
		if err != nil || point == nil {
			inst.modbusErrorMsg("writePoint(): bad response from GetPoint(), ", err)
			return nil, err
		}

		if !isWriteable(point.WriteMode, point.ObjectType) { // if point isn't writeable then reset writeable properties and do `UpdatePoint()`
			point = resetWriteableProperties(point)
			point, err = inst.db.UpdatePoint(pntUUID, point, true, false)
			if err != nil || point == nil {
				inst.modbusDebugMsg("writePoint(): bad response from UpdatePoint() err:", err)
				return nil, err
			}
			return point, nil
	*/

	point, _, isWriteValueChange, _, err := inst.db.PointWrite(pntUUID, body, false)
	if err != nil || point == nil {
		inst.bacnetDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("writePoint(): bad response from GetDevice()")
		return nil, err
	}

	if boolean.IsTrue(point.Enable) {
		if isWriteValueChange {
			if writemode.IsWriteable(point.WriteMode) {
				point.WritePollRequired = boolean.NewTrue()
				if point.WriteMode != model.WriteAlways && point.WriteMode != model.WriteOnce {
					point.ReadPollRequired = boolean.NewTrue()
				} else {
					point.ReadPollRequired = boolean.NewFalse()
				}
			} else {
				point.WritePollRequired = boolean.NewFalse()
			}
			point, err = inst.db.UpdatePoint(point.UUID, point, true, true)
			if err != nil || point == nil {
				inst.bacnetDebugMsg("writePoint(): bad response from UpdatePoint() err:", err)
				inst.pointUpdateErr(point, fmt.Sprint("writePoint(): cannot find PollingPoint for point: ", point.UUID), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
				return point, err
			}
			return point, nil
		}
	}
	return point, nil
}

func (inst *Instance) updatePointName(body *model.Point) error {
	device, err := inst.db.GetDevice(body.DeviceUUID, api.Args{})
	if err != nil {
		return err
	}
	return inst.writeBacnetPointName(body, body.Name, device.NetworkUUID, device.UUID) // update the bacnet point name
}

// initPointsNames on start update all the point names
func (inst *Instance) initPointsNames() error {
	net, err := inst.db.GetNetwork(inst.networkUUID, api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		inst.bacnetErrorMsg(fmt.Sprintf("write-all-point-names: network-UUID%s  err:%s", inst.networkUUID, err.Error()))
		return err
	}
	for _, dev := range net.Devices {
		for _, point := range dev.Points {
			err := inst.writeBacnetPointName(point, point.Name, dev.NetworkUUID, dev.UUID)
			if err != nil {
				inst.bacnetErrorMsg(fmt.Sprintf("write-all-point-name: point-name:%s  err:%s", point.Name, err.Error()))
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.bacnetDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deleteNetwork(): nil network object")
		return
	}
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	ok, err = inst.closeBacnetStoreNetwork(body.UUID)
	return ok, err
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.bacnetDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deleteDevice(): nil device object")
		return
	}

	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.bacnetDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deletePoint(): nil point object")
		return
	}

	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) pointUpdate(point *model.Point, value float64, readSuccess, clearFaults bool) (*model.Point, error) {
	if readSuccess {
		point.OriginalValue = float.New(value)
	}
	point, err := inst.db.UpdatePoint(point.UUID, point, true, clearFaults)
	if err != nil {
		inst.bacnetDebugMsg("UpdatePoint() error: ", err)
		return nil, err
	}
	return point, nil
}

// THIS SHOULD NOT BE USED, CHANGE TO pointUpdate() (Above)
func (inst *Instance) pointWrite(uuid string, value float64) error {
	inst.bacnetErrorMsg("pointWrite() DON'T USE THIS FUNCTION: ")
	priority := map[string]*float64{"_16": &value}
	pointWriter := model.PointWriter{Priority: &priority}
	_, _, _, _, err := inst.db.PointWrite(uuid, &pointWriter, true)
	if err != nil {
		inst.bacnetErrorMsg("bacnet-server: pointWrite()", err)
	}
	return err
}

func (inst *Instance) pointUpdateSuccess(uuid string) error {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	point.InSync = boolean.NewTrue()
	err := inst.db.UpdatePointErrors(uuid, &point)
	if err != nil {
		inst.bacnetErrorMsg("bacnet-server: pointUpdateSuccess()", err)
	}
	return err
}

func (inst *Instance) pointUpdateErr(point *model.Point, message string, messageLevel string, messageCode string) error {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = messageLevel
	point.CommonFault.MessageCode = messageCode
	point.CommonFault.Message = message
	point.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdatePointErrors(point.UUID, point)
	if err != nil {
		inst.bacnetErrorMsg(" pointUpdateErr()", err)
	}
	return err
}

func (inst *Instance) deviceUpdateErr(device *model.Device, message string, messageLevel string, messageCode string) error {
	device.CommonFault.InFault = true
	device.CommonFault.MessageLevel = messageLevel
	device.CommonFault.MessageCode = messageCode
	device.CommonFault.Message = message
	device.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdateDeviceErrors(device.UUID, device)
	if err != nil {
		inst.bacnetErrorMsg(" deviceUpdateErr()", err)
	}
	return err
}

func (inst *Instance) networkUpdateErr(network *model.Network, message string, messageLevel string, messageCode string) error {
	network.CommonFault.InFault = true
	network.CommonFault.MessageLevel = messageLevel
	network.CommonFault.MessageCode = messageCode
	network.CommonFault.Message = message
	network.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdateNetworkErrors(network.UUID, network)
	if err != nil {
		inst.bacnetErrorMsg(" networkUpdateErr()", err)
	}
	return err
}

func (inst *Instance) getNetworks() ([]*model.Network, error) {
	return inst.db.GetNetworks(api.Args{})
}
