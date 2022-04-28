package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	pollqueue "github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/poll-queue"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"go.bug.st/serial"
	"time"
)

//THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
//addDevice add network. Called via API call.
func (i *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body == nil {
		modbusErrorMsg("addNetwork(): nil network object")
		return nil, errors.New("empty network body, no network created")
	}
	modbusDebugMsg("addNetwork(): ", body.Name)
	network, err = i.db.CreateNetwork(body, true)
	if network == nil || err != nil {
		modbusErrorMsg("addNetwork(): failed to create modbus network: ", body.Name)
		return nil, errors.New("failed to create modbus network")
	}

	if utils.BoolIsNil(body.Enable) {
		pollManager := pollqueue.NewPollManager(&i.db, network.UUID, i.pluginUUID)
		pollManager.StartPolling()
		i.NetworkPollManagers = append(i.NetworkPollManagers, pollManager)
	}

	if err != nil {
		return nil, err
	}
	return network, nil
}

//addDevice add device. Called via API call.
func (i *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		modbusErrorMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	modbusDebugMsg("addDevice(): ", body.Name)
	device, err = i.db.CreateDevice(body)
	if device == nil || err != nil {
		modbusErrorMsg("addDevice(): failed to create modbus device: ", body.Name)
		return nil, errors.New("failed to create modbus device")
	}

	modbusDebugMsg("addDevice(): ", body.UUID)
	//NOTHING TO DO ON DEVICE CREATED
	return device, nil
}

