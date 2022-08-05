package main

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/bugs"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	nets, err := inst.db.GetNetworksByPluginName(body.PluginPath, api.Args{})
	if err != nil {
		return nil, err
	}
	for _, net := range nets {
		if net != nil {
			errMsg := fmt.Sprintf("lora: only max one network is allowed with lora")
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	device, _ = inst.db.GetDeviceByArgs(api.Args{AddressUUID: body.AddressUUID})
	if device != nil {
		errMsg := fmt.Sprintf("lora: the lora ID (address_uuid) must be unique: %s", nils.StringIsNil(body.AddressUUID))
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	err = inst.addDevicePoints(device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	point, err = inst.db.CreatePoint(body, true, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	err = inst.updateDevicePointsAddress(device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

// writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	point, _, _, _, err = inst.db.PointWrite(pntUUID, body, false)
	return point, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) networkUpdateSuccess(uuid string) error {
	var network model.Network
	network.CommonFault.InFault = false
	network.CommonFault.MessageLevel = model.MessageLevel.Info
	network.CommonFault.MessageCode = model.CommonFaultCode.Ok
	network.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	network.CommonFault.LastOk = time.Now().UTC()
	err := inst.db.UpdateNetworkErrors(uuid, &network)
	if err != nil {
		log.Error(bugs.DebugPrint(name, inst.networkUpdateSuccess, err))
	}
	return err
}

func (inst *Instance) networkUpdateErr(uuid, port string, e error) error {
	var network model.Network
	network.CommonFault.InFault = true
	network.CommonFault.MessageLevel = model.MessageLevel.Fail
	network.CommonFault.MessageCode = model.CommonFaultCode.NetworkError
	network.CommonFault.Message = fmt.Sprintf(" port: %s message: %s", port, e.Error())
	network.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdateNetworkErrors(uuid, &network)
	if err != nil {
		log.Error(bugs.DebugPrint(name, inst.networkUpdateErr, err))
	}
	return err
}

func (inst *Instance) deviceUpdateSuccess(uuid string) error {
	var device model.Device
	device.CommonFault.InFault = false
	device.CommonFault.MessageLevel = model.MessageLevel.Info
	device.CommonFault.MessageCode = model.CommonFaultCode.Ok
	device.CommonFault.Message = fmt.Sprintf("lastMessage: %s", utilstime.TimeStamp())
	device.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdateDeviceErrors(uuid, &device)
	if err != nil {
		log.Error("lora-app deviceUpdateErr()", err)
	}
	return err
}

func (inst *Instance) deviceUpdateErr(uuid string, err error) error {
	var device model.Device
	device.CommonFault.InFault = true
	device.CommonFault.MessageLevel = model.MessageLevel.Fail
	device.CommonFault.MessageCode = model.CommonFaultCode.DeviceError
	device.CommonFault.Message = fmt.Sprintf(" error: %s", err.Error())
	device.CommonFault.LastFail = time.Now().UTC()
	err = inst.db.UpdateDeviceErrors(uuid, &device)
	if err != nil {
		log.Error("lora-app deviceUpdateErr()", err)
	}
	return err
}

func (inst *Instance) handleSerialPayload(data string) {
	commonData, fullData := decoder.DecodePayload(data)
	deviceUUID := commonData.ID
	if deviceUUID != "" {
		dev, err := inst.db.GetDeviceByArgs(api.Args{AddressUUID: nils.NewString(deviceUUID)})
		if err != nil {
			errMsg := fmt.Sprintf("lora: issue on failed to find device: %v id: %s\n", err.Error(), deviceUUID)
			log.Errorf(errMsg)
			if dev != nil {
				_ = inst.deviceUpdateErr(dev.UUID, errors.New(errMsg))
			}
			return
		}
		if dev != nil {
			log.Println("lora: sensor-found", deviceUUID, "sensor rssi:", commonData.Rssi)
			_ = inst.deviceUpdateSuccess(dev.UUID)
		}
	}

	if fullData != nil {
		inst.updateDevicePointValues(commonData, fullData)
	}
}

// TODO: need better way to add/update CommonValues points instead of adding/updating the rssi point manually in each func
// addDevicePoints add all points related to a device
func (inst *Instance) addDevicePoints(deviceBody *model.Device) error {
	network, err := inst.db.GetNetwork(deviceBody.NetworkUUID, api.Args{})
	if err != nil {
		log.Errorln("lora: addDevicePoints(), get network", err)
		return err
	}
	if network.PluginPath != "lora" {
		log.Errorln("lora: incorrect network plugin type, must be lora, network was:", network.PluginPath)
		return errors.New("lora: incorrect network plugin type, must be lora")
	}

	points := decoder.GetDevicePointsStruct(deviceBody)
	// TODO: should check this before the device is even added in the wizard
	if points == struct{}{} {
		log.Errorln("lora: addDevicePoints() incorrect device model, try THLM", err)
		return errors.New("lora: addDevicePoints() no device description or points found for this device")
	}
	pointsRefl := reflect.ValueOf(points)

	// kinda poor repeating this but oh well
	pointName := getStructFieldJSONNameByName(decoder.CommonValues{}, "Rssi")
	point := new(model.Point)

	inst.setNewPointFields(deviceBody, point, pointName)
	_, err = inst.addPoint(point)
	if err != nil {
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
				_ = inst.addPointsFromStruct(deviceBody, pointsRefl.Field(i))
			}
			continue
		}
		pointName := getReflectFieldJSONName(pointsRefl.Type().Field(i))
		inst.setNewPointFields(deviceBody, point, pointName)
		_, err := inst.addPoint(point)
		if err != nil {
			log.Errorf("lora: issue on addPoint: %v\n", err)
			return err
		}
	}
	return nil
}

func (inst *Instance) setNewPointFields(deviceBody *model.Device, pointBody *model.Point, name string) {
	pointBody.DeviceUUID = deviceBody.UUID
	pointBody.AddressUUID = deviceBody.AddressUUID
	pointBody.IsProducer = boolean.NewFalse()
	pointBody.IsConsumer = boolean.NewFalse()
	pointBody.IsOutput = boolean.NewFalse()
	pointBody.Name = fmt.Sprintf("%s", name)
	pointBody.IoNumber = name
}

// updateDevicePointsAddress by its lora id and type as in temp or lux
func (inst *Instance) updateDevicePointsAddress(body *model.Device) error {
	dev, err := inst.db.GetDevice(body.UUID, api.Args{WithPoints: true})
	if err != nil {
		return err
	}
	for _, pt := range dev.Points {
		pt.AddressUUID = body.AddressUUID
		_, err = inst.db.UpdatePoint(pt.UUID, pt, true, false)
		if err != nil {
			log.Errorf("lora: issue on UpdatePoint updateDevicePointsAddress(): %v\n", err)
			return err
		}
	}
	return nil
}

// TODO: update to make more efficient for updating just the value (incl fault etc.)
func (inst *Instance) updatePointValue(body *model.Point, value float64) error {
	// TODO: fix this so don't need to request the point for the UUID before hand
	pnt, err := inst.db.GetOnePointByArgs(api.Args{AddressUUID: body.AddressUUID, IoNumber: &body.IoNumber})
	if err != nil {
		log.Errorf("lora: issue on failed to find point: %v name: %s IO-ID:%s\n", err, body.AddressUUID, body.IoNumber)
		return err
	}
	priority := map[string]*float64{"_16": &value}
	if pnt.IoType != "" && pnt.IoType != string(model.IOTypeRAW) {
		if body.PresentValue == nil {
			priority["_16"] = nil
		}
		priority["_16"] = float.New(decoder.MicroEdgePointType(pnt.IoType, *body.PresentValue))
	}
	pointWriter := model.PointWriter{Priority: &priority}
	_, _, _, _, err = inst.db.PointWrite(pnt.UUID, &pointWriter, true)
	if err != nil {
		log.Error("lora: UpdatePointValue()", err)
	}
	return err
}

// updateDevicePointValues update all points under a device within commonSensorData and sensorStruct
func (inst *Instance) updateDevicePointValues(commonValues *decoder.CommonValues, sensorStruct interface{}) {
	// manually update rssi + any other CommonValues
	pnt := new(model.Point)
	pnt.AddressUUID = &commonValues.ID
	pnt.IoNumber = getStructFieldJSONNameByName(sensorStruct, "Rssi")
	err := inst.updatePointValue(pnt, float64(commonValues.Rssi))
	if err != nil {
		return
	}
	// update all other fields in sensorStruct
	inst.updateDevicePointValuesStruct(commonValues.ID, sensorStruct)
}

func (inst *Instance) updateDevicePointValuesStruct(deviceID string, sensorStruct interface{}) {
	pnt := new(model.Point)
	pnt.AddressUUID = &deviceID
	sensorRefl := reflect.ValueOf(sensorStruct)

	for i := 0; i < sensorRefl.NumField(); i++ {
		value := 0.0

		// TODO: check if this is needed
		pnt.IoNumber = getReflectFieldJSONName(sensorRefl.Type().Field(i))

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
