package main

import (
	"context"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

type polling struct {
	enable        bool
	loopDelay     time.Duration
	delayNetworks time.Duration
	delayDevices  time.Duration
	delayPoints   time.Duration
	isRunning     bool
}

func delays(networkType string) (deviceDelay, pointDelay time.Duration) {
	deviceDelay = 250 * time.Millisecond
	pointDelay = 500 * time.Millisecond
	if networkType == model.TransType.LoRa {
		deviceDelay = 80 * time.Millisecond
		pointDelay = 6000 * time.Millisecond
	}
	return
}

var poll poller.Poller

func (inst *Instance) polling(p polling) error {
	if p.enable {
		poll = poller.New()
	}
	var counter int
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	log.Infoln("init bacnet master network")
	f := func() (bool, error) {
		nets, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, arg)
		if len(nets) == 0 {
			time.Sleep(2 * time.Second)
			log.Info("bacnet-master: NO NETWORKS FOUND")
		}
		for _, net := range nets { //NETWORKS
			if !inst.pollingEnabled {
				//break
			}

			if net.UUID != "" && net.PluginConfId == inst.pluginUUID {
				timeStart := time.Now()
				_, pointDelay := delays(net.TransportType)
				counter++
				log.Infof("bacnet-master-poll: POLL START: NAME: %s\n", net.Name)
				if boolean.IsFalse(net.Enable) {
					log.Infof("bacnet-master: LOOP NETWORK DISABLED: COUNT %v NAME: %s\n", counter, net.Name)
					continue
				}
				for _, dev := range net.Devices { //DEVICES
					if boolean.IsFalse(net.Enable) {
						log.Infof("bacnet-master-device: DEVICE DISABLED: NAME: %s\n", dev.Name)
						continue
					}

					for _, pnt := range dev.Points { //POINTS
						if boolean.IsFalse(net.Enable) {
							continue
						}
						if pnt.WriteMode == "read_only" {
							readFloat32, err := inst.doReadValue(pnt, net.UUID, dev.UUID)
							if err != nil {
								continue
							} else {
								var b = float64(readFloat32)
								_, err := inst.pointUpdateValue(pnt.UUID, b)
								if err != nil {
									//return false, err
								}
								log.Infof("bacnet-master-point: POINT READ:  %f\n", b)
							}
						} else if pnt.WriteMode == "write_only" || pnt.WriteMode == "write_then_read" {
							err := inst.doWrite(pnt, net.UUID, dev.UUID)
							if err != nil {
								log.Errorln("bacnet-master write error", err)
								continue
							} else {

							}
						}

						time.Sleep(pointDelay) //DELAY between points
					}
					timeEnd := time.Now()
					diff := timeEnd.Sub(timeStart)
					out := time.Time{}.Add(diff)
					log.Infof("bacnet-master-poll-loop: NETWORK-NAME:%s POLL-DURATION: %s  POLL-COUNT: %d\n", net.Name, out.Format("15:04:05.000"), counter)
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
