package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/poller"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

const defaultInterval = 100 * time.Millisecond

type polling struct {
	enable        bool
	loopDelay     time.Duration
	delayNetworks time.Duration
	delayDevices  time.Duration
	delayPoints   time.Duration
	isRunning     bool
}

type devCheck struct {
	devUUID string
	client  Client
}

func checkDevValid(d devCheck) (bool, error) {
	if d.devUUID == "" {
		log.Errorf("modbus: device id is null \n")
		return false, errors.New("modbus: failed to set client")
	}
	return true, nil
}

func valueRaw(responseRaw interface{}) []byte {
	j, err := json.Marshal(responseRaw)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	return j
}

var poll poller.Poller

func (i *Instance) PollingTCP(p polling) error {
	if p.delayNetworks <= 0 {
		p.delayNetworks = defaultInterval
	}
	if p.delayDevices <= 0 {
		p.delayDevices = defaultInterval
	}
	if p.delayPoints <= 0 {
		p.delayPoints = defaultInterval
	}
	if p.enable {
		poll = poller.New()
	}
	var counter int
	var arg api.Args
	arg.Devices = true
	arg.Points = true
	arg.SerialConnection = true
	arg.IpConnection = true
	f := func() (bool, error) {
		log.Infof("modbus: LOOP COUNT: %v\n", counter)
		counter++
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				for _, dev := range net.Devices { //DEVICES
					var client Client
					var dCheck devCheck
					dCheck.devUUID = dev.UUID
					if net.TransportType == model.TransType.Serial {
						client.SerialPort = net.SerialConnection.SerialPort
						client.BaudRate = net.SerialConnection.BaudRate
						client.DataBits = net.SerialConnection.DataBits
						client.StopBits = net.SerialConnection.StopBits
						err = i.setClient(client, net.UUID, true, true)
						if err != nil {
							log.Errorf("modbus: failed to set client %v %s\n", err, dev.CommonIP.Host)
						}
					} else {
						dCheck.client = client
						client.Host = dev.CommonIP.Host
						client.Port = utils.PortAsString(dev.CommonIP.Port)
						err = i.setClient(client, net.UUID, true, false)
						if err != nil {
							log.Errorf("modbus: failed to set client %v %s\n", err, dev.CommonIP.Host)
						}
					}
					validDev, err := checkDevValid(dCheck)
					if err != nil {
						log.Errorf("modbus: failed to vaildate device %v %s\n", err, dev.CommonIP.Host)
					}
					dNet := p.delayNetworks
					time.Sleep(dNet)
					if validDev {
						cli := getClient()
						err := cli.SetUnitId(uint8(dev.AddressId))
						if err != nil {
							log.Errorf("modbus: failed to vaildate SetUnitId %v %d\n", err, dev.AddressId)
						}
						var ops Operation
						ops.UnitId = uint8(dev.AddressId)
						for _, pnt := range dev.Points { //POINTS
							dPnt := dev.PollDelayPointsMS
							if dPnt <= 0 {
								dPnt = 100
							}
							if !isConnected() {
							} else {
								a := utils.IntIsNil(pnt.AddressId)
								ops.Addr = uint16(a)
								l := utils.IntIsNil(pnt.AddressLength)
								ops.Length = uint16(l)
								ops.ObjectType = pnt.ObjectType
								ops.IsHoldingReg = utils.BoolIsNil(pnt.IsOutput)
								ops.ZeroMode = utils.BoolIsNil(dev.ZeroMode)
								if pnt.Priority != nil {
									if (*pnt.Priority).P16 != nil {
										ops.WriteValue = *pnt.Priority.P16
										log.Infof("modbus: WRITE ObjectType: %s  Addr: %d WriteValue: %v\n", ops.ObjectType, ops.Addr, ops.WriteValue)
									}
								}
								request, err := parseRequest(ops)
								if err != nil {
									log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
								}
								responseRaw, responseValue, err := networkRequest(cli, request)
								var _pnt model.Point
								if isWrite(ops.ObjectType) { //IS WRITE
									_pnt.UUID = pnt.UUID
									_pnt.PresentValue = &ops.WriteValue //update point value
									cov := utils.Float64IsNil(pnt.COV)
									covEvent, _ := utils.COV(ops.WriteValue, utils.Float64IsNil(pnt.ValueOriginal), cov)
									if covEvent {
										log.Infof("MODBUS WRITE COV EVENT")
										i.store.Set(pnt.UUID, _pnt, -1) //store point in cache
										if err != nil {
											log.Errorf("modbus-write-cov: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										} else {
											_pnt.InSync = utils.NewTrue()
											_, err = i.pointUpdate(pnt.UUID, &_pnt)
											log.Infof("modbus-write-cov: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										}
									} else {
										if !utils.BoolIsNil(pnt.InSync) {
											log.Infof("MODBUS WRITE SYNC POINT")
											_pnt.UUID = pnt.UUID
											_pnt.PresentValue = &ops.WriteValue //update point value
											if err != nil {
												log.Errorf("modbus-write-sync: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											} else {
												_pnt.InSync = utils.NewTrue()
												_, err = i.pointUpdate(pnt.UUID, &_pnt)
												log.Infof("modbus-write-sync: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											}
										}
										if counter == 1 {
											log.Infof("MODBUS WRITE SYNC ON START")
											_pnt.UUID = pnt.UUID
											_pnt.PresentValue = &ops.WriteValue //update point value
											if err != nil {
												log.Errorf("modbus-write-start-sync: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											} else {
												_pnt.InSync = utils.NewTrue()
												_, err = i.pointUpdate(pnt.UUID, &_pnt)
												log.Infof("modbus-write-start-sync: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											}
										}
									}
								} else { //READ
									_pnt.UUID = pnt.UUID
									rs := responseValue
									_pnt.PresentValue = &rs //update point value
									_pnt.ValueRaw = valueRaw(responseRaw)
									cov := utils.Float64IsNil(pnt.COV)
									covEvent, _ := utils.COV(ops.WriteValue, utils.Float64IsNil(pnt.ValueOriginal), cov)
									if covEvent {
										_, err = i.pointUpdate(pnt.UUID, &_pnt)
										i.store.Set(pnt.UUID, _pnt, -1) //store point in cache
										if err != nil {
											log.Errorf("modbus: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										} else {
											_pnt.InSync = utils.NewTrue()
											log.Infof("modbus: ObjectType---------: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										}

									} else {
										_pnt.UUID = pnt.UUID
										rs = responseValue
										_pnt.PresentValue = &rs //update point value
										_pnt.ValueRaw = valueRaw(responseRaw)
										_, err = i.pointUpdate(pnt.UUID, &_pnt)
										i.store.Set(pnt.UUID, _pnt, -1) //store point in cache
										if err != nil {
											log.Errorf("modbus: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										} else {
											_pnt.InSync = utils.NewTrue()
											log.Infof("modbus: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										}

									}
								}
								time.Sleep(dPnt * time.Millisecond)
							}
						}
					}
				}
			}
		}
		if !p.enable { //TODO the disable of the polling isn't working
			return true, nil
		} else {
			return false, nil
		}
	}
	err := poll.Poll(context.Background(), f)
	if err != nil {
		return nil
	}
	return nil
}
