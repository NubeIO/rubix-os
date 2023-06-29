package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/services/pollqueue"
	"github.com/NubeIO/rubix-os/utils/array"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/NubeIO/rubix-os/utils/writemode"
	"go.bug.st/serial"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	inst.modbusDebugMsg("addNetwork(): ", body.Name)

	// indicates that ui should display polling statistics
	body.HasPollingStatistics = true

	network, err = inst.db.CreateNetwork(body)
	if err != nil {
		inst.modbusErrorMsg("addNetwork(): failed to create modbus network: ", body.Name)
		return nil, errors.New("failed to create modbus network")
	}

	if boolean.IsTrue(network.Enable) {
		conf := inst.GetConfig().(*Config)
		pollQueueConfig := pollqueue.Config{EnablePolling: conf.EnablePolling, LogLevel: conf.LogLevel}
		pollManager := NewPollManager(&pollQueueConfig, &inst.db, network.UUID, network.Name, inst.pluginUUID, inst.pluginName, float.NonNil(network.MaxPollRate))
		pollManager.StartPolling()
		inst.NetworkPollManagers = append(inst.NetworkPollManagers, pollManager)
	} else {
		err = inst.networkUpdateErr(network, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError)
		err = inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError, true)
	}
	return network, nil
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.modbusDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.modbusDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.modbusDebugMsg("addDevice(): failed to create modbus device: ", body.Name)
		return nil, errors.New("failed to create modbus device")
	}

	inst.modbusDebugMsg("addDevice(): ", body.UUID)

	if boolean.IsFalse(device.Enable) {
		err = inst.deviceUpdateErr(device, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
		err = inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	}

	// NOTHING TO DO ON DEVICE CREATED
	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.modbusDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.modbusDebugMsg("addPoint(): ", body.Name)

	if isWriteable(body.WriteMode, body.ObjectType) {
		body.WritePollRequired = boolean.NewTrue()
		body.EnableWriteable = boolean.NewTrue()
	} else {
		body = writemode.ResetWriteableProperties(body)
	}
	body.ReadPollRequired = boolean.NewTrue()

	isTypeBool := checkForBooleanType(body.ObjectType, body.DataType)
	body.IsTypeBool = nils.NewBool(isTypeBool)

	isOutput := checkForOutputType(body.ObjectType)
	body.IsOutput = nils.NewBool(isOutput)

	point, err = inst.db.CreatePoint(body, true)
	if point == nil || err != nil {
		inst.modbusDebugMsg("addPoint(): failed to create modbus point: ", body.Name)
		return nil, errors.New(fmt.Sprint("failed to create modbus point. err: ", err))
	}
	inst.modbusDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	dev, err := inst.db.GetDevice(point.DeviceUUID, args.Args{})
	if err != nil || dev == nil {
		inst.modbusDebugMsg("addPoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		return
	}

	if boolean.IsTrue(point.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
		// DO POLLING ENABLE ACTIONS FOR POINT
		pp := pollqueue.NewPollingPoint(point.UUID, point.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true, pollqueue.NORMAL_RETRY, false, false, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
		// netPollMan.PollQueue.AddPollingPoint(pp)
		// netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	} else {
		err = inst.pointUpdateErr(point, "point disabled", model.MessageLevel.Warning, model.CommonFaultCode.PointError)
	}
	return point, nil

}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.modbusDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("updateNetwork():  nil network object")
		return
	}

	// indicates that ui should display polling statistics
	body.HasPollingStatistics = true

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Warning
		body.CommonFault.MessageCode = model.CommonFaultCode.NetworkError
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

	restartPolling := false
	if body.MaxPollRate != network.MaxPollRate {
		restartPolling = true
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(network.UUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("updateNetwork(): cannot find NetworkPollManager for network: ", network.UUID)
		return
	}

	if netPollMan.NetworkName != network.Name {
		netPollMan.NetworkName = network.Name
	}

	if boolean.IsFalse(network.Enable) && netPollMan.Enable == true {
		// DO POLLING DISABLE ACTIONS
		netPollMan.StopPolling()
		inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError, true)
	} else if restartPolling || (boolean.IsTrue(network.Enable) && netPollMan.Enable == false) {
		if restartPolling {
			netPollMan.StopPolling()
		}
		// DO POLLING Enable ACTIONS
		netPollMan.StartPolling()
		inst.db.ClearErrorsForAllDevicesOnNetwork(network.UUID, true)
	}

	network, err = inst.db.UpdateNetwork(body.UUID, network)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.modbusDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("updateDevice(): nil device object")
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
	if err != nil || device == nil {
		return nil, err
	}

	if boolean.IsTrue(device.Enable) { // If Enabled we need to GetDevice so we get Points
		device, err = inst.db.GetDevice(device.UUID, args.Args{WithPoints: true})
		if err != nil || device == nil {
			return nil, err
		}
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(device.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("updateDevice(): cannot find NetworkPollManager for network: ", device.NetworkUUID)
		return
	}
	if boolean.IsFalse(device.Enable) && netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(device.UUID) {
		// DO POLLING DISABLE ACTIONS FOR DEVICE
		inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(device.UUID)

	} else if boolean.IsTrue(device.Enable) && !netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(device.UUID) {
		// DO POLLING ENABLE ACTIONS FOR DEVICE
		err = inst.db.ClearErrorsForAllPointsOnDevice(device.UUID)
		if err != nil {
			inst.modbusDebugMsg("updateDevice(): error on ClearErrorsForAllPointsOnDevice(): ", err)
		}
		for _, pnt := range device.Points {
			if boolean.IsTrue(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, device.NetworkUUID, netPollMan.FFPluginUUID)
				netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true, pollqueue.NORMAL_RETRY, false, false, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
				// netPollMan.PollQueue.AddPollingPoint(pp)  //This is the original tested way, above is new so that on device update, it will re-poll write-once points
			}
		}

	} else if boolean.IsTrue(device.Enable) {
		// TODO: Currently on every device update, all device points are removed, and re-added.
		device.CommonFault.InFault = false
		device.CommonFault.MessageLevel = model.MessageLevel.Info
		device.CommonFault.MessageCode = model.CommonFaultCode.Ok
		device.CommonFault.Message = ""
		device.CommonFault.LastOk = time.Now().UTC()
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(device.UUID)
		for _, pnt := range device.Points {
			if boolean.IsTrue(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, device.NetworkUUID, netPollMan.FFPluginUUID)
				netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true, pollqueue.NORMAL_RETRY, false, false, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
				// netPollMan.PollQueue.AddPollingPoint(pp)  //This is the original tested way, above is new so that on device update, it will re-poll write-once points
			}
		}
	}
	// TODO: NEED TO ACCOUNT FOR OTHER CHANGES ON DEVICE.  It would be useful to have a way to know if the device polling rates were changed.

	device, err = inst.db.UpdateDevice(device.UUID, device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.modbusDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("updatePoint(): nil point object")
		return
	}

	/*
		pnt, err := inst.db.GetPoint(body.UUID, api.Args{WithPriority: true})
		if pnt == nil || err != nil {
			inst.modbusErrorMsg("could not find pointID: ", pp.FFPointUUID)
			netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
			continue
		}
	*/

	if isWriteable(body.WriteMode, body.ObjectType) {
		body.WritePollRequired = boolean.NewTrue()
		body.EnableWriteable = boolean.NewTrue()
	} else {
		body = writemode.ResetWriteableProperties(body)
	}

	isTypeBool := checkForBooleanType(body.ObjectType, body.DataType)
	body.IsTypeBool = nils.NewBool(isTypeBool)

	inst.modbusDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.modbusDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "point disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	}
	body.CommonFault.InFault = false
	body.CommonFault.MessageLevel = model.MessageLevel.Info
	body.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	body.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	body.CommonFault.LastOk = time.Now().UTC()
	point, err = inst.db.UpdatePoint(body.UUID, body)
	if err != nil || point == nil {
		inst.modbusErrorMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}

	dev, err := inst.db.GetDevice(point.DeviceUUID, args.Args{})
	if err != nil || dev == nil {
		inst.modbusErrorMsg("updatePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusErrorMsg("updatePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		_ = inst.pointUpdateErr(point, "cannot find NetworkPollManager for network", model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
		return
	}

	if boolean.IsTrue(point.Enable) && boolean.IsTrue(dev.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
		// DO POLLING ENABLE ACTIONS FOR POINT
		// TODO: review these steps to check that UpdatePollingPointByUUID might work better?
		pp := pollqueue.NewPollingPoint(point.UUID, point.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true, pollqueue.NORMAL_RETRY, false, false, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
		// netPollMan.PollQueue.AddPollingPoint(pp)
		// netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	} else {
		// DO POLLING DISABLE ACTIONS FOR POINT
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
	}
	return point, nil
}

func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.

	point = nil
	inst.modbusDebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.modbusDebugMsg("writePoint(): nil point object")
		return
	}

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

	point, _, isWriteValueChange, _, err := inst.db.PointWrite(pntUUID, body)
	if err != nil {
		inst.modbusDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	dev, err := inst.db.GetDevice(point.DeviceUUID, args.Args{})
	if err != nil || dev == nil {
		inst.modbusDebugMsg("writePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("writePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		_ = inst.pointUpdateErr(point, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
		return nil, err
	}

	if boolean.IsTrue(point.Enable) {
		if isWriteValueChange || point.WriteMode == model.WriteOnceReadOnce || point.WriteMode == model.WriteOnce || (point.WriteMode == model.WriteOnceThenRead && *point.WriteValue != *point.OriginalValue) { // if the write value has changed, we need to re-add the point so that it is polled asap (if required)
			pp, _ := netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
			if pp == nil {
				if netPollMan.PollQueue.OutstandingPollingPoints.GetPollingPointIndexByPointUUID(point.UUID) > -1 {
					if writemode.IsWriteable(point.WriteMode) {
						netPollMan.PollQueue.PointsUpdatedWhilePolling[point.UUID] = true // this triggers a write post at ASAP priority (for writeable points).
						point.WritePollRequired = boolean.NewTrue()
						if point.WriteMode != model.WriteAlways && point.WriteMode != model.WriteOnce {
							point.ReadPollRequired = boolean.NewTrue()
						} else {
							point.ReadPollRequired = boolean.NewFalse()
						}
					} else {
						netPollMan.PollQueue.PointsUpdatedWhilePolling[point.UUID] = false
						point.WritePollRequired = boolean.NewFalse()
					}
					point.CommonFault.InFault = false
					point.CommonFault.MessageLevel = model.MessageLevel.Info
					point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
					point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
					point.CommonFault.LastOk = time.Now().UTC()
					point, err = inst.db.UpdatePoint(point.UUID, point)
					if err != nil || point == nil {
						inst.modbusDebugMsg("writePoint(): bad response from UpdatePoint() err:", err)
						inst.pointUpdateErr(point, fmt.Sprint("writePoint(): cannot find PollingPoint for point: ", point.UUID), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
						return point, err
					}
					return point, nil
				} else {
					inst.modbusDebugMsg("writePoint(): cannot find PollingPoint for point (could be out for polling: ", point.UUID)
					_ = inst.pointUpdateErr(point, "writePoint(): cannot find PollingPoint for point: ", model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
					return point, err
				}
			}
			if writemode.IsWriteable(point.WriteMode) {
				point.WritePollRequired = boolean.NewTrue()
			} else {
				point.WritePollRequired = boolean.NewFalse()
			}
			if point.WriteMode != model.WriteAlways && point.WriteMode != model.WriteOnce {
				point.ReadPollRequired = boolean.NewTrue()
			} else {
				point.ReadPollRequired = boolean.NewFalse()
			}
			point.CommonFault.InFault = false
			point.CommonFault.MessageLevel = model.MessageLevel.Info
			point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
			point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
			point.CommonFault.LastOk = time.Now().UTC()
			point, err = inst.db.UpdatePoint(point.UUID, point)
			if err != nil || point == nil {
				inst.modbusDebugMsg("writePoint(): bad response from UpdatePoint() err:", err)
				inst.pointUpdateErr(point, fmt.Sprint("writePoint(): bad response from UpdatePoint() err:", err), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
				return point, err
			}

			// pp.PollPriority = model.PRIORITY_ASAP   // TODO: THIS NEEDS TO BE IMPLEMENTED SO THAT ONLY MANUAL WRITES ARE PROMOTED TO ASAP PRIORITY
			netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, false, pollqueue.IMMEDIATE_RETRY, false, false, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
			// netPollMan.PollQueue.AddPollingPoint(pp)
			// netPollMan.PollQueue.UpdatePollingPointByPointUUID(point.UUID, model.PRIORITY_ASAP)

			/*
				netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
				//DO POLLING ENABLE ACTIONS FOR POINT
				//TODO: review these steps to check that UpdatePollingPointByUUID might work better?
				pp := pollqueue.NewPollingPoint(body.UUID, body.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
				netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
				//netPollMan.PollQueue.AddPollingPoint(pp)
				//netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
			*/
		}
	} else {
		// DO POLLING DISABLE ACTIONS FOR POINT
		netPollMan.PollQueue.RemovePollingPointByPointUUID(pntUUID)
	}
	return point, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.modbusDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("deleteNetwork(): nil network object")
		return
	}
	found := false
	for index, netPollMan := range inst.NetworkPollManagers {
		if netPollMan.FFNetworkUUID == body.UUID {
			netPollMan.StopPolling()
			// Next remove the NetworkPollManager from the slice in polling instance
			inst.NetworkPollManagers[index] = inst.NetworkPollManagers[len(inst.NetworkPollManagers)-1]
			inst.NetworkPollManagers = inst.NetworkPollManagers[:len(inst.NetworkPollManagers)-1]
			found = true
		}
	}
	if !found {
		inst.modbusDebugMsg("deleteNetwork(): cannot find NetworkPollManager for network: ", body.UUID)
	}
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.modbusDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("deleteDevice(): nil device object")
		return
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(body.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("deleteDevice(): cannot find NetworkPollManager for network: ", body.NetworkUUID)
		_ = inst.deviceUpdateErr(body, "cannot find NetworkPollManager for network", model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
		return
	}
	netPollMan.PollQueue.RemovePollingPointByDeviceUUID(body.UUID)
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.modbusDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.modbusDebugMsg("deletePoint(): nil point object")
		return
	}

	dev, err := inst.db.GetDevice(body.DeviceUUID, args.Args{})
	if err != nil || dev == nil {
		inst.modbusDebugMsg("addPoint(): bad response from GetDevice()")
		return false, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.modbusDebugMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		_ = inst.pointUpdateErr(body, "cannot find NetworkPollManager for network", model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
		return
	}

	netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
	otherPointsOnSameDeviceExist := netPollMan.PollQueue.CheckPollingQueueForDevUUID(body.DeviceUUID)
	if !otherPointsOnSameDeviceExist {
		netPollMan.PollQueue.RemoveDeviceFromActiveDevicesList(body.DeviceUUID)
	}
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// THE FOLLOWING FUNCTIONS ARE CALLED FROM WITHIN THE PLUGIN
func (inst *Instance) pointUpdate(point *model.Point, value float64, readSuccess bool) (*model.Point, error) {
	if readSuccess {
		point.OriginalValue = float.New(value)
	}
	_, err := inst.db.UpdatePoint(point.UUID, point)
	if err != nil {
		inst.modbusDebugMsg("MODBUS UPDATE POINT UpdatePointPresentValue() error: ", err)
		return nil, err
	}
	return point, nil
}

func (inst *Instance) pointUpdateErr(point *model.Point, message string, messageLevel string, messageCode string) error {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = messageLevel
	point.CommonFault.MessageCode = messageCode
	point.CommonFault.Message = fmt.Sprintf("modbus: %s", message)
	point.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdatePointErrors(point.UUID, point)
	if err != nil {
		inst.modbusErrorMsg(" pointUpdateErr()", err)
	}
	return err
}

func (inst *Instance) deviceUpdateErr(device *model.Device, message string, messageLevel string, messageCode string) error {
	device.CommonFault.InFault = true
	device.CommonFault.MessageLevel = messageLevel
	device.CommonFault.MessageCode = messageCode
	device.CommonFault.Message = fmt.Sprintf("modbus: %s", message)
	device.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdateDeviceErrors(device.UUID, device)
	if err != nil {
		inst.modbusErrorMsg(" deviceUpdateErr()", err)
	}
	return err
}

func (inst *Instance) networkUpdateErr(network *model.Network, message string, messageLevel string, messageCode string) error {
	network.CommonFault.InFault = true
	network.CommonFault.MessageLevel = messageLevel
	network.CommonFault.MessageCode = messageCode
	network.CommonFault.Message = fmt.Sprintf("modbus: %s", message)
	network.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdateNetworkErrors(network.UUID, network)
	if err != nil {
		inst.modbusErrorMsg(" networkUpdateErr()", err)
	}
	return err
}

func (inst *Instance) listSerialPorts() (*array.Array, error) {
	ports, err := serial.GetPortsList()
	p := array.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}
