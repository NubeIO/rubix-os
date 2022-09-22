package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/array"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"go.bug.st/serial"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	inst.tmvDebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		inst.tmvErrorMsg("addNetwork(): failed to create tmv network: ", body.Name)
		return nil, errors.New("failed to create tmv network")
	}

	if boolean.IsFalse(network.Enable) {
		err = inst.networkUpdateErr(network, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError)
		err = inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError, true)
	}
	return network, nil
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.tmvDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.tmvDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.tmvDebugMsg("addDevice(): failed to create tmv device: ", body.Name)
		return nil, errors.New("failed to create tmv device")
	}

	inst.tmvDebugMsg("addDevice(): ", body.UUID)

	if boolean.IsFalse(device.Enable) {
		err = inst.deviceUpdateErr(device, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
		err = inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	}

	// NOTHING TO DO ON DEVICE CREATED
	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.tmvDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.tmvDebugMsg("addPoint(): ", body.Name)

	point, err = inst.db.CreatePoint(body, true, true)
	if point == nil || err != nil {
		inst.tmvDebugMsg("addPoint(): failed to create tmv point: ", body.Name)
		return nil, errors.New("failed to create tmv point")
	}
	inst.tmvDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	if boolean.IsFalse(point.Enable) {
		err = inst.pointUpdateErr(point, "point disabled", model.MessageLevel.Warning, model.CommonFaultCode.PointError)
	}
	return point, nil

}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.tmvDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("updateNetwork():  nil network object")
		return
	}

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

	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil || network == nil {
		return nil, err
	}

	if boolean.IsFalse(network.Enable) {
		// DO POLLING DISABLE ACTIONS
		inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError, true)
	}

	network, err = inst.db.UpdateNetwork(body.UUID, network, true)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.tmvDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("updateDevice(): nil device object")
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
	if err != nil || device == nil {
		return nil, err
	}

	if boolean.IsFalse(device.Enable) {
		// DO POLLING DISABLE ACTIONS FOR DEVICE
		inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	} else {
		// TODO: Currently on every device update, all device points are removed, and re-added.
		device.CommonFault.InFault = false
		device.CommonFault.MessageLevel = model.MessageLevel.Info
		device.CommonFault.MessageCode = model.CommonFaultCode.Ok
		device.CommonFault.Message = ""
		device.CommonFault.LastOk = time.Now().UTC()
	}

	device, err = inst.db.UpdateDevice(device.UUID, device, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.tmvDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("updatePoint(): nil point object")
		return
	}

	inst.tmvDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.tmvDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "point disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	}

	point, err = inst.db.UpdatePoint(body.UUID, body, true, true)
	if err != nil || point == nil {
		inst.tmvDebugMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}
	return point, nil
}

func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.

	point = nil
	inst.tmvDebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.tmvDebugMsg("writePoint(): nil point object")
		return
	}

	inst.tmvDebugMsg(fmt.Sprintf("writePoint() body: %+v", body))
	inst.tmvDebugMsg(fmt.Sprintf("writePoint() priority: %+v", body.Priority))

	point, _, _, _, err = inst.db.PointWrite(pntUUID, body, false)
	if err != nil {
		inst.tmvDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	return point, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.tmvDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("deleteNetwork(): nil network object")
		return
	}

	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.tmvDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("deleteDevice(): nil device object")
		return
	}

	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.tmvDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("deletePoint(): nil point object")
		return
	}

	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// THE FOLLOWING FUNCTIONS ARE CALLED FROM WITHIN THE PLUGIN
func (inst *Instance) updatePointNames() error {
	nets, err := inst.db.GetNetworksByPlugin("lorawan", api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}
	for _, net := range nets {
		for _, dev := range net.Devices {
			for _, pnt := range dev.Points {
				newPointName := ""
				switch pnt.Name {
				case "digital-1":
					newPointName = "APP-FAULT"
				case "digital-2":
					newPointName = "FLOW_STATUS"
				case "temp-3":
					newPointName = "FLOW_TEMPERATURE"
				case "digital-4":
					newPointName = "OVER_TEMPERATURE_WARN"
				case "digital-5":
					newPointName = "OVER_TEMPERATURE_ALERT"
				case "int_8-6":
					newPointName = "DAILY_TEMP_TEST_1"
				case "int_8-7":
					newPointName = "DAILY_TEMP_TEST_2"
				case "int_8-8":
					newPointName = "DAILY_TEMP_TEST_3"
				case "int_8-9":
					newPointName = "MONTHLY_MEAN_TEMP_TEST"
				case "uint_32-10":
					newPointName = "TOTAL_FLOW_ACCUMULATION"
				case "uint_16-11":
					newPointName = "ONE_DAY_FLOW_ACCUMULATION"
				case "digital-12":
					newPointName = "ONE_DAY_LOW_FLOW_ALERT"
				case "int_8-13":
					newPointName = "DAYS_OF_LOW_FLOW"
				case "digital-14":
					newPointName = "FIVE_DAY_LOW_FLOW_ALERT"
				case "digital-15":
					newPointName = "MONTHLY_HOT_FLUSH_STATUS"
				case "int_8-16":
					newPointName = "OVER_TEMPERATURE_WARN_COUNT"
				case "int_8-17":
					newPointName = "OVER_TEMPERATURE_ALERT_COUNT"
				case "digital-18":
					newPointName = "SOLENOID_STATUS"
				case "digital-19":
					newPointName = "ENABLE"
				case "temp-20":
					newPointName = "TEMPERATURE_SP"
				case "temp-21":
					newPointName = "OVER_TEMPERATURE_OFFSET"
				case "int_16-22":
					newPointName = "LOW_FLOW_THRESHOLD"
				case "temp-23":
					newPointName = "HOT_FLUSH_SP"
				case "uint_16-24":
					newPointName = "HOT_FLUSH_DELAY"
				case "digital-25":
					newPointName = "RESET_ALL"
				case "digital-26":
					newPointName = "SOLENOID_ALLOW"
				case "uint_16-27":
					newPointName = "OVERTEMP_ALERT_DURATION_SP"
				case "temp-28":
					newPointName = "TEMP_CALIBRATION_OFFSET"
				}
				if newPointName != "" {
					pnt.Name = newPointName
					pnt, err = inst.db.UpdatePoint(pnt.UUID, pnt, true, false)
				}
			}
		}
	}
	return nil
}

func (inst *Instance) pointUpdateErr(point *model.Point, message string, messageLevel string, messageCode string) error {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = messageLevel
	point.CommonFault.MessageCode = messageCode
	point.CommonFault.Message = message
	point.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdatePointErrors(point.UUID, point)
	if err != nil {
		inst.tmvErrorMsg(" pointUpdateErr()", err)
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
		inst.tmvErrorMsg(" deviceUpdateErr()", err)
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
		inst.tmvErrorMsg(" networkUpdateErr()", err)
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
