package main

import (
	"context"
	"errors"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/poller"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

const defaultInterval = 2000 * time.Millisecond

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
		for _, net := range nets { //networks
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				for _, dev := range net.Devices { //devices
					var client Client
					var dCheck devCheck
					dCheck.devUUID = dev.UUID
					if net.TransportType == model.TransType.Serial {
						client.SerialPort = net.SerialConnection.SerialPort
						client.BaudRate = net.SerialConnection.BaudRate
						client.DataBits = net.SerialConnection.DataBits
						client.StopBits = net.SerialConnection.StopBits
						err = setClientSerial(client)
						if err != nil {
							log.Errorf("modbus: failed to set client %v %s\n", err, dev.CommonIP.Host)
						}
					} else {
						dCheck.client = client
						client.Host = dev.CommonIP.Host
						client.Port = utils.PortAsString(dev.CommonIP.Port)
						err = setClient(client)
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
						var ops Operation
						ops.UnitId = uint8(dev.AddressId)
						for _, pnt := range dev.Points { //points
							dPnt := p.delayPoints
							if !isConnected() {
							} else {
								ops.Addr = uint16(pnt.AddressId)
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
								r, err := DoOperations(cli, request)
								log.Infof("modbus: ObjectType: %s  Addr: %d Response: %v\n", ops.ObjectType, ops.Addr, r)
							}
							time.Sleep(dPnt)
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
