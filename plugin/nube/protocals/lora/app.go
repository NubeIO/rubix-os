package main

import (
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

// addDevicePoints add all points related to a device
func (inst *Instance) addDevicePoints(deviceBody *model.Device) (*model.Point, error) {
	point := new(model.Point)
	point.DeviceUUID = deviceBody.UUID
	point.AddressUUID = deviceBody.AddressUUID
	point.IsProducer = utils.NewFalse()
	point.IsConsumer = utils.NewFalse()
	point.IsOutput = utils.NewFalse()

	for _, pointName := range getDevicePointList(deviceBody) {
		point.Name = fmt.Sprintf("%s_%s_%s", model.TransProtocol.Lora, deviceBody.AddressUUID, pointName)
		point.IoID = pointName
		if inst.addPoint(point) != nil {
			log.Errorf("lora: issue on addPoint: %v\n", err)
			return nil, err
		}
	}
	return nil, nil
}

// addPoint add a pnt
func (inst *Instance) addPoint(body *model.Point) error {
	_, err := inst.db.CreatePoint(body)
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
func (inst *Instance) updatePointValue(body *model.Point, value float64, sensorType string) error {
	addr := body.AddressUUID
	pnt, err := inst.db.GetPointByFieldAndIOID("address_uuid", addr, body)
	if err != nil {
		log.Errorf("lora: issue on failed to find point: %v\n", err)
		return err
	}

	pnt.PresentValue = &value
	pnt.CommonFault.InFault = false
	pnt.CommonFault.MessageLevel = model.MessageLevel.Info
	pnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
	pnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	pnt.CommonFault.LastOk = time.Now().UTC()

	// TODO: fix this for all points if they need conversion
	if pnt.IoType != "" && pnt.IoType != string(model.IOType.RAW) {
		*body.PresentValue = decoder.MicroEdgePointType(pnt.IoType, *body.PresentValue)
	}

	_, err = inst.db.UpdatePoint(pnt.UUID, body, true)
	log.Infof("lora UpdatePoint { AddressUUID: %s value:%v IoID:%s Type:%s }\n", addr, *body.PresentValue, body.IoID, sensorType)
	if err != nil {
		log.Errorf("lora: issue on UpdatePoint: %v\n", err)
		return err
	}

	return nil
}

// updateDevicePointValues update all points under a device within commonSensorData and sensorStruct
func (inst *Instance) updateDevicePointValues(commonSensorData *decoder.CommonValues, sensorStruct interface{}) {
	pnt := new(model.Point)
	pnt.AddressUUID = commonSensorData.Id

	sensorRefl := reflect.ValueOf(sensorStruct)

	// TODO: check this isn't already done in CommonValues
	pnt.IoID = getStructFieldJSONNameByName(sensorRefl, "Rssi")
	err := inst.updatePointValue(pnt, float64(commonSensorData.Rssi), commonSensorData.Sensor)
	if err != nil {
		return
	}

	for i := 0; i < sensorRefl.NumField(); i++ {
		var value float64 = 0.0

		// TODO: check if this is needed
		pnt.IoID = getStructFieldJSONNameByIndex(sensorRefl, i)

		switch sensorRefl.Field(i).Kind() {
		case reflect.String:
			// TODO: handle strings
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			value = sensorRefl.Field(i).Float()
		case reflect.Bool:
			value = BoolToFloat(sensorRefl.Field(i).Bool())
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Int:
			value = float64(sensorRefl.Field(i).Int())
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			value = float64(sensorRefl.Field(i).Uint())
		}

		err := inst.updatePointValue(pnt, value, commonSensorData.Sensor)
		if err != nil {
			return
		}
	}
}