//addPoint add point. Called via API call.
func (i *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		modbusErrorMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	modbusDebugMsg("addPoint(): ", body.Name)

	if isWriteable(body.WriteMode) {
		body.WritePollRequired = utils.NewTrue()
	} else {
		body.WritePollRequired = utils.NewFalse()
	}
	body.ReadPollRequired = utils.NewTrue()

	point, err = i.db.CreatePoint(body, true, false)
	if point == nil || err != nil {
		modbusErrorMsg("addPoint(): failed to create modbus point: ", body.Name)
		return nil, errors.New("failed to create modbus point")
	}
	modbusDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	modbusDebugMsg(fmt.Sprintf("Instance(): %+v\n", i))

	//net, err := i.db.DB.GetNetworkByDeviceUUID(point.DeviceUUID, api.Args{})
	dev, err := i.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		modbusErrorMsg("addPoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := i.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		modbusErrorMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		return
	}

	if utils.BoolIsNil(point.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(point.UUID)
		//DO POLLING ENABLE ACTIONS FOR POINT
		//TODO: review these steps to check that UpdatePollingPointByUUID might work better?
		pp := pollqueue.NewPollingPoint(point.UUID, point.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		pp.PollPriority = point.PollPriority
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
		//netPollMan.PollQueue.AddPollingPoint(pp)
		//netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	}
	return point, nil

}

//updateNetwork update network. Called via API call.
func (i *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	modbusDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		modbusErrorMsg("updateNetwork():  nil network object")
		return
	}
	netPollMan, err := i.getNetworkPollManagerByUUID(body.UUID)
	if netPollMan == nil || err != nil {
		modbusErrorMsg("updateNetwork(): cannot find NetworkPollManager for network: ", body.UUID)
		return
	}

	if utils.BoolIsNil(body.Enable) == false && netPollMan.Enable == true {
		//DO POLLING DISABLE ACTIONS
		netPollMan.StopPolling()

	} else if utils.BoolIsNil(body.Enable) == true && netPollMan.Enable == false {
		//DO POLLING Enable ACTIONS
		netPollMan.StartPolling()
	}

	network, err = i.db.UpdateNetwork(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//updateDevice update device. Called via API call.
func (i *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	modbusDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		modbusErrorMsg("updateDevice(): nil device object")
		return
	}

	netPollMan, err := i.getNetworkPollManagerByUUID(body.NetworkUUID)
	if netPollMan == nil || err != nil {
		modbusErrorMsg("updateDevice(): cannot find NetworkPollManager for network: ", body.NetworkUUID)
		return
	}

	if utils.BoolIsNil(body.Enable) == false && netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(body.UUID) {
		//DO POLLING DISABLE ACTIONS FOR DEVICE
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(body.UUID)

	} else if utils.BoolIsNil(body.Enable) == true && !netPollMan.PollQueue.CheckIfActiveDevicesListIncludes(body.UUID) {
		//DO POLLING ENABLE ACTIONS FOR DEVICE
		for _, pnt := range body.Points {
			if utils.BoolIsNil(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, body.NetworkUUID, netPollMan.FFPluginUUID)
				pp.PollPriority = pnt.PollPriority
				netPollMan.PollQueue.AddPollingPoint(pp)
			}
		}

	} else if utils.BoolIsNil(body.Enable) == true {
		//TODO: Currently on every device update, all device points are removed, and re-added.
		netPollMan.PollQueue.RemovePollingPointByDeviceUUID(body.UUID)
		for _, pnt := range body.Points {
			if utils.BoolIsNil(pnt.Enable) {
				pp := pollqueue.NewPollingPoint(pnt.UUID, pnt.DeviceUUID, body.NetworkUUID, netPollMan.FFPluginUUID)
				pp.PollPriority = pnt.PollPriority
				netPollMan.PollQueue.AddPollingPoint(pp)
			}
		}
	}
	// TODO: NEED TO ACCOUNT FOR OTHER CHANGES ON DEVICE.  It would be useful to have a way to know if the device polling rates were changed.

	device, err = i.db.UpdateDevice(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

//updatePoint update point. Called via API call.
func (i *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	modbusDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		modbusErrorMsg("updatePoint(): nil point object")
		return
	}

	/*
		pnt, err := i.db.GetPoint(body.UUID, api.Args{WithPriority: true})
		if pnt == nil || err != nil {
			modbusErrorMsg("could not find pointID: ", pp.FFPointUUID)
			netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
			continue
		}

	*/

	modbusDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	modbusDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	dev, err := i.db.GetDevice(body.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		modbusErrorMsg("updatePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := i.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		modbusErrorMsg("updatePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		i.pointUpdateErr(body, err)
		return
	}

	if utils.BoolIsNil(body.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
		//DO POLLING ENABLE ACTIONS FOR POINT
		//TODO: review these steps to check that UpdatePollingPointByUUID might work better?
		pp := pollqueue.NewPollingPoint(body.UUID, body.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		pp.PollPriority = body.PollPriority
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
		//netPollMan.PollQueue.AddPollingPoint(pp)
		//netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	} else {
		//DO POLLING DISABLE ACTIONS FOR POINT
		netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
	}

	point, err = i.db.UpdatePoint(body.UUID, body, true)
	if err != nil || point == nil {
		modbusErrorMsg("updatePoint(): bad response from UpdatePoint()")
		return nil, err
	}
	return point, nil
}

//writePoint update point. Called via API call.
func (i *Instance) writePoint(pntUUID string, body *model.Point) (point *model.Point, err error) {

	//TODO: check for PointWriteByName calls that might not flow through the plugin.

	modbusDebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		modbusErrorMsg("writePoint(): nil point object")
		return
	}

	modbusDebugMsg(fmt.Sprintf("writePoint() body: %+v\n", body))
	modbusDebugMsg(fmt.Sprintf("writePoint() priority: %+v\n", body.Priority))

	/* TODO: ONLY NEEDED IF THE WRITE VALUE IS WRITTEN ON COV (CURRENTLY IT IS WRITTEN ANYTIME THERE IS A WRITE COMMAND).
	point, err = i.db.GetPoint(pntUUID, api.Args{})
	if err != nil || point == nil {
		modbusErrorMsg("writePoint(): bad response from GetPoint(), ", err)
		return nil, err
	}

	previousWriteVal := -1.11
	if isWriteable(point.WriteMode) {
		previousWriteVal = utils.Float64IsNil(point.WriteValue)
	}
	*/

	body.WritePollRequired = utils.NewTrue()

	point, err = i.db.WritePoint(pntUUID, body, true)
	fmt.Println(fmt.Sprintf("writePoint err %+v", err))
	fmt.Println(fmt.Sprintf("writePoint point %+v", point))
	if err != nil || point == nil {
		modbusErrorMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	/*  TODO: THIS SECTION MIGHT BE USEFUL IF WE ADD ASAP PRIORITY FOR IMMEDIATE POINT WRITES
	dev, err := i.db.GetDevice(body.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		modbusErrorMsg("writePoint(): bad response from GetDevice()")
		return nil, err
	}

	netPollMan, err := i.getNetworkPollManagerByUUID(dev.NetworkUUID)
	if netPollMan == nil || err != nil {
		modbusErrorMsg("writePoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		i.pointUpdateErr(body, err)
		return
	}

	if utils.BoolIsNil(body.Enable) {
		netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
		//DO POLLING ENABLE ACTIONS FOR POINT
		//TODO: review these steps to check that UpdatePollingPointByUUID might work better?
		pp := pollqueue.NewPollingPoint(body.UUID, body.DeviceUUID, dev.NetworkUUID, netPollMan.FFPluginUUID)
		pp.PollPriority = body.PollPriority
		netPollMan.PollingPointCompleteNotification(pp, false, false, 0, true) // This will perform the queue re-add actions based on Point WriteMode. TODO: check function of pointUpdate argument.
		//netPollMan.PollQueue.AddPollingPoint(pp)
		//netPollMan.SetPointPollRequiredFlagsBasedOnWriteMode(pnt)
	} else {
		//DO POLLING DISABLE ACTIONS FOR POINT
		netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
	}
	*/

	return point, nil
}

//deleteNetwork delete network. Called via API call.
func (i *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	modbusDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		modbusErrorMsg("deleteNetwork(): nil network object")
		return
	}
	found := false
	for index, netPollMan := range i.NetworkPollManagers {
		if netPollMan.FFNetworkUUID == body.UUID {
			netPollMan.StopPolling()
			//Next remove the NetworkPollManager from the slice in polling instance
			i.NetworkPollManagers[index] = i.NetworkPollManagers[len(i.NetworkPollManagers)-1]
			i.NetworkPollManagers = i.NetworkPollManagers[:len(i.NetworkPollManagers)-1]
			found = true
		}
	}
	if !found {
		modbusErrorMsg("deleteNetwork(): cannot find NetworkPollManager for network: ", body.UUID)
	}
	ok, err = i.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//deleteNetwork delete device. Called via API call.
func (i *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	modbusDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		modbusErrorMsg("deleteDevice(): nil device object")
		return
	}

	netPollMan, err := i.getNetworkPollManagerByUUID(body.NetworkUUID)
	if netPollMan == nil || err != nil {
		modbusErrorMsg("deleteDevice(): cannot find NetworkPollManager for network: ", body.NetworkUUID)
		return
	}
	netPollMan.PollQueue.RemovePollingPointByDeviceUUID(body.UUID)
	ok, err = i.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//deletePoint delete point. Called via API call.
func (i *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	modbusDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		modbusErrorMsg("deletePoint(): nil point object")
		return
	}

	dev, err := i.db.GetDevice(body.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		modbusErrorMsg("addPoint(): bad response from GetDevice()")
		return false, err
	}

	netPollMan, err := i.getNetworkPollManagerByUUID(dev.NetworkUUID)

	if netPollMan == nil || err != nil {
		modbusErrorMsg("addPoint(): cannot find NetworkPollManager for network: ", dev.NetworkUUID)
		i.pointUpdateErr(body, err)
		return
	}

	netPollMan.PollQueue.RemovePollingPointByPointUUID(body.UUID)
	otherPointsOnSameDeviceExist := netPollMan.PollQueue.CheckPollingQueueForDevUUID(body.DeviceUUID)
	if !otherPointsOnSameDeviceExist {
		netPollMan.PollQueue.RemoveDeviceFromActiveDevicesList(body.DeviceUUID)
	}
	ok, err = i.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//pointUpdate update point. Called from within plugin.
func (i *Instance) pointUpdate(point *model.Point, value float64, writeSuccess, readSuccess, clearFaults bool) (*model.Point, error) {
	if clearFaults {
		point.CommonFault.InFault = false
		point.CommonFault.MessageLevel = model.MessageLevel.Info
		point.CommonFault.MessageCode = model.CommonFaultCode.Ok
		point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
		point.CommonFault.LastOk = time.Now().UTC()
	}

	if readSuccess {
		if value != utils.Float64IsNil(point.OriginalValue) {
			point.ValueUpdatedFlag = utils.NewTrue() //Flag so that UpdatePointValue() will broadcast new value to producers. TODO: MAY NOT BE NEEDED.
		}
		point.OriginalValue = utils.NewFloat64(value)
	}
	point.InSync = utils.NewTrue() //TODO: MAY NOT BE NEEDED.

	_, err = i.db.UpdatePoint(point.UUID, point, true)
	if err != nil {
		modbusErrorMsg("MODBUS UPDATE POINT UpdatePointPresentValue() error: ", err)
		return nil, err
	}
	return point, nil
}

//pointUpdateErr update point with errors. Called from within plugin.
func (i *Instance) pointUpdateErr(point *model.Point, err error) (*model.Point, error) {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	_, err = i.db.UpdatePoint(point.UUID, point, true)
	if err != nil {
		modbusErrorMsg(" pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}

//listSerialPorts list all serial ports on host
func (inst *Instance) listSerialPorts() (*utils.Array, error) {
	ports, err := serial.GetPortsList()
	p := utils.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}
