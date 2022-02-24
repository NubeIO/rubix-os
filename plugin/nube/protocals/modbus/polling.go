package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uurl"
	log "github.com/sirupsen/logrus"
	"time"
)

const defaultInterval = 100 * time.Millisecond

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

func delays(networkType string) (networkDelay, deviceDelay, pointDelay time.Duration) {
	networkDelay = 100 * time.Millisecond
	deviceDelay = 100 * time.Millisecond
	pointDelay = 100 * time.Millisecond
	if networkType == model.TransType.LoRa {
		networkDelay = 100 * time.Millisecond
		deviceDelay = 100 * time.Millisecond
		pointDelay = 6000 * time.Millisecond
	}
	return
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
		networkDelay, _, _ := delays("")
		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				_, deviceDelay, pointDelay := delays(net.TransportType)
				counter++
				if !utils.BoolIsNil(net.Enable) {
					log.Infof("modbus: LOOP NETWORK DISABLED: COUNT %v NAME: %s\n", counter, net.Name)
					continue
				} else {
					log.Infof("modbus: LOOP COUNT: %v\n", counter)
				}
				for _, dev := range net.Devices { //DEVICES
					if !utils.BoolIsNil(dev.Enable) {
						log.Infof("modbus-device: DEVICE DISABLED: NAME: %s\n", dev.Name)
						continue
					}
					var mbClient smod.ModbusClient
					var dCheck devCheck
					dCheck.devUUID = dev.UUID
					mbClient, err = i.setClient(net, dev, true)
					if err != nil {
						log.Errorf("modbus: failed to set client error: %v network name:%s\n", err, net.Name)
						continue
					}
					//log.Infof("modbus-device: DEVICE DISABLED: NAME: %s\n", mbClient.RTUClientHandler.Address)
					if net.TransportType == model.TransType.Serial || net.TransportType == model.TransType.LoRa {
						if dev.AddressId >= 1 {
							mbClient.RTUClientHandler.SlaveID = byte(dev.AddressId)
						}
					} else if dev.TransportType == model.TransType.IP {
						url, err := uurl.JoinIpPort(dev.Host, dev.Port)
						if err != nil {
							log.Errorf("modbus: failed to validate device IP %s\n", url)
							continue
						}
						mbClient.TCPClientHandler.Address = url
						mbClient.TCPClientHandler.SlaveID = byte(dev.AddressId)
					} else {
						log.Errorf("modbus: failed to validate device and network %v %s\n", err, dev.Name)
						continue
					}
					time.Sleep(deviceDelay)          //DELAY between devices
					for _, pnt := range dev.Points { //POINTS
						if !utils.BoolIsNil(pnt.Enable) {
							log.Infof("modbus-point: POINT DISABLED: NAME: %s\n", pnt.Name)
							continue
						}
						write := isWrite(pnt.ObjectType)
						if write { //IS WRITE
							//get existing
							if !utils.BoolIsNil(pnt.InSync) {
								fmt.Println("WRITE", pnt.Name)
								_, responseValue, err := networkRequest(mbClient, pnt, true)
								if err != nil {
									_, err = i.pointUpdateErr(pnt.UUID, err)
									continue
								}
								_, err = i.pointUpdate(pnt.UUID, responseValue)
							}
						} else { //READ
							_, responseValue, err := networkRequest(mbClient, pnt, false)
							if err != nil {
								_, err = i.pointUpdateErr(pnt.UUID, err)
								continue
							}
							//simple cov
							fmt.Println(pnt.Name, *pnt.PresentValue, responseValue)
							isChange := !utils.CompareFloatPtr(pnt.PresentValue, &responseValue)
							if isChange {
								_, err = i.pointUpdate(pnt.UUID, responseValue)
								if err != nil {
									continue
								}
							}
						}
						time.Sleep(pointDelay) //DELAY between points
					}
				}
			}
		}
		time.Sleep(networkDelay) //DELAY between networks
		if !p.enable {           //TODO the disable of the polling isn't working
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
