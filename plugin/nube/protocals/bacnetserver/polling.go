package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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

func (inst *Instance) BACnetServerPolling() error {
	poll = poller.New()
	var counter = 0
	f := func() (bool, error) {
		counter++
		// fmt.Println("\n \n")
		inst.bacnetDebugMsg("LOOP COUNT: ", counter)
		var arg api.Args
		arg.WithDevices = true
		arg.WithPoints = true
		nets, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, arg)
		if len(nets) == 0 {
			time.Sleep(5 * time.Second)
			inst.bacnetDebugMsg("NO NETWORKS FOUND")
		} else {
			inst.bacnetDebugMsg("NETWORKS FOUND", len(nets))
		}
		for _, net := range nets { // NETWORKS
			if boolean.IsFalse(net.Enable) {
				inst.bacnetDebugMsg("NETWORK DISABLED: NAME: ", net.Name)
				continue
			} else {
				if net.UUID != "" && net.PluginConfId == inst.pluginUUID {
					timeStart := time.Now()
					devDelay, pointDelay := delays(net.TransportType)
					// counter++
					for _, dev := range net.Devices { // DEVICES
						if boolean.IsFalse(net.Enable) {
							inst.bacnetDebugMsg("DEVICE DISABLED: NAME: ", dev.Name)
							continue
						}
						time.Sleep(devDelay)             // DELAY between devices
						for _, pnt := range dev.Points { // POINTS
							if boolean.IsFalse(net.Enable) {
								continue
							}
							time.Sleep(pointDelay)            // DELAY between points
							if pnt.WriteMode == "read_only" { // Only need to write value from FF Point to BACnet Server
								pnt.PointPriorityArrayMode = model.PriorityArrayToPresentValue
								writeVal := float.NonNil(pnt.PresentValue)
								err := inst.doWrite(pnt, net.UUID, dev.UUID, writeVal)
								if err != nil {
									err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
									continue
								}
								// val := float.NonNil(pnt.WriteValue) //TODO not sure if this should then update the PV of the point
								pnt, err = inst.pointUpdate(pnt, writeVal, true, true)
								if err != nil {
									continue
								}

							} else if pnt.WriteMode == "write_once_then_read" {
								if boolean.IsTrue(pnt.WritePollRequired) { // DO WRITE IF REQUIRED
									pnt.PointPriorityArrayMode = model.PriorityArrayToWriteValue
									if pnt.WriteValue != nil {
										writeVal := float.NonNil(pnt.WriteValue)
										err := inst.doWrite(pnt, net.UUID, dev.UUID, writeVal)
										if err != nil {
											err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
											continue
										}
										pnt, err = inst.pointUpdate(pnt, writeVal, true, true)
										if err != nil {
											continue
										}
									}
								}
								readFloat, err := inst.doReadValue(pnt, net.UUID, dev.UUID)
								if err != nil {
									err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
									continue
								} else {
									pnt, err = inst.pointUpdate(pnt, readFloat, true, true)
									if err != nil {
										continue
									}
								}
							}
						}
						timeEnd := time.Now()
						diff := timeEnd.Sub(timeStart)
						out := time.Time{}.Add(diff)
						inst.bacnetDebugMsg(fmt.Sprintf("poll-loop: NETWORK-NAME:%s POLL-DURATION: %s  POLL-COUNT: %d\n", net.Name, out.Format("15:04:05.000"), counter))
					}
				}
			}
		}
		return false, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	inst.pollingCancel = cancel
	go poll.GoPoll(ctx, f)
	return nil
}
