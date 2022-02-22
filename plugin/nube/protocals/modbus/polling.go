package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
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
					var mbClient smod.ModbusClient
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
						mbClient, err = i.setClient(net, true, true)
						if err != nil {
							log.Errorf("modbus: failed to set client %v %s\n", err, *net.SerialPort)
							break
						}
					} else {
						mbClient, err = i.setClient(net, true, false)
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
						if dev.AddressId == 0 {
							log.Errorf("modbus: AddressId=0 is not valid")
							break
						}
						mbClient.RTUClientHandler.SlaveID = byte(dev.AddressId)
						mbClient.DeviceZeroMode = utils.BoolIsNil(dev.ZeroMode)
						if err != nil {
							log.Errorf("modbus: failed to vaildate SetUnitId %v %d\n", err, dev.AddressId)
						}
						for _, pnt := range dev.Points { //POINTS
							dPnt := dev.PollDelayPointsMS
							if dPnt <= 0 {
								dPnt = 100
							}
							write := isWrite(pnt.ObjectType)
							if write && !utils.BoolIsNil(pnt.WriteValueOnceSync) { //IS WRITE
								_, responseValue, err := networkRequest(mbClient, pnt)
								if err != nil {
									_, err = i.pointUpdateErr(pnt.UUID, pnt, err)
									break
								}
								pnt.PresentValue = &responseValue
								_, err = i.pointUpdate(pnt.UUID, pnt)

							} else { //READ
								_, responseValue, err := networkRequest(mbClient, pnt)
								if err != nil {
									_, err = i.pointUpdateErr(pnt.UUID, pnt, err)
									break
								}
								pnt.PresentValue = &responseValue
								_, err = i.pointUpdate(pnt.UUID, pnt)
								if err != nil {
									break
								}
							}
							time.Sleep(dPnt * time.Millisecond)

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
