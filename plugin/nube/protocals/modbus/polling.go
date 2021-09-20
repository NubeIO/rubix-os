package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"time"
)

type polling struct {
	enable         bool
	loopDelay      time.Duration
	delayAfterLoop time.Duration
	isRunning      bool
}

func (i *Instance) PollingTCP(p polling) error {
	if p.loopDelay <= 0 {
		p.loopDelay = 100
	}
	if p.delayAfterLoop <= 0 {
		p.delayAfterLoop = 100
	}

	c := time.Tick(p.loopDelay * time.Millisecond)
	i.pollingEnabled = true
	for range c {
		var arg api.Args
		arg.Devices = true
		arg.Points = true
		arg.IpConnection = true
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if err != nil {
			return err
		}
		for cnt, net := range nets { //networks
			//fmt.Println(cnt, net)
			for _, dev := range net.Devices { //devices
				//fmt.Println(cnt, dev.UUID, dev.Name)
				for _, pnt := range dev.Points { //points
					fmt.Println(cnt, pnt.UUID, pnt.Name)
				}
			}
		}
		fmt.Printf("Modbus Polling Loop %s\n", time.Now())
		jitter := p.delayAfterLoop * time.Millisecond
		time.Sleep(jitter)
	}

	return nil
}

func (i *Instance) PollingRTU(p polling) error {
	if p.loopDelay <= 0 {
		p.loopDelay = 100
	}
	if p.delayAfterLoop <= 0 {
		p.delayAfterLoop = 100
	}

	c := time.Tick(p.loopDelay * time.Millisecond)
	i.pollingEnabled = true
	for range c {
		var arg api.Args
		arg.Devices = true
		arg.Points = true
		arg.SerialConnection = true
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if err != nil {
			return err
		}
		for cnt, net := range nets { //networks
			//fmt.Println(cnt, net)
			for _, dev := range net.Devices { //devices
				//fmt.Println(cnt, dev.UUID, dev.Name)
				for _, pnt := range dev.Points { //points
					fmt.Println(cnt, pnt.UUID, pnt.Name)
				}
			}

		}
		fmt.Printf("Modbus Polling Loop %s\n", time.Now())
		jitter := p.delayAfterLoop * time.Millisecond
		time.Sleep(jitter)
	}

	return nil
}
