package main

import (
	"errors"
	"fmt"
	argspkg "github.com/NubeIO/rubix-os/args"
	"reflect"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/NubeIO/rubix-os/utils/integer"

	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/bugs"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/plugin/nube/protocals/lora/decoder"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	nets, err := inst.db.GetNetworksByPluginName(body.PluginPath, argspkg.Args{})
	if err != nil {
		return nil, err
	}
	for _, net := range nets {
		if net != nil {
			errMsg := "loraraw: only max one network is allowed with lora"
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	body.TransportType = "serial"
	if integer.IsUnitNil(body.SerialBaudRate) {
		body.SerialBaudRate = integer.NewUint(38400)
	}
	network, err = inst.db.CreateNetwork(body)
	if err != nil {
		return nil, err
	}
	inst.networkUUID = network.UUID
	go inst.run()
	return body, nil
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	*body.AddressUUID = strings.ToUpper(*body.AddressUUID)
	device, _ = inst.db.GetDeviceByArgs(argspkg.Args{AddressUUID: body.AddressUUID})
	if device != nil {
		errMsg := fmt.Sprintf("loraraw: the lora ID (address_uuid) must be unique: %s", nils.StringIsNil(body.AddressUUID))
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	err = inst.addDevicePoints(device)
	if err != nil {
		inst.db.DeleteDevice(device.UUID)
		return nil, err
	}
	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	body.ObjectType = "analog_input"
	body.IoType = string(model.IOTypeRAW)
	body.Name = strings.ToLower(body.Name)
	body.EnableWriteable = boolean.NewFalse()
	point, err = inst.db.CreatePoint(body, true)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (inst *Instance) deletePoint(body *model.Point) (success bool, err error) {
	// TODO: For now this db call has been removed, so that point deletes of lora points is not allowed by the user; can only be deleted by the whole device.
	/*
		success, err = inst.db.DeletePoint(body.UUID)
		if err != nil {
			return false, err
		}
	*/
	return success, nil
}

func (inst *Instance) networkUpdateSuccess(uuid string) error {
	var network model.Network
	network.InFault = false
	network.MessageLevel = model.MessageLevel.Info
	network.MessageCode = model.CommonFaultCode.Ok
	network.Message = model.CommonFaultMessage.NetworkMessage
	network.LastOk = time.Now().UTC()
	err := inst.db.UpdateNetworkErrors(uuid, &network)
	if err != nil {
		log.Error(bugs.DebugPrint(name, inst.networkUpdateSuccess, err))
	}
	return err
}

func (inst *Instance) networkUpdateErr(uuid, port string, e error) error {
	var network model.Network
	network.InFault = true
	network.MessageLevel = model.MessageLevel.Fail
	network.MessageCode = model.CommonFaultCode.NetworkError
	network.Message = fmt.Sprintf(" port: %s message: %s", port, e.Error())
	network.LastFail = time.Now().UTC()
	err := inst.db.UpdateNetworkErrors(uuid, &network)
	if err != nil {
		log.Error(bugs.DebugPrint(name, inst.networkUpdateErr, err))
	}
	return err
}

func (inst *Instance) deviceUpdateSuccess(uuid string) error {
	var device model.Device
	device.InFault = false
	device.MessageLevel = model.MessageLevel.Info
	device.MessageCode = model.CommonFaultCode.Ok
	device.Message = fmt.Sprintf("lastMessage: %s", utilstime.TimeStamp())
	device.LastFail = time.Now().UTC()
	err := inst.db.UpdateDeviceErrors(uuid, &device)
	if err != nil {
		log.Error("lora-app deviceUpdateErr()", err)
	}
	return err
}

func (inst *Instance) deviceUpdateErr(uuid string, err error) error {
	var device model.Device
	device.InFault = true
	device.MessageLevel = model.MessageLevel.Fail
	device.MessageCode = model.CommonFaultCode.DeviceError
	device.Message = fmt.Sprintf(" error: %s", err.Error())
	device.LastFail = time.Now().UTC()
	err = inst.db.UpdateDeviceErrors(uuid, &device)
	if err != nil {
		log.Error("lora-app deviceUpdateErr()", err)
	}
	return err
}

func (inst *Instance) pointUpdateSuccess(point *model.Point) error {
	if point == nil {
		return errors.New("lora-plugin: nil point to pointUpdateSuccess()")
	}
	point.InFault = false
	point.MessageLevel = model.MessageLevel.Info
	point.MessageCode = model.CommonFaultCode.Ok
	point.Message = fmt.Sprintf("lastMessage: %s", utilstime.TimeStamp())
	err := inst.db.UpdatePointSuccess(point.UUID, point)
	if err != nil {
		log.Error("lora-app UpdatePointSuccess()", err)
	}
	return err
}

func (inst *Instance) handleSerialPayload(data string) {
	if inst.networkUUID == "" {
		return
	}
	if !decoder.ValidPayload(data) {
		return
	}
	log.Debug("loraraw: Uplink: ", data)
	device := inst.getDeviceByLoRaAddress(decoder.DecodeAddress(data))
	if device == nil {
		id := decoder.DecodeAddress(data) // show user messages from lora
		rssi := decoder.DecodeRSSI(data)
		log.Infof("lora-raw: message from sensor id: %s rssi: %d", id, rssi)
		return
	}
	devDesc := decoder.GetDeviceDescription(device)
	if devDesc == &decoder.NilLoRaDeviceDescription {
		return
	}
	commonData, fullData := decoder.DecodePayload(data, devDesc)
	if commonData == nil {
		return
	}
	deviceId := commonData.ID
	if deviceId != "" {
		dev, err := inst.db.GetDeviceByArgs(argspkg.Args{AddressUUID: nils.NewString(deviceId)})
		if err != nil {
			errMsg := fmt.Sprintf("lora-raw: issue on failed to find device: %v id: %s\n", err.Error(), deviceId)
			log.Errorf(errMsg)
			if dev != nil {
				_ = inst.deviceUpdateErr(dev.UUID, errors.New(errMsg))
			}
			return
		}
		if dev != nil {
			log.Infof("lora-raw: sensor found id: %s rssi: %d type: %s", commonData.ID, commonData.Rssi, commonData.Sensor)
			_ = inst.deviceUpdateSuccess(dev.UUID)
		}
	}
	if fullData != nil {
		inst.updateDevicePointValues(commonData, fullData, device)
	}
}

func (inst *Instance) getDeviceByLoRaAddress(address string) *model.Device {
	device, err := inst.db.GetDeviceByArgs(argspkg.Args{AddressUUID: &address})
	if err != nil {
		return nil
	}
	return device
}

// TODO: need better way to add/update CommonValues points instead of adding/updating the rssi point manually in each func
// addDevicePoints add all points related to a device
func (inst *Instance) addDevicePoints(deviceBody *model.Device) error {
	network, err := inst.db.GetNetwork(deviceBody.NetworkUUID, argspkg.Args{})
	if err != nil {
		log.Errorln("loraraw: addDevicePoints(), get network", err)
		return err
	}
	if network.PluginPath != "lora" {
		log.Errorln("loraraw: incorrect network plugin type, must be lora, network was:", network.PluginPath)
		return errors.New("loraraw: incorrect network plugin type, must be lora")
	}

	points := decoder.GetDevicePointsStruct(deviceBody)
	// TODO: should check this before the device is even added in the wizard
	if points == struct{}{} {
		log.Errorln("loraraw: addDevicePoints() incorrect device model, try THLM", err)
		return errors.New("loraraw: addDevicePoints() no device description or points found for this device")
	}
	pointsRefl := reflect.ValueOf(points)
	inst.addPointsFromName(deviceBody, "Rssi", "Snr")
	inst.addPointsFromStruct(deviceBody, pointsRefl, "")
	return nil
}

func (inst *Instance) addPointsFromName(deviceBody *model.Device, names ...string) {
	var points []*model.Point
	for _, name := range names {
		pointName := getStructFieldJSONNameByName(decoder.CommonValues{}, name)
		point := new(model.Point)
		inst.setNewPointFields(deviceBody, point, pointName)
		point.EnableWriteable = boolean.NewFalse()
		points = append(points, point)
	}
	inst.savePoints(points)
}

func (inst *Instance) addPointsFromStruct(deviceBody *model.Device, pointsRefl reflect.Value, postfix string) {
	var points []*model.Point
	for i := 0; i < pointsRefl.NumField(); i++ {
		field := pointsRefl.Field(i)
		if field.Kind() == reflect.Struct {
			if _, ok := field.Interface().(decoder.CommonValues); !ok {
				inst.addPointsFromStruct(deviceBody, pointsRefl.Field(i), postfix)
			}
			continue
		} else if field.Kind() == reflect.Array || field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				pf := fmt.Sprintf("%s_%d", postfix, j+1)
				v := field.Index(j)
				inst.addPointsFromStruct(deviceBody, v, pf)
			}
			continue
		}
		pointName := getReflectFieldJSONName(pointsRefl.Type().Field(i))
		if postfix != "" {
			pointName = fmt.Sprintf("%s%s", pointName, postfix)
		}
		point := new(model.Point)
		inst.setNewPointFields(deviceBody, point, pointName)
		point.EnableWriteable = boolean.NewFalse()
		points = append(points, point)
	}
	inst.savePoints(points)
}

func (inst *Instance) savePoints(points []*model.Point) {
	var wg sync.WaitGroup
	for _, point := range points {
		wg.Add(1)
		point := point
		go func() {
			defer wg.Done()
			point.EnableWriteable = boolean.NewFalse()
			_, err := inst.addPoint(point)
			if err != nil {
				log.Errorf("loraraw: issue on addPoint: %v\n", err)
			}
		}()
	}
	wg.Wait()
}

func (inst *Instance) setNewPointFields(deviceBody *model.Device, pointBody *model.Point, name string) {
	pointBody.Enable = boolean.NewTrue()
	pointBody.DeviceUUID = deviceBody.UUID
	pointBody.AddressUUID = deviceBody.AddressUUID
	pointBody.IsOutput = boolean.NewFalse()
	pointBody.Name = cases.Title(language.English).String(name)
	pointBody.IoNumber = name
	pointBody.ThingType = "point"
	pointBody.WriteMode = model.ReadOnly
}

// updateDevicePointsAddress by its lora id and type as in temp or lux
func (inst *Instance) updateDevicePointsAddress(body *model.Device) error {
	dev, err := inst.db.GetDevice(body.UUID, argspkg.Args{WithPoints: true})
	if err != nil {
		return err
	}
	for _, pt := range dev.Points {
		pt.AddressUUID = body.AddressUUID
		pt.EnableWriteable = boolean.NewFalse()
		_, err = inst.db.UpdatePoint(pt.UUID, pt)
		if err != nil {
			log.Errorf("loraraw: issue on UpdatePoint updateDevicePointsAddress(): %v\n", err)
			return err
		}
	}
	return nil
}

// TODO: update to make more efficient for updating just the value (incl fault etc.)
func (inst *Instance) updatePointValue(body *model.Point, value float64, device *model.Device) error {
	// TODO: fix this so don't need to request the point for the UUID before hand
	pnt, err := inst.db.GetOnePointByArgs(argspkg.Args{AddressUUID: body.AddressUUID, IoNumber: &body.IoNumber})
	if err != nil {
		log.Errorf("loraraw: issue on failed to find point: %v name: %s IO-ID:%s\n", err, body.AddressUUID, body.IoNumber)
		return err
	}

	priority := map[string]*float64{"_16": &value}
	if pnt.IoType != "" && pnt.IoType != string(model.IOTypeRAW) {
		priority["_16"] = float.New(decoder.MicroEdgePointType(pnt.IoType, value, device.Model))

	}
	pointWriter := model.PointWriter{Priority: &priority}
	point, _, _, _, err := inst.db.PointWrite(pnt.UUID, &pointWriter) // TODO: look on it, faults messages were cleared out
	if err != nil {
		log.Error("lora-raw: UpdatePointValue()", err)
		return err
	}
	err = inst.pointUpdateSuccess(point)
	return err
}

// updateDevicePointValues update all points under a device within commonSensorData and sensorStruct
func (inst *Instance) updateDevicePointValues(commonValues *decoder.CommonValues, sensorStruct interface{}, device *model.Device) {
	// manually update rssi + any other CommonValues
	pnt := new(model.Point)
	pnt.AddressUUID = &commonValues.ID
	pnt.IoNumber = getStructFieldJSONNameByName(sensorStruct, "Rssi")
	err := inst.updatePointValue(pnt, float64(commonValues.Rssi), device)
	if err != nil {
		return
	}
	pnt.IoNumber = getStructFieldJSONNameByName(sensorStruct, "Snr")
	err = inst.updatePointValue(pnt, float64(commonValues.Snr), device)
	if err != nil {
		return
	}
	// update all other fields in sensorStruct
	inst.updateDevicePointValuesStruct(commonValues.ID, sensorStruct, "", device)
}

func (inst *Instance) updateDevicePointValuesStruct(deviceID string, sensorStruct interface{}, postfix string, device *model.Device) {
	pnt := new(model.Point)
	pnt.AddressUUID = &deviceID
	sensorRefl := reflect.ValueOf(sensorStruct)

	for i := 0; i < sensorRefl.NumField(); i++ {
		value := 0.0
		pnt.IoNumber = fmt.Sprintf("%s%s", getReflectFieldJSONName(sensorRefl.Type().Field(i)), postfix)
		field := sensorRefl.Field(i)

		switch field.Kind() {
		case reflect.Float32, reflect.Float64:
			value = field.Float()
		case reflect.Bool:
			value = BoolToFloat(field.Bool())
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			value = float64(field.Int())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			value = float64(field.Uint())
		case reflect.Struct:
			if _, ok := field.Interface().(decoder.CommonValues); !ok {
				inst.updateDevicePointValuesStruct(deviceID, field.Interface(), postfix, device)
			}
			continue
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				pf := fmt.Sprintf("%s_%d", postfix, j+1)
				v := field.Index(j).Interface()
				inst.updateDevicePointValuesStruct(deviceID, v, pf, device)
			}
			continue
		default:
			continue
		}

		err := inst.updatePointValue(pnt, value, device)
		if err != nil {
			return
		}
	}
}
