package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

//wizard make a network/dev/pnt
func (i *Instance) wizardSerial(body wizard) (string, error) {
	sp := "/dev/ttyACM0"
	if body.SerialPort != "" {
		sp = body.SerialPort
	}
	id := "AAB296C4"
	if body.SensorID != "" {
		id = body.SensorID
	}
	st := string(decoder.THLM)
	if body.SensorType != "" {
		st = body.SensorType
	}
	var ser model.SerialConnection
	ser.SerialPort = sp
	ser.BaudRate = 38400

	var net model.Network
	net.Name = model.TransProtocol.Lora
	net.TransportType = model.TransType.Serial
	net.PluginPath = model.TransProtocol.Lora
	net.SerialConnection = &ser

	var dev model.Device
	dev.Name = model.TransProtocol.Lora
	dev.AddressUUID = id
	dev.Manufacture = model.CommonNaming.NubeIO
	dev.Model = st

	var pnt model.Point
	_, err = i.db.WizardNewNetDevPnt("lora", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add lora serial network wizard", err
	}
	return "pass: added network and points", err
}

//listSerialPorts list all serial ports on host
func (i *Instance) listSerialPorts() (*utils.Array, error) {
	ports, err := serial.GetPortsList()
	p := utils.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}

// PublishSensor close serial port
func (i *Instance) publishSensor(commonSensorData decoder.CommonValues, sensorStruct interface{}) {
	pnt := new(model.Point)
	pnt.AddressUUID = commonSensorData.Id
	if commonSensorData.Sensor == string(decoder.THLM) {
		s := sensorStruct.(decoder.TDropletTHLM)
		log.Infof("lora decode as DropletTHLM: AddressUUID: %s Sensor:%s voltage:%v \n", pnt.AddressUUID, commonSensorData.Sensor, float64(s.Voltage))
		for _, e := range THLM {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Voltage:
				f := float64(s.Voltage)
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Temp:
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, s.Temperature)
				if err != nil {
					return
				}
			case model.PointTags.Humidity:
				pnt.UnitType = e //set point type
				f := float64(s.Humidity)
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Light:
				pnt.UnitType = e //set point type
				f := float64(s.Light)
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Motion:
				pnt.UnitType = e //set point type
				f := float64(s.Motion)
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			}
		}
	}
	if commonSensorData.Sensor == string(decoder.THL) {
		s := sensorStruct.(decoder.TDropletTHLM)
		log.Infof("lora decode as DropletTHLM: AddressUUID: %s Sensor:%s voltage:%v \n", pnt.AddressUUID, commonSensorData.Sensor, float64(s.Voltage))
		for _, e := range THL {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Voltage:
				f := float64(s.Voltage)
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Temp:
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, s.Temperature)
				if err != nil {
					return
				}
			case model.PointTags.Humidity:
				pnt.UnitType = e //set point type
				f := float64(s.Humidity)
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Light:
				pnt.UnitType = e //set point type
				f := float64(s.Light)
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			}
		}
	}
	if commonSensorData.Sensor == string(decoder.TH) {
		s := sensorStruct.(decoder.TDropletTHLM)
		log.Infof("lora decode as DropletTHLM: AddressUUID: %s Sensor:%s voltage:%v \n", pnt.AddressUUID, commonSensorData.Sensor, float64(s.Voltage))
		for _, e := range TH {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Voltage:
				f := float64(s.Voltage)
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Temp:
				pnt.UnitType = e //set point type
				err := i.devTHLM(pnt, s.Temperature)
				if err != nil {
					return
				}
			case model.PointTags.Humidity:
				pnt.UnitType = e //set point type
				f := float64(s.Humidity)
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			}
		}
	}
}
