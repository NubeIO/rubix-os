package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
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
	return false, nil
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
	arg.IpConnection = true
	f := func() (bool, error) {
		fmt.Println(counter)
		counter++
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		for cnt, net := range nets { //networks
			if net.UUID != "" {
				for _, dev := range net.Devices { //devices
					var client Client
					var dCheck devCheck
					dCheck.devUUID = dev.UUID
					dCheck.client = client
					client.Host = dev.CommonIP.Host
					client.Port = utils.PortAsString(dev.CommonIP.Port)
					err := setClient(client)
					if err != nil {
						log.Errorf("modbus: failed to set client %v %s\n", err, dev.CommonIP.Host)
					}
					fmt.Println(cnt, dev.UUID, dev.CommonIP.Host, dev.CommonIP.Port, dev.AddressId)
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
							fmt.Println(cnt, pnt.UUID, pnt.Name, pnt.ObjectType)
							dPnt := p.delayPoints
							if !isConnected() {
								fmt.Println("isConnected")
							} else {
								ops.Addr = uint16(pnt.AddressId)
								ops.ObjectType = pnt.ObjectType
								ops.IsHoldingReg = utils.BoolIsNil(pnt.IsOutput)
								ops.ZeroMode = utils.BoolIsNil(pnt.ZeroMode)
								request, err := parseRequest(ops)
								if err != nil {
									fmt.Println(err)
								}
								r, err := DoOperations(cli, request)
								fmt.Println(r)
								//cli.SetEncoding()
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
