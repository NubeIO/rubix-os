package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

const defaultInterval = 1000 * time.Millisecond

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
	arg.WithDevices = true
	arg.WithPoints = true
	f := func() (bool, error) {
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if len(nets) == 0 {
			time.Sleep(2 * time.Second)
			log.Info("modbus: NO MODBUS NETWORKS FOUND")
		}
		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				log.Infof("modbus: LOOP COUNT: %v\n", counter)
				counter++
				for _, dev := range net.Devices { //DEVICES
					var client Client
					var dCheck devCheck
					dCheck.devUUID = dev.UUID
					if net.TransportType == model.TransType.Serial {
						if net.SerialPort == nil {
							log.Errorln("invalid serial connection details", "SerialPort")
							break
						}
						if net.SerialBaudRate == nil {
							log.Errorln("invalid serial connection details", "SerialBaudRate")
							break
						}
						if net.SerialDataBits == nil {
							log.Errorln("invalid serial connection details", "SerialDataBits")
							break
						}
						if net.SerialStopBits == nil {
							log.Errorln("invalid serial connection details", "SerialStopBits")
							break
						}
						if net.SerialParity == nil {
							log.Errorln("invalid serial connection details", "SerialParity")
							break
						}
						client.SerialPort = *net.SerialPort
						client.BaudRate = *net.SerialBaudRate
						client.DataBits = *net.SerialDataBits
						client.StopBits = *net.SerialStopBits
						client.Parity = *net.SerialParity
						err = i.setClient(client, net.UUID, true, true)
						if err != nil {
							log.Errorf("modbus: failed to set client %v %s\n", err, dev.CommonIP.Host)
							break
						}
					} else {
						dCheck.client = client
						client.Host = dev.CommonIP.Host
						client.Port = utils.PortAsString(dev.CommonIP.Port)
						err = i.setClient(client, net.UUID, true, false)
						if err != nil {
							log.Errorf("modbus: failed to set client %v %s\n", err, dev.CommonIP.Host)
							break
						}
					}
					validDev, err := checkDevValid(dCheck)
					if err != nil {
						log.Errorf("modbus: failed to vaildate device %v %s\n", err, dev.CommonIP.Host)
						break
					}
					dNet := p.delayNetworks
					time.Sleep(dNet)
					if validDev {
						cli := getClient()
						if dev.AddressId == 0 {
							log.Errorf("modbus: AddressId=0 is not valid")
							break
						}
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
								a := utils.IntIsNil(pnt.AddressID)
								ops.Addr = uint16(a)
								l := utils.IntIsNil(pnt.AddressLength)
								ops.Length = uint16(l)
								ops.ObjectType = pnt.ObjectType
								ops.Encoding = pnt.ObjectEncoding
								ops.IsHoldingReg = utils.BoolIsNil(pnt.IsOutput)
								ops.ZeroMode = utils.BoolIsNil(dev.ZeroMode)
								_isWrite := isWrite(ops.ObjectType)
								var _pnt model.Point
								if _isWrite && !utils.BoolIsNil(pnt.WriteValueOnceSync) || counter == 1 { //IS WRITE
									if pnt.Priority != nil {
										if (*pnt.Priority).P16 != nil {
											ops.WriteValue = *pnt.Priority.P16
											log.Infof("modbus-write: WRITE ObjectType: %s  Addr: %d WriteValue: %v\n", ops.ObjectType, ops.Addr, ops.WriteValue)
										}
									}
									request, err := parseRequest(ops)
									if err != nil {
										log.Errorf("modbus-write: failed to read holding/input registers: %v\n", err)
										break
									}
									responseRaw, responseValue, err := networkRequest(cli, request)
									log.Infof("modbus-write: ObjectType: %s  Addr: %d ARRAY: %v\n", ops.ObjectType, ops.Addr, responseRaw)
									_pnt.UUID = pnt.UUID
									_pnt.PresentValue = &ops.WriteValue //update point value
									cov := utils.Float64IsNil(pnt.COV)
									covEvent, _ := utils.COV(ops.WriteValue, utils.Float64IsNil(pnt.OriginalValue), cov)
									if covEvent {
										if err != nil {
											log.Errorf("modbus-write-cov: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											_, err = i.pointUpdateErr(pnt.UUID, &_pnt, err)
											break
										} else {
											log.Infof("modbus-write: MODBUS WRITE COV EVENT: COV value is %v\n", cov)
											_pnt.InSync = utils.NewTrue()
											if utils.BoolIsNil(pnt.WriteValueOnce) {
												_pnt.WriteValueOnceSync = utils.NewTrue()
											}
											_, err = i.pointUpdate(pnt.UUID, &_pnt)
											log.Infof("modbus-write-cov: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										}
									} else {
										if !utils.BoolIsNil(pnt.InSync) {
											log.Infof("modbus-write: MODBUS WRITE SYNC POINT")
											_pnt.UUID = pnt.UUID
											_pnt.PresentValue = &ops.WriteValue //update point value
											if err != nil {
												log.Errorf("modbus-write-sync: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
												break
											} else {
												_pnt.InSync = utils.NewTrue()
												if utils.BoolIsNil(pnt.WriteValueOnce) {
													_pnt.WriteValueOnceSync = utils.NewTrue()
												}
												_, err = i.pointUpdate(pnt.UUID, &_pnt)
												log.Infof("modbus-write-sync: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											}
										}
										if counter == 1 {
											log.Infof("modbus: MODBUS WRITE SYNC ON START")
											_pnt.UUID = pnt.UUID
											_pnt.PresentValue = &ops.WriteValue //update point value
											if err != nil {
												log.Errorf("modbus-write-start-sync: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
												break
											} else {
												_pnt.InSync = utils.NewTrue()
												if utils.BoolIsNil(pnt.WriteValueOnce) {
													_pnt.WriteValueOnceSync = utils.NewTrue()
												}
												_, err = i.pointUpdate(pnt.UUID, &_pnt)
												log.Infof("modbus-write-start-sync: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											}
										}
									}
								} else if !_isWrite { //READ
									request, err := parseRequest(ops)
									if err != nil {
										log.Errorf("modbus-read: failed to read holding/input registers: %v\n", err)
										break
									}
									_, responseValue, err := networkRequest(cli, request)
									if err != nil {
										log.Errorf("modbus-read: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										_, err = i.pointUpdateErr(pnt.UUID, &_pnt, err)
										break
									}
									_pnt.UUID = pnt.UUID
									rs := responseValue
									_pnt.PresentValue = &rs //update point value
									cov := utils.Float64IsNil(pnt.COV)
									covEvent, _ := utils.COV(ops.WriteValue, utils.Float64IsNil(pnt.OriginalValue), cov)
									if covEvent {
										_, err = i.pointUpdate(pnt.UUID, &_pnt)
										i.store.Set(pnt.UUID, _pnt, -1) //store point in cache
										if err != nil {
											log.Errorf("modbus-read: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											break
										} else {
											_pnt.InSync = utils.NewTrue()
											log.Infof("modbus-read: ObjectType---------: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)

										}
									} else {
										_pnt.UUID = pnt.UUID
										rs = responseValue
										_pnt.PresentValue = &rs //update point value
										//_pnt.ValueRaw = valueRaw(responseRaw)
										_, err = i.pointUpdate(pnt.UUID, &_pnt)
										i.store.Set(pnt.UUID, _pnt, -1) //store point in cache
										if err != nil {
											log.Errorf("modbus:-read ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
											break
										} else {
											_pnt.InSync = utils.NewTrue()
											log.Infof("modbus-read: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, responseValue)
										}
									}
								}
								time.Sleep(dPnt * time.Millisecond)
							}
						}
					}
				}
			}
			time.Sleep(1 * time.Second)
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
