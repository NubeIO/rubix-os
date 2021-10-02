package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"time"
)

/*
user adds a network
user adds a device
- create device and send plugin the uuid
- ask the plugin do you want to add pre-made points for example
- add points
*/

var err error

// SerialOpen open serial port
func (i *Instance) SerialOpen() error {
	s := new(SerialSetting)
	var arg api.Args
	arg.SerialConnection = true
	net, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if err != nil {
		return err
	}
	if net.SerialConnection == nil {
		return err
	}
	fmt.Println(i.pluginUUID, net)
	s.SerialPort = net.SerialConnection.SerialPort
	s.BaudRate = int(net.SerialConnection.BaudRate)
	connected := false
	go func() error {
		sc := New(s)
		connected, err = sc.NewSerialConnection()
		if err != nil {
			log.Errorf("lora: issue on SerialOpenAndRead: %v\n", err)
		}
		sc.Loop()
		return nil
	}()
	return nil

}

// SerialClose close serial port
func (i *Instance) SerialClose() error {
	err := Disconnect()
	if err != nil {
		return err
	}
	return nil
}

/*
POINTS
*/

// addPoints close serial port
func (i *Instance) addPoints(deviceBody *model.Device) (*model.Point, error) {
	p := new(model.Point)
	p.DeviceUUID = deviceBody.UUID
	p.AddressUUID = deviceBody.AddressUUID
	p.IsProducer = utils.NewFalse()
	p.IsConsumer = utils.NewFalse()
	p.IsOutput = utils.NewFalse()
	if deviceBody.Model == string(decoder.THLM) {
		for _, e := range THLM {
			n := fmt.Sprintf("%s_%s_%s", model.TransProtocol.Lora, deviceBody.AddressUUID, e)
			p.Name = n
			p.Description = deviceBody.Model
			p.UnitType = e //temp
			err := i.addPoint(p)
			if err != nil {
				log.Errorf("lora: issue on addPoint: %v\n", err)
				return nil, err
			}
		}
	}
	return nil, nil
}

// addPoints add a pnt
func (i *Instance) addPoint(body *model.Point) error {
	_, err := i.db.CreatePoint(body)
	if err != nil {
		log.Errorf("lora: issue on CreatePoint: %v\n", err)
		return err
	}
	return nil
}

// updatePoints update the point values
func (i *Instance) updatePoints(deviceBody *model.Device) (*model.Point, error) {
	p := new(model.Point)
	p.UUID = deviceBody.UUID
	code := deviceBody.AddressUUID
	if code == string(decoder.THLM) {
		for _, e := range THLM {
			p.ThingType = e
			err := i.updatePoint(p)
			if err != nil {
				log.Errorf("lora: issue on updatePoint: %v\n", err)
				return nil, err
			}
		}
	}
	return nil, nil

}

// updatePoint by its lora id and type as in temp or lux
func (i *Instance) updatePoint(body *model.Point) error {
	addr := body.AddressUUID
	_, err := i.db.UpdatePointByFieldAndUnit("address_uuid", addr, body)
	log.Infof("lora UpdatePointByFieldAndUnit: AddressUUID: %s value:%v UnitType:%s \n", addr, body.PresentValue, body.UnitType)
	if err != nil {
		log.Errorf("lora: issue on UpdatePointByFieldAndType: %v\n", err)
		return err
	}
	return nil
}

// updatePoint by its lora id
func (i *Instance) devTHLM(pnt *model.Point, value float64) error {
	pnt.PresentValue = value
	pnt.CommonFault.InFault = false
	pnt.CommonFault.MessageLevel = model.MessageLevel.Info
	pnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
	pnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	pnt.CommonFault.LastOk = time.Now().UTC()
	err := i.updatePoint(pnt)
	if err != nil {
		log.Errorf("lora: issue on update points %v\n", err)
	}
	return nil
}

var THLM = []string{"rssi", "voltage", "temperature", "humidity", "light", "motion"}

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

			}
		}
	}
}

//wizard make a network/dev/pnt
func (i *Instance) wizardSerial() (string, error) {
	var ser model.SerialConnection
	ser.SerialPort = "/dev/ttyACM0"
	ser.BaudRate = 38400

	var net model.Network
	net.Name = model.TransProtocol.Lora
	net.TransportType = model.TransType.Serial
	net.PluginPath = model.TransProtocol.Lora
	net.SerialConnection = &ser

	var dev model.Device
	dev.Name = model.TransProtocol.Lora
	dev.AddressUUID = "AAB296C4"
	dev.Manufacture = model.CommonNaming.NubeIO
	dev.Model = string(decoder.THLM)

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
