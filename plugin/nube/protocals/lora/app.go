package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

var err error

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
			p.IoID = e //temp
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
			p.IoID = e //temp
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
			p.IoID = e //temp
			err := i.addPoint(p)
			if err != nil {
				log.Errorf("lora: issue on addPoint: %v\n", err)
				return nil, err
			}
		}
	} else if deviceBody.Model == string(decoder.ME) {
		for _, e := range ME {
			n := fmt.Sprintf("%s_%s_%s", model.TransProtocol.Lora, deviceBody.AddressUUID, e)
			p.Name = n
			p.Description = deviceBody.Model
			p.IoID = e                  //temp
			p.IoType = model.IOType.RAW //raw
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
			err := i.updatePoint(p, code)
			if err != nil {
				log.Errorf("lora: issue on updatePoint: %v\n", err)
				return nil, err
			}
		}
	} else if code == string(decoder.THL) {
		for _, e := range THL {
			p.ThingType = e
			err := i.updatePoint(p, code)
			if err != nil {
				log.Errorf("lora: issue on updatePoint: %v\n", err)
				return nil, err
			}
		}
	} else if code == string(decoder.TH) {
		for _, e := range THL {
			p.ThingType = e
			err := i.updatePoint(p, code)
			if err != nil {
				log.Errorf("lora: issue on updatePoint: %v\n", err)
				return nil, err
			}
		}
	} else if code == string(decoder.ME) {
		for _, e := range ME {
			p.ThingType = e
			err := i.updatePoint(p, code)
			if err != nil {
				log.Errorf("lora: issue on updatePoint: %v\n", err)
				return nil, err
			}
		}
	}
	return nil, nil

}

// updatePoint by its lora id and type as in temp or lux
func (i *Instance) updatePoint(body *model.Point, sensorType string) error {
	addr := body.AddressUUID
	pnt, err := i.db.GetPointByFieldAndIOID("address_uuid", addr, body)
	if err != nil {
		log.Errorf("lora: issue on failed to find point: %v\n", err)
		return err
	}
	if sensorType == string(decoder.ME) {
		*body.PresentValue = decoder.MicroEdgePointType(pnt.IoType, *body.PresentValue)
	}
	_, err = i.db.UpdatePoint(pnt.UUID, body, true)
	log.Infof("lora UpdatePoint: AddressUUID: %s value:%v IoID:%s Type:%s \n", addr, *body.PresentValue, body.IoID, sensorType)
	if err != nil {
		log.Errorf("lora: issue on UpdatePoint: %v\n", err)
		return err
	}

	return nil
}

// updatePoint by its lora id
func (i *Instance) devTHLM(pnt *model.Point, value float64, sensorType string) error {
	pnt.PresentValue = &value
	pnt.CommonFault.InFault = false
	pnt.CommonFault.MessageLevel = model.MessageLevel.Info
	pnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
	pnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	pnt.CommonFault.LastOk = time.Now().UTC()
	err := i.updatePoint(pnt, sensorType)
	if err != nil {
		log.Errorf("lora: issue on update points %v\n", err)
	}
	return nil
}

var ME = []string{"rssi", "voltage", "pulse", "ai1", "ai2", "ai3"}
var TH = []string{"rssi", "voltage", "temperature", "humidity"}
var THL = []string{"rssi", "voltage", "temperature", "humidity", "light"}
var THLM = []string{"rssi", "voltage", "temperature", "humidity", "light", "motion"}
