package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/poller"
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
	a := poller.New()
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
			//fmt.Println(cnt, net)
			if net.UUID != "" {
				for _, dev := range net.Devices { //devices
					var client Client
					client.Host = dev.CommonIP.Host
					client.Port = utils.PortAsString(dev.CommonIP.Port)
					err := setClient(client)
					if err != nil {
						log.Info(err, "ERROR ON set modbus client")
					}

					fmt.Println(cnt, dev.UUID, dev.CommonIP.Host, dev.CommonIP.Port, dev.AddressId)
					if dev.UUID != "" {
						for _, pnt := range dev.Points { //points
							fmt.Println(cnt, pnt.UUID, pnt.Name, pnt.ObjectType)
							dPnt := p.delayPoints
							cli := getClient()
							if !isConnected() {
								fmt.Println("isConnected")
							} else {
								var ops Operation
								ops.UnitId = 1
								ops.Addr = 1
								ops.IsHoldingReg = true
								ops.ObjectType = pnt.ObjectType
								ops.ZeroMode = true
								request, err := parseRequest(ops)
								if err != nil {
									fmt.Println(err)
								}
								r, err := DoOperations(cli, request)
								fmt.Println(r)
								//cli.SetEncoding()
							}
							time.Sleep(dPnt)
							fmt.Println(cnt, pnt.UUID, pnt.Name)
						}
					}
				}
			}

		}

		return false, nil
	}
	err := a.Poll(context.Background(), f)
	if err != nil {
		return nil
	}
	return nil
}
