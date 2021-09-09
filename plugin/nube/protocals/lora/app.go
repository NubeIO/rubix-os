package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	log "github.com/sirupsen/logrus"
	"time"
)

/*
user adds a network
user adds a device
- create device and send plugin the uuid
- ask the plugin do you want to add pre-made points for example
- add points
*/

func SerialOpenAndRead() error {
	s := new(SerialSetting)
	s.SerialPort = "/dev/ttyACM1"
	s.BaudRate = 38400
	sc := New(s)
	err := sc.NewSerialConnection()
	if err != nil {
		return err
	}
	sc.Loop()
	return nil
}

// SerialOpen open serial port
func (i *Instance) SerialOpen() error {
	go func() error {
		err := SerialOpenAndRead()
		if err != nil {
			return err
		}
		return nil
	}()
	log.Info("LORA: open serial port")
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
NETWORK
*/
// addPoints close serial port
func (i *Instance) validateNetwork(body *model.Network) error {
	//rules
	// max one network for lora-raw, if user adds a 2nd network it will be put into fault
	// serial port must be set

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
	if deviceBody.Model == string(decoder.THLM) {
		for _, e := range THLM {
			p.PointType = e //temp
			err := i.addPoint(p)
			if err != nil {
				log.Error("LORA: issue on add points", " ", err)
				return nil, err
			}
		}
	}
	return nil, nil

}

// addPoints close serial port
func (i *Instance) addPoint(body *model.Point) error {
	_, err := i.db.CreatePoint(body)
	if err != nil {
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
			p.PointType = e
			err := i.updatePoint(p)
			if err != nil {
				log.Error("LORA: issue on add points", " ", err)
				return nil, err
			}
		}
	}
	return nil, nil

}

// updatePoint by its lora id
func (i *Instance) updatePoint(body *model.Point) error {
	addr := body.AddressUUID
	_, err := i.db.UpdatePointByField("address_uuid", addr, body)
	if err != nil {
		return err
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
		for _, e := range THLM {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.PresentValue = f
				pnt.CommonFault.InFault = false
				pnt.CommonFault.MessageLevel = model.MessageLevel.Info
				pnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
				pnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
				pnt.CommonFault.LastOk = time.Now().UTC()
				err := i.updatePoint(pnt)
				if err != nil {
					fmt.Println("err", err, s.Id)
				}
			case model.PointTags.Voltage:
				f := float64(s.Voltage)
				pnt.PresentValue = f
				pnt.CommonFault.InFault = false
				pnt.CommonFault.MessageLevel = model.MessageLevel.Info
				pnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
				pnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
				pnt.CommonFault.LastOk = time.Now().UTC()
				err := i.updatePoint(pnt)
				if err != nil {
					fmt.Println("err", err, s.Id)
				}
			case model.PointTags.Temp:
				pnt.PresentValue = s.Temperature
				pnt.CommonFault.InFault = false
				pnt.CommonFault.MessageLevel = model.MessageLevel.Info
				pnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
				pnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
				pnt.CommonFault.LastOk = time.Now().UTC()
				err := i.updatePoint(pnt)
				if err != nil {
					fmt.Println("err", err, s.Id)
				}
			case model.PointTags.Humidity:
				f := float64(s.Humidity)
				pnt.PresentValue = f
				pnt.CommonFault.InFault = false
				pnt.CommonFault.MessageLevel = model.MessageLevel.Info
				pnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
				pnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
				pnt.CommonFault.LastOk = time.Now().UTC()
				err := i.updatePoint(pnt)
				if err != nil {
					fmt.Println("err", err, s.Id)
				}

			}
		}
	}
}
