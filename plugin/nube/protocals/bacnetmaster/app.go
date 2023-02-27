package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/services/pollqueue"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
	"github.com/NubeIO/flow-framework/utils/writemode"
	address "github.com/NubeIO/lib-networking/ip"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body.NetworkInterface == "" {
		return nil, errors.New("network interface can not be empty try, eth0")
	}

	inst.bacnetDebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body, true)
	if network == nil || err != nil {
		inst.bacnetErrorMsg("addNetwork(): failed to create bacnet network: ", body.Name)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("failed to create bacnet network")
	}
	err = inst.makeBacnetStoreNetwork(network)
	if err != nil {
		inst.bacnetErrorMsg("addNetwork(): issue on add bacnet-network to store err ", err.Error())
		// fmt.Sprintf("issue on add bacnet-device to store err:%s", err.Error())
	}
	body.MaxPollRate = float.New(0.1)
	body.TransportType = "ip"
	if boolean.IsTrue(network.Enable) {
		conf := inst.GetConfig().(*Config)
		pollQueueConfig := pollqueue.Config{EnablePolling: conf.EnablePolling, LogLevel: conf.LogLevel}
		pollManager := NewPollManager(&pollQueueConfig, &inst.db, network.UUID, inst.pluginUUID, inst.pluginName, float.NonNil(network.MaxPollRate))
		pollManager.StartPolling()
		inst.NetworkPollManagers = append(inst.NetworkPollManagers, pollManager)
	} else {
		err = inst.networkUpdateErr(network, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError)
		err = inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError, true)
	}
	return network, err
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.bacnetDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.bacnetDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.bacnetDebugMsg("addDevice(): failed to create bacnet device: ", body.Name)
		return nil, errors.New("failed to create bacnet device")
	}
	if body.Host == "" {
		body.Host = "192.168.15.100"
	}
	if body.Host == "0.0.0.0" {
		body.Host = "192.168.15.100"
	}
	if body.Port == 0 {
		body.Port = 47808
	}
	if float.IsNil(body.FastPollRate) {
		body.FastPollRate = float.New(1)
	}
	if float.IsNil(body.NormalPollRate) {
		body.NormalPollRate = float.New(15)
	}
	if float.IsNil(body.SlowPollRate) {
		body.SlowPollRate = float.New(120)
	}
	err = address.New().IsIPAddrErr(body.Host)
	if body == nil {
		inst.bacnetDebugMsg("addDevice(): nil device object")
		return nil, errors.New(fmt.Sprintf("invalid ip addr %s", body.Host))
	}
	err = inst.bacnetStoreDevice(device)
	if err != nil {
		inst.bacnetDebugMsg(fmt.Sprintf("bacnet-device: add device to store err: %s", err.Error()))
		return nil, errors.New(fmt.Sprintf("bacnet-device: add device to store err: %s", err.Error()))
	}
	inst.bacnetDebugMsg("addDevice(): ", body.UUID)

	if boolean.IsFalse(device.Enable) {
		err = inst.deviceUpdateErr(device, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
		inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	}

	// NOTHING TO DO ON DEVICE CREATED
	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.bacnetDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.bacnetDebugMsg("addPoint(): ", body.Name)

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

	point, err = inst.db.CreatePoint(body, true, true)
	if point == nil || err != nil {
		inst.bacnetDebugMsg("addPoint(): failed to create bacnet point: ", body.Name)
		return nil, err
	}
	inst.bacnetDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("addPoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		return nil, err
	}

	if boolean.IsTrue(point.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
		// DO POLLING ENABLE ACTIONS FOR POINT
		pp := pollqueue.NewPollingPoint(point.UUID, point.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true, pollqueue.NORMAL_RETRY, false, false) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
		// netPollMan.PollQueue.AddPollingPoint(pp)
		// netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	} else {
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

	restartPolling := false
	if body.MaxPollRate != network.MaxPollRate {
		restartPolling = true
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(network.UUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("updateNetwork(): cannot find NetworkPollManager for network: ", network.UUID)
		return
	}
	err = inst.makeBacnetStoreNetwork(network)
	if err != nil {
		inst.bacnetDebugMsg("updateNetwork(): makeBacnetStoreNetwork: , err: ", network.UUID, err)
	}

	if boolean.IsFalse(network.Enable) && netPollMan.Enable {
		// DO POLLING DISABLE ACTIONS
		netPollMan.StopPolling()
		inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError, true)
	} else if restartPolling || (boolean.IsTrue(network.Enable) == true && netPollMan.Enable == false) {
		if restartPolling {
			netPollMan.StopPolling()
		}
		// DO POLLING Enable ACTIONS
		netPollMan.StartPolling()
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
		inst.bacnetDebugMsg(fmt.Sprintf("bacnet-device: update device to store err: %s", err.Error()))
		return nil, err
	}

	if boolean.IsTrue(device.Enable) { // If Enabled we need to GetDevice so we get Points
		device, err = inst.db.GetDevice(device.UUID, api.Args{WithPoints: true})
		if err != nil || device == nil {
			return nil, err
		}
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(device.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("updateDevice(): cannot find NetworkPollManager for network: ", device.NetworkUUID)
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
			inst.bacnetDebugMsg("updateDevice(): error on ClearErrorsForAllPointsOnDevice(): ", err)
		}
		for _, pnt := range device.Points {
			if boolean.IsTrue(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, device.NetworkUUID, netPollMan.FFPluginUUID)
				netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true, pollqueue.NORMAL_RETRY, false, false) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
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
				netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true, pollqueue.NORMAL_RETRY, false, false) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
				// netPollMan.PollQueue.AddPollingPoint(pp)  //This is the original tested way, above is new so that on device update, it will re-poll write-once points
			}
		}
	}
	// TODO: NEED TO ACCOUNT FOR OTHER CHANGES ON DEVICE.  It would be useful to have a way to know if the device polling rates were changed.

	device, err = inst.db.UpdateDevice(device.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.bacnetDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("updatePoint(): nil point object")
		return
	}

	/*
		pnt, err := inst.db.GetPoint(body.UUID, api.Args{WithPriority: true})
		if pnt == nil || err != nil {
			inst.bacnetErrorMsg("could not find pointID: ", pp.FFPointUUID)
			netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
			continue
		}

	*/

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
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	point, err = inst.db.UpdatePoint(body.UUID, body)
	if err != nil || point == nil {
		inst.bacnetDebugMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}
	// err = inst.updatePointName(body)  //TODO: Does this need to be added (from BACnet Server)
	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("updatePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("updatePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		_ = inst.pointUpdateErr(point, "cannot find NetworkPollManager for network", model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
		return
	}

	if boolean.IsTrue(point.Enable) && boolean.IsTrue(dev.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
		// DO POLLING ENABLE ACTIONS FOR POINT
		// TODO: review these steps to check that UpdatePollingPointByUUID might work better?
		pp := pollqueue.NewPollingPoint(point.UUID, point.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, true, pollqueue.NORMAL_RETRY, false, false) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
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

	point, _, isWriteValueChange, _, err := inst.db.PointWrite(pntUUID, body)
	if err != nil || point == nil {
		inst.bacnetDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("writePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if err != nil {
		inst.bacnetDebugMsg("writePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		_ = inst.pointUpdateErr(point, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
		return
	}

	if boolean.IsTrue(point.Enable) {
		if isWriteValueChange { // if the write value has changed, we need to re-add the point so that it is polled asap (if required)
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
						inst.bacnetDebugMsg("writePoint(): bad response from UpdatePoint() err:", err)
						inst.pointUpdateErr(point, fmt.Sprint("writePoint(): cannot find PollingPoint for point: ", point.UUID), model.MessageLevel.Fail, model.CommonFaultCode.SystemError)
						return point, err
					}
					return point, nil
				} else {
					inst.bacnetDebugMsg("writePoint(): cannot find PollingPoint for point (could be out for polling: ", point.UUID)
					_ = inst.pointUpdateErr(point, "writePoint(): cannot find PollingPoint for point: ", model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
					return point, err
				}
			}
			pp.PollPriority = model.PRIORITY_ASAP
			netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true, false, pollqueue.NORMAL_RETRY, false, false) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
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
	inst.bacnetDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deleteNetwork(): nil network object")
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
		inst.bacnetDebugMsg("deleteNetwork(): cannot find NetworkPollManager for network: ", body.UUID)
	}
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	ok, err = inst.closeBacnetStoreNetwork(body.UUID)
	return ok, nil
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.bacnetDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deleteDevice(): nil device object")
		return
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(body.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("deleteDevice(): cannot find NetworkPollManager for network: ", body.NetworkUUID)
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
	inst.bacnetDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.bacnetDebugMsg("deletePoint(): nil point object")
		return
	}

	dev, err := inst.db.GetDevice(body.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.bacnetDebugMsg("addPoint(): bad response from GetDevice()")
		return false, err
	}

	netPollMan, err := inst.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		inst.bacnetDebugMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
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

func (inst *Instance) pointUpdate(point *model.Point, value *float64, readSuccess bool) (*model.Point, error) {
	if readSuccess {
		point.OriginalValue = value
	}
	point, err := inst.db.UpdatePoint(point.UUID, point)
	if err != nil {
		inst.bacnetDebugMsg("UpdatePoint() error: ", err)
		return nil, err
	}
	return point, nil
}

func (inst *Instance) pointUpdateFromPriorityArray(point *model.Point, priorityArray map[string]*float64, presentValue *float64) (*model.Point, error) {
	point, err := priorityarray.ApplyMapToPriorityArray(point, &priorityArray)
	point.OriginalValue = presentValue
	point, err = inst.db.UpdatePoint(point.UUID, point)
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
