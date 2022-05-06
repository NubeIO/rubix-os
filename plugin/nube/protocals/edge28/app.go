package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

//addDevice add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//addDevice add device
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	network, err := inst.db.GetNetwork(body.NetworkUUID, api.Args{WithDevices: true})
	if err != nil {
		return nil, err
	}
	if len(network.Devices) >= 1 {
		errMsg := fmt.Sprintf("edge-28: only one device is allowed per network")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	return device, nil
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

//addPoint add point
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body.IoNumber == "" {
		body.IoNumber = pointList.UI1
	}
	if body.IoType == "" {
		body.IoType = UITypes.DIGITAL
	}
	objectType, isOutput, isTypeBool := selectObjectType(body.IoNumber)
	body.ObjectType = objectType
	if objectType == "" {
		errMsg := fmt.Sprintf("edge28: point object type can not be empty")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	body.IsOutput = nils.NewBool(isOutput)
	body.IsTypeBool = nils.NewBool(isTypeBool)
	point, err = inst.db.CreatePoint(body, true, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

//updateNetwork update network
func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//updateDevice update device
func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

//updatePoint update point
func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	point, err = inst.db.UpdatePoint(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

//deleteNetwork delete network
func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	//TODO: check for PointWriteByName calls that might not flow through the plugin.
	if body == nil {
		return
	}
	point, err = inst.db.WritePoint(pntUUID, body, true)
	if err != nil || point == nil {
		return nil, err
	}
	return point, nil
}

//deleteNetwork delete device
func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//deletePoint delete point
func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//pointUpdate update point present value
func (inst *Instance) pointUpdate(uuid string) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	_, err := inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("edge28-app: UpdatePoint()", err)
		return nil, err
	}
	return nil, nil
}

//pointUpdate update point present value
func (inst *Instance) pointUpdateValue(uuid string, value float64) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	priority := map[string]*float64{"_16": &value}
	point.InSync = boolean.NewTrue()
	_, err := inst.db.UpdatePointValue(uuid, &point, &priority, true)
	if err != nil {
		log.Error("edge28-app: pointUpdateValue()", err)
		return nil, err
	}
	return nil, nil
}

//pointUpdate update point present value
func (inst *Instance) pointUpdateErr(uuid string, err error) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("edge28-app: pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}

var Rls = []string{"R1", "R2"}
var DOs = []string{"DO1", "DO2", "DO3", "DO4", "DO5"}
var UOs = []string{"UO1", "UO2", "UO3", "UO4", "UO5", "UO6", "UO7"}
var UIs = []string{"UI1", "UI2", "UI3", "UI4", "UI5", "UI6", "UI7"}
var DIs = []string{"DI1", "DI2", "DI3", "DI4", "DI5", "DI6", "DI7"}

type Point struct {
	IoNumber   string //R1
	ObjectType string //binary_output
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
	//RAW  string
	DIGITAL string
	PERCENT string
	VOLTSDC string
	//MILLIAMPS  string
}{
	//RAW:  "RAW",
	DIGITAL: "DIGITAL",
	PERCENT: "PERCENT",
	VOLTSDC: "0-10VDC",
	//MILLIAMPS:  "4-20mA",
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
