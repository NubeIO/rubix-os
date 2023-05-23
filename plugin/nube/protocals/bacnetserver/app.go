package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/NubeIO/rubix-os/utils/integer"
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
	network, err = inst.db.CreateNetwork(body)
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
		inst.bacnetErrorMsg(err)
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
	body.Host = inst.getIp(network.NetworkInterface)
	if body.Host == "" {
		body.Host = "192.168.15.100"
	}
	if body.Host == "0.0.0.0" {
		body.Host = "192.168.15.100"
	}
	if integer.IsNil(body.DeviceObjectId) {
		body.DeviceObjectId = integer.New(2508)
	}
	if body.Port == 0 {
		body.Port = 47808
	}
	inst.bacnetDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}

	err = inst.bacnetStoreDevice(device)
	if err != nil {
		inst.bacnetErrorMsg("issue on add bacnet-device to store")
		inst.bacnetErrorMsg(err)
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
		// return nil, errors.New(errMsg) // Use if we should error out and not create the point
		body.ObjectType = "analog_value" // Use if we should default to analog_value
	}
	_, _, isBool := setObjectType(body.ObjectType)
	body.IsTypeBool = boolean.New(isBool)
	body.PointPriorityArrayMode = model.PriorityArrayToWriteValue
	body.WritePriority = integer.New(16)
	inst.bacnetDebugMsg("addPoint(): ", body.Name)
	point, err = inst.db.CreatePoint(body, true)
	if point == nil || err != nil {
		inst.bacnetDebugMsg("addPoint(): failed to create bacnet point: ", body.Name)
		return nil, errors.New("failed to create bacnet point")
	}
	inst.bacnetDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))
	err = inst.updatePointName(point)
	if err != nil {
		inst.bacnetErrorMsg("addPoint(): failed to set bacnet point name: ", point.Name)
	}
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

	network, err = inst.db.UpdateNetwork(body.UUID, body)
	if err != nil || network == nil {
		return nil, err
	}

	if boolean.IsFalse(network.Enable) {
		// DO POLLING DISABLE ACTIONS
		inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError, true)
	} else {
		inst.db.ClearErrorsForAllDevicesOnNetwork(network.UUID, true)
	}
	err = inst.bacnetStoreNetwork(network)
	if err != nil {
		inst.bacnetDebugMsg("updateNetwork(): bacnetStoreNetwork: ", network.UUID)
	}

	network, err = inst.db.UpdateNetwork(body.UUID, network)
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

	device, err = inst.db.UpdateDevice(body.UUID, body)
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

	device, err = inst.db.UpdateDevice(device.UUID, body)
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

	if body.ObjectType == "" {
		errMsg := fmt.Sprintf("point object type can not be empty")
		inst.bacnetErrorMsg(errMsg)
		// return nil, errors.New(errMsg) // Use if we should error out and not create the point
		body.ObjectType = "analog_value" // Use if we should default to analog_value
	}
	_, _, isBool := setObjectType(body.ObjectType)
	body.IsTypeBool = boolean.New(isBool)
	body.PointPriorityArrayMode = model.PriorityArrayToWriteValue

	inst.bacnetDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.bacnetDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "point disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	}

	point, err := inst.db.UpdatePoint(body.UUID, body)
	if err != nil {
		inst.bacnetDebugMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}
	err = inst.updatePointName(point)
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

	point, _, isWriteValueChange, _, err := inst.db.PointWrite(pntUUID, body)
	if err != nil || point == nil {
		inst.bacnetDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	if isWriteValueChange {
		point.WritePollRequired = boolean.NewTrue()
	}

	if boolean.IsTrue(point.Enable) {
		dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
		if err != nil {
			inst.bacnetErrorMsg("BACnetServerPolling(): Device not found")
			return point, nil
		}

		_, err = inst.SyncFFPointWithBACnetServerPoint(point, point.DeviceUUID, dev.NetworkUUID, isWriteValueChange)
		if err != nil {
			inst.bacnetErrorMsg(err)
			if isWriteValueChange {
				point.WritePollRequired = boolean.NewTrue()
				point, err = inst.db.UpdatePoint(point.UUID, point)
			}
		}
		return point, nil
	}
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	point, err = inst.db.UpdatePoint(point.UUID, point)
	if err != nil || point == nil {
		inst.bacnetDebugMsg("writePoint(): bad response from UpdatePoint() err:", err)
		inst.pointUpdateErr(point, fmt.Sprint("writePoint(): bad response from UpdatePoint() err:", err), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
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

// THE FOLLOWING FUNCTIONS ARE CALLED FROM WITHIN THE PLUGIN
func (inst *Instance) pointUpdate(point *model.Point, value float64, readWriteSuccess bool) (*model.Point, error) {
	if readWriteSuccess {
		point.OriginalValue = float.New(value)
		point.WritePollRequired = boolean.NewFalse()
	}
	point, err := inst.db.UpdatePoint(point.UUID, point)
	if err != nil {
		inst.bacnetDebugMsg("UpdatePoint() error: ", err)
		return nil, err
	}
	return point, nil
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
