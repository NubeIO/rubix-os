package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
// addNetwork add network. Called via API call (or wizard)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	if body == nil {
		inst.edge28ErrorMsg("addNetwork(): nil network object")
		return nil, errors.New("empty network body, no network created")
	}
	inst.edge28DebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body, true)
	if network == nil || err != nil {
		inst.edge28ErrorMsg("addNetwork(): failed to create modbus network: ", body.Name)
		return nil, errors.New("failed to create modbus network")
	}
	return network, nil
}

// addDevice add device. Called via API call (or wizard)
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.edge28DebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.edge28DebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.edge28DebugMsg("addDevice(): failed to create modbus device: ", body.Name)
		return nil, errors.New("failed to create modbus device")
	}

	inst.edge28DebugMsg("addDevice(): ", body.UUID)
	// NOTHING TO DO ON DEVICE CREATED
	return device, nil
}

// addPoint add point. Called via API call (or wizard)
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
	objectType, isOutput, isTypeBool := selectObjectType(body.IoNumber)
	body.ObjectType = objectType
	if objectType == "" {
		errMsg := fmt.Sprintf("point object type can not be empty")
		inst.edge28ErrorMsg(errMsg)
		return nil, errors.New(errMsg)
	}
	body.IsOutput = nils.NewBool(isOutput)
	body.IsTypeBool = nils.NewBool(isTypeBool)

	// point, err = inst.db.CreatePoint(body, true, false)
	point, err = inst.db.CreatePoint(body, true, true)
	if point == nil || err != nil {
		inst.edge28DebugMsg("addPoint(): failed to create edge28 point: ", body.Name)
		return nil, errors.New("failed to create edge28 point")
	}
	inst.edge28DebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	return point, nil
}

// updateNetwork update network. Called via API call.
func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.edge28DebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("updateNetwork():  nil network object")
		return
	}
	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

// updateDevice update device. Called via API call.
func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.edge28DebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("updateDevice(): nil device object")
		return
	}

	dev, err := inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil || dev == nil {
		return nil, err
	}
	return dev, nil
}

// updatePoint update point. Called via API call.
func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.edge28DebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.edge28DebugMsg("updatePoint(): nil point object")
		return
	}

	inst.edge28DebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.edge28DebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	point, err = inst.db.UpdatePoint(body.UUID, body, true)
	if err != nil || point == nil {
		inst.edge28DebugMsg("updatePoint(): bad response from UpdatePoint()")
		return nil, err
	}

	dev, err := inst.db.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil || dev == nil {
		inst.edge28DebugMsg("updatePoint(): bad response from GetDevice()")
		return nil, err
	}
	return point, nil
}

// writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {

	// TODO: check for PointWriteByName calls that might not flow through the plugin.

	inst.edge28DebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.edge28DebugMsg("writePoint(): nil point object")
		return
	}

	inst.edge28DebugMsg(fmt.Sprintf("writePoint() body: %+v", body))
	inst.edge28DebugMsg(fmt.Sprintf("writePoint() priority: %+v", body.Priority))

	/* TODO: ONLY NEEDED IF THE WRITE VALUE IS WRITTEN ON COV (CURRENTLY IT IS WRITTEN ANYTIME THERE IS A WRITE COMMAND).
	point, err = inst.db.GetPoint(pntUUID, apinst.Args{})
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

	point, err = inst.db.WritePoint(pntUUID, body, true)
	if err != nil || point == nil {
		inst.edge28DebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}
	return point, nil
}

// deleteNetwork delete network. Called via API call.
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

// deletePoint delete point. Called via API call.
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

// pointUpdate update point. Called from within plugin.
func (inst *Instance) pointUpdate(point *model.Point, value float64, writeSuccess, readSuccess, clearFaults bool) (*model.Point, error) {
	if clearFaults {
		point.CommonFault.InFault = false
		point.CommonFault.MessageLevel = model.MessageLevel.Info
		point.CommonFault.MessageCode = model.CommonFaultCode.Ok
		point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
		point.CommonFault.LastOk = time.Now().UTC()
	}

	if readSuccess {
		if value != float.NonNil(point.OriginalValue) {
			point.ValueUpdatedFlag = boolean.NewTrue() // Flag so that UpdatePointValue() will broadcast new value to producers. TODO: MAY NOT BE NEEDED.
		}
		point.OriginalValue = float.New(value)
	}
	point.InSync = boolean.NewTrue() // TODO: MAY NOT BE NEEDED.

	_, err = inst.db.UpdatePoint(point.UUID, point, true)
	if err != nil {
		inst.edge28DebugMsg("EDGE28 UPDATE POINT UpdatePointPresentValue() error: ", err)
		return nil, err
	}
	return point, nil
}

// pointUpdateErr update point with errors. Called from within plugin.
func (inst *Instance) pointUpdateErr(point *model.Point, err error) (*model.Point, error) {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdatePoint(point.UUID, point, true)
	if err != nil {
		inst.edge28DebugMsg(" pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}

func selectObjectType(selectedPlugin string) (objectType string, isOutput, isTypeBool bool) {
	isOutput = false
	isOutput = false
	switch selectedPlugin {
	case PointsList.R1.IoNumber, PointsList.R2.IoNumber:
		objectType = PointsList.R1.ObjectType
		isOutput = true
		isTypeBool = true
	case PointsList.UO1.IoNumber:
		objectType = PointsList.UO1.ObjectType
		isOutput = true
	case PointsList.DO1.IoNumber:
		objectType = PointsList.DO1.ObjectType
		isOutput = true
		isTypeBool = true
	case PointsList.UI1.IoNumber:
		objectType = PointsList.UI1.ObjectType
	case PointsList.DI1.IoNumber:
		objectType = PointsList.DI1.ObjectType
		isTypeBool = true
	}
	return

}

var Rls = []string{"R1", "R2"}
var DOs = []string{"DO1", "DO2", "DO3", "DO4", "DO5"}
var UOs = []string{"UO1", "UO2", "UO3", "UO4", "UO5", "UO6", "UO7"}
var UIs = []string{"UI1", "UI2", "UI3", "UI4", "UI5", "UI6", "UI7"}
var DIs = []string{"DI1", "DI2", "DI3", "DI4", "DI5", "DI6", "DI7"}

type Point struct {
	IoNumber   string // R1
	ObjectType string // binary_output
	IsOutput   *bool
	IsTypeBool *bool
}

var PointsList = struct {
	R1  Point  `json:"R1"`
	R2  Point  `json:"R2"`
	DO1 Point  `json:"DO1"`
	DO2 string `json:"DO2"`
	DO3 string `json:"DO3"`
	DO4 string `json:"DO4"`
	DO5 string `json:"DO5"`
	UO1 Point  `json:"UO1"`
	UO2 Point  `json:"UO2"`
	UO3 string `json:"UO3"`
	UO4 string `json:"UO4"`
	UO5 string `json:"UO5"`
	UO6 string `json:"UO6"`
	UO7 string `json:"UO7"`
	UI1 Point  `json:"UI1"`
	UI2 string `json:"UI2"`
	UI3 string `json:"UI3"`
	UI4 string `json:"UI4"`
	UI5 string `json:"UI5"`
	UI6 string `json:"UI6"`
	UI7 string `json:"UI7"`
	DI1 Point  `json:"DI1"`
	DI2 string `json:"DI2"`
	DI3 string `json:"DI3"`
	DI4 string `json:"DI4"`
	DI5 string `json:"DI5"`
	DI6 string `json:"DI6"`
	DI7 string `json:"DI7"`
}{
	R1:  Point{IoNumber: "RI", ObjectType: "binary_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	R2:  Point{IoNumber: "R2", ObjectType: "binary_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	DO1: Point{IoNumber: "RI", ObjectType: "binary_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	DO2: "DO2",
	DO3: "DO3",
	DO4: "DO4",
	DO5: "DO5",
	UO1: Point{IoNumber: "UO1", ObjectType: "analog_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	UO2: Point{IoNumber: "UO2", ObjectType: "analog_output", IsOutput: nils.NewTrue(), IsTypeBool: nils.NewTrue()},
	UO3: "UO3",
	UO4: "UO4",
	UO5: "UO5",
	UO6: "UO6",
	UO7: "UO7",
	UI1: Point{IoNumber: "UI1", ObjectType: "analog_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	UI2: "UI2",
	UI3: "UI3",
	UI4: "UI4",
	UI5: "UI5",
	UI6: "UI6",
	UI7: "UI7",
	DI1: Point{IoNumber: "DI1", ObjectType: "binary_value", IsOutput: nils.NewFalse(), IsTypeBool: nils.NewFalse()},
	DI2: "DI2",
	DI3: "DI3",
	DI4: "DI4",
	DI5: "DI5",
	DI6: "DI6",
	DI7: "DI7",
}

var pointList = struct {
	R1  string `json:"R1"`
	R2  string `json:"R2"`
	DO1 string `json:"DO1"`
	DO2 string `json:"DO2"`
	DO3 string `json:"DO3"`
	DO4 string `json:"DO4"`
	DO5 string `json:"DO5"`
	UO1 string `json:"UO1"`
	UO2 string `json:"UO2"`
	UO3 string `json:"UO3"`
	UO4 string `json:"UO4"`
	UO5 string `json:"UO5"`
	UO6 string `json:"UO6"`
	UO7 string `json:"UO7"`
	UI1 string `json:"UI1"`
	UI2 string `json:"UI2"`
	UI3 string `json:"UI3"`
	UI4 string `json:"UI4"`
	UI5 string `json:"UI5"`
	UI6 string `json:"UI6"`
	UI7 string `json:"UI7"`
	DI1 string `json:"DI1"`
	DI2 string `json:"DI2"`
	DI3 string `json:"DI3"`
	DI4 string `json:"DI4"`
	DI5 string `json:"DI5"`
	DI6 string `json:"DI6"`
	DI7 string `json:"DI7"`
}{
	R1:  "R1",
	R2:  "R2",
	DO1: "DO1",
	DO2: "DO2",
	DO3: "DO3",
	DO4: "DO4",
	DO5: "DO5",
	UO1: "UO1",
	UO2: "UO2",
	UO3: "UO3",
	UO4: "UO4",
	UO5: "UO5",
	UO6: "UO6",
	UO7: "UO7",
	UI1: "UI1",
	UI2: "UI2",
	UI3: "UI3",
	UI4: "UI4",
	UI5: "UI5",
	UI6: "UI6",
	UI7: "UI7",
	DI1: "DI1",
	DI2: "DI2",
	DI3: "DI3",
	DI4: "DI4",
	DI5: "DI5",
	DI6: "DI6",
	DI7: "DI7",
}

func pointsAll() []string {
	out := append(Rls, DOs...)
	out = append(out, UOs...)
	out = append(out, DIs...)
	out = append(out, UIs...)
	return out
}

var UOTypes = struct {
	// RAW  string
	DIGITAL string
	PERCENT string
	VOLTSDC string
	// MILLIAMPS  string
}{
	// RAW:  "RAW",
	DIGITAL: "DIGITAL",
	PERCENT: "PERCENT",
	VOLTSDC: "0-10VDC",
	// MILLIAMPS:  "4-20mA",
}

var UITypes = struct {
	RAW              string
	DIGITAL          string
	PERCENT          string
	VOLTSDC          string
	MILLIAMPS        string
	RESISTANCE       string
	THERMISTOR10KT2  string
	THERMISTOR10KT3  string
	THERMISTOR20KT1  string
	THERMISTORPT100  string
	THERMISTORPT1000 string
}{
	RAW:              "RAW",
	DIGITAL:          "DIGITAL",
	PERCENT:          "PERCENT",
	VOLTSDC:          "0-10VDC",
	MILLIAMPS:        "4-20mA",
	RESISTANCE:       "RESISTANCE",
	THERMISTOR10KT2:  "THERMISTOR_10K_TYPE2",
	THERMISTOR10KT3:  "THERMISTOR_10K_TYPE3",
	THERMISTOR20KT1:  "THERMISTOR_20K_TYPE1",
	THERMISTORPT100:  "THERMISTOR_PT100",
	THERMISTORPT1000: "THERMISTOR_PT1000",
}
