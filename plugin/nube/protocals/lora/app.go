package main

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeDev/flow-framework/utils"
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

var err error

// SerialOpen open serial port
func (i *Instance) SerialOpen() error {
	s := new(SerialSetting)
	var arg api.Args
	arg.WithSerialConnection = true
	net, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if err != nil {
		return err
	}
	if net.SerialConnection == nil {
		return err
	}
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

	checkDeviceNotExist, err := i.db.GetDeviceByField("address_uuid", deviceBody.AddressUUID, false)
	if err != nil {
		return nil, err
	}
	if checkDeviceNotExist.UUID != "" {
		log.Errorf("lora: a device with the same lora ID (address_uuid) exists: %v\n", checkDeviceNotExist.UUID)
		return nil, errors.New("a device with the same lora ID (address_uuid) exists")
	}

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
	} else if deviceBody.Model == string(decoder.THL) {
		for _, e := range THL {
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
	} else if deviceBody.Model == string(decoder.TH) {
		for _, e := range TH {
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
	} else if code == string(decoder.THL) {
		for _, e := range THL {
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
	log.Infof("lora UpdatePointByFieldAndUnit: AddressUUID: %s value:%v UnitType:%s \n", addr, *body.PresentValue, body.UnitType)
	if err != nil {
		log.Errorf("lora: issue on UpdatePointByFieldAndType: %v\n", err)
		return err
	}
	return nil
}

// updatePoint by its lora id
func (i *Instance) devTHLM(pnt *model.Point, value float64) error {
	pnt.PresentValue = &value
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

var TH = []string{"rssi", "voltage", "temperature", "humidity"}
var THL = []string{"rssi", "voltage", "temperature", "humidity", "light"}
var THLM = []string{"rssi", "voltage", "temperature", "humidity", "light", "motion"}
