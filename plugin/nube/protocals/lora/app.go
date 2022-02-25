package main

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

var err error

func (inst *Instance) handleSerialPayload(data string) {
	commonData, fullData := decoder.DecodePayload(data)
	if fullData != nil {
		inst.updateDevicePointValues(commonData, fullData)
	}
}

// TODO: need better way to add/update CommonValues points instead of
//    adding/updating the rssi point manually in each func

// addDevicePoints add all points related to a device
func (inst *Instance) addDevicePoints(deviceBody *model.Device) error {
	points := decoder.GetDevicePointsStruct(deviceBody)
	// TODO: should check this before the device is even added in the wizard
	if points == struct{}{} {
		return errors.New("no device description or points found for this device")
	}
	pointsRefl := reflect.ValueOf(points)

	// kinda poor repeating this but oh well
	pointName := getStructFieldJSONNameByName(decoder.CommonValues{}, "Rssi")
	point := new(model.Point)
	inst.setnewPointFields(deviceBody, point, pointName)
	if inst.addPoint(point) != nil {
		log.Errorf("lora: issue on addPoint: %v\n", err)
		return err
	}

	return inst.addPointsFromStruct(deviceBody, pointsRefl)
}

func (inst *Instance) addPointsFromStruct(deviceBody *model.Device, pointsRefl reflect.Value) error {

	point := new(model.Point)

	for i := 0; i < pointsRefl.NumField(); i++ {
		if pointsRefl.Field(i).Kind() == reflect.Struct {
			if _, ok := pointsRefl.Field(i).Interface().(decoder.CommonValues); !ok {
				inst.addPointsFromStruct(deviceBody, pointsRefl.Field(i))
			}
			continue
		}

		pointName := getReflectFieldJSONName(pointsRefl.Type().Field(i))
		inst.setnewPointFields(deviceBody, point, pointName)
		if inst.addPoint(point) != nil {
			log.Errorf("lora: issue on addPoint: %v\n", err)
			return err
		}
	}
	return nil
}

func (inst *Instance) setnewPointFields(deviceBody *model.Device, pointBody *model.Point, name string) {
	pointBody.DeviceUUID = deviceBody.UUID
	pointBody.AddressUUID = deviceBody.AddressUUID
	pointBody.IsProducer = utils.NewFalse()
	pointBody.IsConsumer = utils.NewFalse()
	pointBody.IsOutput = utils.NewFalse()
	pointBody.Name = fmt.Sprintf("%s_%s_%s", model.TransProtocol.Lora, deviceBody.AddressUUID, name)
	pointBody.IoID = name
}

// addPoint add a pnt
func (inst *Instance) addPoint(body *model.Point) error {
	_, err := inst.db.CreatePoint(body, false, true)
	if err != nil {
		log.Errorf("lora: issue on CreatePoint: %v\n", err)
		return err
	}
	return nil
}

// updateDevicePointsAddress by its lora id and type as in temp or lux
func (inst *Instance) updateDevicePointsAddress(body *model.Device) error {
	var pnt model.Point
	pnt.AddressUUID = body.AddressUUID
	var arg api.Args
	arg.WithPoints = true
	dev, err := inst.db.GetDevice(body.UUID, arg)
	if err != nil {
		return err
	}
	for _, pt := range dev.Points {
		_, err = inst.db.UpdatePoint(pt.UUID, &pnt, true)
		if err != nil {
			log.Errorf("lora: issue on UpdatePoint: %v\n", err)
			return err
		}
	}
	return nil
}

// TODO: update to make more efficient for updating just the value (incl fault etc.)
func (inst *Instance) updatePointValue(body *model.Point, value float64) error {
	// TODO: fix this so don't need to request the point for the UUID before hand
	pnt, err := inst.db.GetOnePointByArgs(api.Args{AddressUUID: &body.AddressUUID, IoId: &body.IoID})
	if err != nil {
		log.Errorf("lora: issue on failed to find point: %v name: %s IO-ID:%s\n", err, body.AddressUUID, body.IoID)
		return err
	}
	body.CommonFault.InFault = false
	body.CommonFault.MessageLevel = model.MessageLevel.Info
	body.CommonFault.MessageCode = model.CommonFaultCode.Ok
	body.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	body.CommonFault.LastOk = time.Now().UTC()

	var pri model.Priority
	pri.P16 = &value
	body.Priority = &pri
	body.InSync = utils.NewTrue()
	fmt.Println(body)
	if pnt.IoType != "" && pnt.IoType != string(model.IOTypeRAW) {
		*pri.P16 = decoder.MicroEdgePointType(pnt.IoType, *body.PresentValue)
	}
	_, _ = inst.db.UpdatePointValue(pnt.UUID, body, true)
	if err != nil {
		log.Error("lora: UpdatePointValue()", err)
		return err
	}
	return nil
}

// updateDevicePointValues update all points under a device within commonSensorData and sensorStruct
func (inst *Instance) updateDevicePointValues(commonValues *decoder.CommonValues, sensorStruct interface{}) {
	// manually update rssi + any other CommonValues
	pnt := new(model.Point)
	pnt.AddressUUID = commonValues.ID
	pnt.IoID = getStructFieldJSONNameByName(sensorStruct, "Rssi")
	err := inst.updatePointValue(pnt, float64(commonValues.Rssi))
	if err != nil {
		return
	}

	// update all other fields in sensorStruct
	inst.updateDevicePointValuesStruct(commonValues.ID, sensorStruct)
}

func (inst *Instance) updateDevicePointValuesStruct(deviceID string, sensorStruct interface{}) {
	pnt := new(model.Point)
	pnt.AddressUUID = deviceID
	sensorRefl := reflect.ValueOf(sensorStruct)

	for i := 0; i < sensorRefl.NumField(); i++ {
		value := 0.0

		// TODO: check if this is needed
		pnt.IoID = getReflectFieldJSONName(sensorRefl.Type().Field(i))

		switch sensorRefl.Field(i).Kind() {
		case reflect.String:
			// TODO: handle strings
			continue
		case reflect.Float32, reflect.Float64:
			value = sensorRefl.Field(i).Float()
		case reflect.Bool:
			value = BoolToFloat(sensorRefl.Field(i).Bool())
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			value = float64(sensorRefl.Field(i).Int())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			value = float64(sensorRefl.Field(i).Uint())
		case reflect.Struct:
			if _, ok := sensorRefl.Field(i).Interface().(decoder.CommonValues); !ok {
				inst.updateDevicePointValuesStruct(deviceID, sensorRefl.Field(i).Interface())
			}
			continue
		}

		err := inst.updatePointValue(pnt, value)
		if err != nil {
			return
		}
	}
}
