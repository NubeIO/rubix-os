package main

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeIO/flow-framework/utils"
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
	var net model.Network
	net.Name = model.TransProtocol.Lora
	net.TransportType = model.TransType.Serial
	net.PluginPath = model.TransProtocol.Lora
	net.SerialPort = utils.NewStringAddress(sp)
	net.SerialBaudRate = utils.NewUint(38400)

	var dev model.Device
	dev.Name = model.TransProtocol.Lora
	dev.AddressUUID = id
	dev.Manufacture = model.CommonNaming.NubeIO
	dev.Model = st

	var pnt model.Point
	_, err = i.db.WizardNewNetworkDevicePoint("lora", &net, &dev, &pnt)
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
func B2i(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// PublishSensor close serial port
func (i *Instance) publishSensor(commonSensorData decoder.CommonValues, sensorStruct interface{}) {
	pnt := new(model.Point)
	pnt.AddressUUID = commonSensorData.Id
	if commonSensorData.Sensor == string(decoder.THLM) {
		s := sensorStruct.(decoder.TDropletTHLM)
		log.Infof("lora decode as THLM: AddressUUID: %s Sensor:%s voltage:%v \n", pnt.AddressUUID, commonSensorData.Sensor, s.Voltage)
		for _, e := range THLM {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.THLM))
				if err != nil {
					return
				}
			case model.PointTags.Voltage:
				f := s.Voltage
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.THLM))
				if err != nil {
					return
				}
			case model.PointTags.Temp:
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, s.Temperature, string(decoder.THLM))
				if err != nil {
					return
				}
			case model.PointTags.Humidity:
				pnt.IoID = e //set point type
				f := float64(s.Humidity)
				err := i.devTHLM(pnt, f, string(decoder.THLM))
				if err != nil {
					return
				}
			case model.PointTags.Light:
				pnt.IoID = e //set point type
				f := float64(s.Light)
				err := i.devTHLM(pnt, f, string(decoder.THLM))
				if err != nil {
					return
				}
			case model.PointTags.Motion:
				pnt.IoID = e //set point type
				f := B2i(s.Motion)
				err := i.devTHLM(pnt, f, string(decoder.THLM))
				if err != nil {
					return
				}
			}
		}
	}
	if commonSensorData.Sensor == string(decoder.THL) {
		s := sensorStruct.(decoder.TDropletTHLM)
		log.Infof("lora decode as THL: AddressUUID: %s Sensor:%s voltage:%v \n", pnt.AddressUUID, commonSensorData.Sensor, s.Voltage)
		for _, e := range THL {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.THL))
				if err != nil {
					return
				}
			case model.PointTags.Voltage:
				f := s.Voltage
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.THL))
				if err != nil {
					return
				}
			case model.PointTags.Temp:
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, s.Temperature, string(decoder.THL))
				if err != nil {
					return
				}
			case model.PointTags.Humidity:
				pnt.IoID = e //set point type
				f := float64(s.Humidity)
				err := i.devTHLM(pnt, f, string(decoder.THL))
				if err != nil {
					return
				}
			case model.PointTags.Light:
				pnt.IoID = e //set point type
				f := float64(s.Light)
				err := i.devTHLM(pnt, f, string(decoder.THL))
				if err != nil {
					return
				}
			}
		}
	}
	if commonSensorData.Sensor == string(decoder.TH) {
		s := sensorStruct.(decoder.TDropletTHLM)
		log.Infof("lora decode as TH: AddressUUID: %s Sensor:%s voltage:%v \n", pnt.AddressUUID, commonSensorData.Sensor, s.Voltage)
		for _, e := range TH {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.TH))
				if err != nil {
					return
				}
			case model.PointTags.Voltage:
				f := s.Voltage
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.TH))
				if err != nil {
					return
				}
			case model.PointTags.Temp:
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, s.Temperature, string(decoder.TH))
				if err != nil {
					return
				}
			case model.PointTags.Humidity:
				pnt.IoID = e //set point type
				f := float64(s.Humidity)
				err := i.devTHLM(pnt, f, string(decoder.TH))
				if err != nil {
					return
				}
			}
		}
	}
	if commonSensorData.Sensor == string(decoder.ME) {
		s := sensorStruct.(decoder.TMicroEdge)
		log.Infof("lora decode as ME: AddressUUID: %s Sensor:%s voltage:%v \n", pnt.AddressUUID, commonSensorData.Sensor, s.Voltage)
		for _, e := range ME {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.ME))
				if err != nil {
					return
				}
			case model.PointTags.Voltage:
				f := s.Voltage
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.ME))
				if err != nil {
					return
				}
			case model.PointTags.Pulse:
				f := s.Pulse
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, float64(f), string(decoder.ME))
				if err != nil {
					return
				}
			case model.PointTags.AI1:
				f := s.AI1
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.ME))
				if err != nil {
					return
				}
			case model.PointTags.AI2:
				f := s.AI2
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.ME))
				if err != nil {
					return
				}
			case model.PointTags.AI3:
				f := s.AI3
				pnt.IoID = e //set point type
				err := i.devTHLM(pnt, f, string(decoder.ME))
				if err != nil {
					return
				}
			}
		}
	}
}
