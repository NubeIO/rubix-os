package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	edgerest "github.com/NubeDev/flow-framework/plugin/nube/protocals/edge28/restclient"
	"github.com/NubeDev/flow-framework/src/poller"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
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

var poll poller.Poller

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var getUI *edgerest.UI
var getDI *edgerest.DI

func (i *Instance) polling(p polling) error {
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
	arg.WithSerialConnection = true
	arg.WithIpConnection = true
	f := func() (bool, error) {
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		if len(nets) == 0 {
			time.Sleep(15000 * time.Millisecond)
			log.Info("edge-28: NO NETWORKS FOUND")
		}
		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				log.Infof("edge-28: LOOP COUNT: %v\n", counter)
				counter++
				for _, dev := range net.Devices { //DEVICES
					if err != nil {
						log.Errorf("edge-28: failed to vaildate device %v %s\n", err, dev.CommonIP.Host)
					}
					rest := edgerest.NewNoAuth(dev.CommonIP.Host, utils.ToString(dev.CommonIP.Port))
					getUI, err = rest.GetUIs()
					if err != nil {
						return false, err
					}
					getDI, err = rest.GetDIs()
					dNet := p.delayNetworks
					time.Sleep(dNet)
					for _, pnt := range dev.Points { //POINTS
						switch pnt.IoID {
						case pointList.UI1:
							fmt.Println(getUI.Val.UI1.Val, pointList.UI1, pnt.UUID, pnt.IoID)
						case pointList.UI2:
							fmt.Println(getUI.Val.UI2.Val, pointList.UI2, pnt.UUID, pnt.IoID)
						case pointList.DI1:
							fmt.Println(getDI.Val.DI1.Val, pointList.DI1, pnt.UUID, pnt.IoID)
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
