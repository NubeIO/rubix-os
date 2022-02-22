package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uurl"
	"time"

	"github.com/NubeIO/flow-framework/api"
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
					mbClient, err = i.setClient(net, dev, true)
					if err != nil {
						log.Errorf("modbus: failed to set client %v %s\n", err, net.Name)
						break
					}
					if net.TransportType == model.TransType.Serial {
						mbClient.RTUClientHandler.SlaveID = byte(dev.AddressId)
					} else if dev.TransportType == model.TransType.IP {
						url, err := uurl.JoinIpPort(dev.Host, dev.Port)
						if err != nil {
							log.Errorf("modbus: failed to validate device IP %s\n", url)
							break
						}
						mbClient.TCPClientHandler.Address = url
						mbClient.TCPClientHandler.SlaveID = byte(dev.AddressId)
					} else {
						log.Errorf("modbus: failed to validate device and network %v %s\n", err, dev.Name)
						break
					}
					dNet := p.delayNetworks
					time.Sleep(dNet)
					for _, pnt := range dev.Points { //POINTS
						dPnt := dev.PollDelayPointsMS
						if dPnt <= 0 {
							dPnt = 100
						}
						write := isWrite(pnt.ObjectType)
						if write && !utils.BoolIsNil(pnt.WriteValueOnceSync) { //IS WRITE
							fmt.Println("WRITE", *pnt.Priority.P16)
							_, responseValue, err := networkRequest(mbClient, pnt)
							if err != nil {
								_, err = i.pointUpdateErr(pnt.UUID, pnt, err)
								break
							}
							fmt.Println("WRITE responseValue", responseValue)
							_, err = i.pointUpdate(pnt, responseValue)
						} else { //READ
							_, responseValue, err := networkRequest(mbClient, pnt)
							if err != nil {
								_, err = i.pointUpdateErr(pnt.UUID, pnt, err)
								break
							}
							pnt.PresentValue = &responseValue
							_, err = i.pointUpdate(pnt, responseValue)
							if err != nil {
								break
							}
						}
						time.Sleep(dPnt * time.Millisecond)
					}

				}
			}
			time.Sleep(5 * time.Second)
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
