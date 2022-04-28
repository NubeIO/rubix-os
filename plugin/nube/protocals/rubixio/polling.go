package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/nube/rubixio"
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/rest"
	log "github.com/sirupsen/logrus"
	"time"
)

const defaultInterval = 2000 * time.Millisecond //default polling is 2.5 sec
const pollName = "polling"

type polling struct {
	enable        bool
	loopDelay     time.Duration
	delayNetworks time.Duration
	delayDevices  time.Duration
	delayPoints   time.Duration
	isRunning     bool
}

var poll poller.Poller

func (inst *Instance) syncInputs(pnt *model.Point, inputs *rubixio.Inputs) {
	if inputs == nil {
		return
	}
	IoNumber := pnt.IoNumber
	pntType := pnt.IoType //10k

	if IoNumber == inputs.UI1.IoNum {
		temp := inputs.UI1.Temp10K
		_, err := inst.pointUpdateValue(pnt.UUID, temp)
		if err != nil {
			//return
		}

		fmt.Println(pntType, temp)

	}

}

func (inst *Instance) getInputs() *rubixio.Inputs {
	restService := &rest.Service{}
	restService.Url = "192.168.15.194"
	restService.Port = 5001
	restOptions := &rest.Options{}
	restService.Options = restOptions
	restService = rest.New(restService)

	nubeProxy := &rest.NubeProxy{}
	restService.NubeProxy = nubeProxy

	bacnetClient := rubixio.New(restService)
	inputs, err := bacnetClient.GetInputs()
	//inputs

	fmt.Println(inputs, err)
	return inputs

}

func (inst *Instance) polling(p polling) error {
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
	var counter float64
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	f := func() (bool, error) {
		nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		if len(nets) == 0 {
			time.Sleep(15000 * time.Millisecond)
			log.Info("rubixio-polling: NO NETWORKS FOUND")
		}
		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == inst.pluginUUID {
				log.Infof("rubixio-polling: LOOP COUNT: %v\n", counter)
				counter++
				for _, dev := range net.Devices { //DEVICES
					if err != nil {
						log.Errorf("rubixio-polling: failed to vaildate device %v %s\n", err, dev.CommonIP.Host)
					}
					dNet := p.delayNetworks
					time.Sleep(dNet)
					inputs := inst.getInputs()
					for _, pnt := range dev.Points { //POINTS
						inst.syncInputs(pnt, inputs)

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
		return err
	}
	return nil
}
