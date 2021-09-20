package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/poller"
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
		fmt.Println("return false")
		fmt.Println(counter)
		counter++
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		for cnt, net := range nets { //networks
			//fmt.Println(cnt, net)
			for _, dev := range net.Devices { //devices
				//fmt.Println(cnt, dev.UUID, dev.Name)

				for _, pnt := range dev.Points { //points
					fmt.Println(cnt, pnt.UUID, pnt.Name)
					dPnt := p.delayPoints
					time.Sleep(dPnt)
					fmt.Println(cnt, pnt.UUID, pnt.Name)
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
