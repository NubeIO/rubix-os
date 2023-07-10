package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body == nil {
		inst.edge28ErrorMsg("addNetwork(): nil network object")
		return nil, errors.New("empty network body, no network created")
	}
	inst.edge28DebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body)
	if network == nil || err != nil {
		inst.edge28ErrorMsg("addNetwork(): failed to create edge28 network: ", body.Name)
		return nil, errors.New("failed to create edge28 network")
	}
	return network, nil
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.edge28DebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.edge28DebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.edge28DebugMsg("addDevice(): failed to create edge28 device: ", body.Name)
		return nil, errors.New("failed to create edge28 device")
	}

	inst.edge28DebugMsg("addDevice(): ", body.UUID)

	// CREATE ALL EDGE28 POINTS
	inst.edge28DebugMsg("addDevice(): ADDING ALL POINTS")
	// for _, e := range pointsAll() {
	//	inst.edge28DebugMsg(e)
	//	var pnt model.Point
	//	pnt.DeviceUUID = device.UUID
	//	pName := e
	//	pnt.Name = pName
	//	pnt.Description = pName
	//	pnt.IoNumber = e
	//	pnt.Fallback = float.New(0)
	//	pnt.COV = float.New(0.1)
	//	pnt.IoType = UITypes.DIGITAL
	//	pnt.Enable = boolean.NewTrue()
	//	inst.edge28DebugMsg(fmt.Sprintf("%+v\n", pnt))
	//	point, err := inst.addPoint(&pnt)
	//	if err != nil || point.UUID == "" {
	//		inst.edge28ErrorMsg("addDevice(): ", "failed to create a new point")
	//	}
	// }
	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.edge28DebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.edge28DebugMsg("addPoint(): ", body.Name)

	if body.IoNumber == "" {
		body.IoNumber = pointList.UI1
	}
	if body.IoType == "" {
		body.IoType = UITypes.DIGITAL
	}
	objectType, isOutput := selectObjectType(body.IoNumber)
	isTypeBool := checkForBooleanType(body.IoType)

	if isOutput {
		body.EnableWriteable = boolean.NewTrue()
	} else {
		body.EnableWriteable = boolean.NewFalse()
	}

	body.ObjectType = objectType
	if objectType == "" {
		errMsg := fmt.Sprintf("point object type can not be empty")
		inst.edge28ErrorMsg(errMsg)
		return nil, errors.New(errMsg)
	}
	body.IsOutput = nils.NewBool(isOutput)
	if isOutput {
		body.PointPriorityArrayMode = model.PriorityArrayToWriteValue
	} else {
		body.PointPriorityArrayMode = model.ReadOnlyNoPriorityArrayRequired
	}

	body.IsTypeBool = nils.NewBool(isTypeBool)
	body.WritePollRequired = nils.NewBool(isOutput) // write value immediately if it is an output point
	body.ReadPollRequired = nils.NewBool(!isOutput) // only read value if it isn't an output point

	if body.WriteValue != nil {
		body.WriteValue = limitValueByEdge28Type(body.IoType, body.WriteValue)
	}

	point, err = inst.db.CreatePoint(body, true)
	if point == nil || err != nil {
		inst.edge28DebugMsg("addPoint(): failed to create edge28 point: ", body.Name)
		return nil, errors.New("failed to create edge28 point")
	}
	inst.edge28DebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	return point, nil
}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.edge28DebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("updateNetwork():  nil network object")
		return
	}
	network, err = inst.db.UpdateNetwork(body.UUID, body)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.edge28DebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("updateDevice(): nil device object")
		return
	}

	dev, err := inst.db.UpdateDevice(body.UUID, body)
	if err != nil || dev == nil {
		return nil, err
	}
	return dev, nil
}

func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.edge28DebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("updatePoint(): nil point object")
		return
	}
	if body.IoNumber == "" {
		body.IoNumber = pointList.UI1
	}
	if body.IoType == "" {
		body.IoType = UITypes.DIGITAL
	}
	_, isOutput := selectObjectType(body.IoNumber)
	isTypeBool := checkForBooleanType(body.IoType)
	body.IsOutput = nils.NewBool(isOutput)
	body.IsTypeBool = nils.NewBool(isTypeBool)

	if isOutput {
		body.EnableWriteable = boolean.NewTrue()
	} else {
		body.EnableWriteable = boolean.NewFalse()
	}

	if body.WriteValue != nil {
		body.WriteValue = limitValueByEdge28Type(body.IoType, body.WriteValue)
	}

	inst.edge28DebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.edge28DebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	point, err = inst.db.UpdatePoint(body.UUID, body)
	if err != nil || point == nil {
		inst.edge28DebugMsg("updatePoint(): bad response from UpdatePoint()")
	}
	return point, err
}

func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.
	inst.edge28DebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.edge28DebugMsg("writePoint(): nil point object")
		return
	}

	// TODO: add code to check through priority array and limit the values by IoType.
	pnt, err := inst.db.GetPoint(pntUUID, argspkg.Args{})
	if err == nil {
		body.Priority = limitPriorityArrayByEdge28Type(pnt.IoType, body)
	}

	/* TODO: ONLY NEEDED IF THE WRITE VALUE IS WRITTEN ON COV (CURRENTLY IT IS WRITTEN ANYTIME THERE IS A WRITE COMMAND).
	point, err = inst.db.GetPoint(pntUUID, inst.Args{})
	if err != nil || point == nil {
		inst.edge28ErrorMsg("writePoint(): bad response from GetPoint(), ", err)
		return nil, err
	}

	previousWriteVal := -1.11
	if isWriteable(point.WriteMode) {
		previousWriteVal = utils.Float64IsNil(point.WriteValue)
	}
	*/

	// body.WritePollRequired = utils.NewTrue() // TODO: commented out this section, seems like useless

	point, _, _, _, err = inst.db.PointWrite(pntUUID, body)
	if err != nil {
		inst.edge28DebugMsg("writePoint(): bad response from WritePoint(), ", err)
	}
	return point, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.edge28DebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("deleteNetwork(): nil network object")
		return
	}
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// deleteNetwork delete device. Called via API call.
func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.edge28DebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("deleteDevice(): nil device object")
		return
	}
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.edge28DebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("deletePoint(): nil point object")
		return
	}
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) pointUpdate(point *model.Point, value float64, readSuccess bool) (*model.Point, error) {
	if readSuccess {
		point.OriginalValue = float.New(value)
	}
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	_, err := inst.db.UpdatePoint(point.UUID, point)
	if err != nil {
		inst.edge28DebugMsg("EDGE28 UPDATE POINT UpdatePointPresentValue() error: ", err)
		return nil, err
	}
	return point, nil
}

func (inst *Instance) pointUpdateErr(point *model.Point, err error) error {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	err = inst.db.UpdatePointErrors(point.UUID, point)
	if err != nil {
		inst.edge28DebugMsg(" pointUpdateErr()", err)
		return err
	}
	return nil
}
