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

func delays() (deviceDelay, pointDelay time.Duration) {
	deviceDelay = 100 * time.Millisecond
	pointDelay = 100 * time.Millisecond
	return
}

var poll poller.Poller
var lastPingFailed = "start"
var rsyncWrite = 0

func (inst *Instance) BACnetServerPolling() error {
	poll = poller.New()
	var counter = 0
	f := func() (bool, error) {
		counter++
		// fmt.Println("\n \n")
		inst.bacnetDebugMsg("LOOP COUNT: ", counter)
		var err error
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
					devDelay, pointDelay := delays()
					// counter++
					for _, dev := range net.Devices { // DEVICES
						dev, err = inst.db.GetDevice(dev.UUID, api.Args{WithPoints: true})
						if err != nil {
							inst.bacnetErrorMsg("BACnetServerPolling(): Device not found")
							continue
						}
						if boolean.IsFalse(net.Enable) {
							inst.bacnetDebugMsg("DEVICE DISABLED: NAME: ", dev.Name)
							continue
						}
						err = inst.pingDevice(net, dev)
						if err != nil {
							lastPingFailed = "fail"
						}
						if lastPingFailed == "fail" && err != nil {
							lastPingFailed = "rsync"
						}
						if rsyncWrite == 0 || lastPingFailed == "rsync" {
							inst.massUpdateServer(net, dev)
							lastPingFailed = "in-sync"
						}
						time.Sleep(devDelay)             // DELAY between devices
						for _, pnt := range dev.Points { // POINTS
							// pnt, err = inst.db.GetPoint(pnt.UUID, api.Args{WithPriority: true})
							pnt, err = inst.db.GetPoint(pnt.UUID, api.Args{})
							if err != nil {
								inst.bacnetErrorMsg("BACnetServerPolling(): Point not found")
								continue
							}
							inst.bacnetDebugMsg("BACnetServerPolling() pnt.ObjectType: ", pnt.ObjectType)
							inst.bacnetDebugMsg("BACnetServerPolling(): pnt.WritePollRequired: ", boolean.IsTrue(pnt.WritePollRequired))
							if boolean.IsFalse(net.Enable) {
								continue
							}
							time.Sleep(pointDelay) // DELAY between points
							if !isWriteableObjectType(pnt.ObjectType) || boolean.IsTrue(pnt.WritePollRequired) {
								writeVal := float.NonNil(pnt.WriteValue)
								err := inst.doWrite(pnt, net.UUID, dev.UUID, writeVal)
								if err != nil {
									err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
									continue
								}
								// pnt, err = inst.pointUpdate(pnt, writeVal, true, true)
								// if err != nil {
								//	continue
								// }
							}
							// TODO: below could be optimized by not doing read on successful write.  currently I found cases where the write didn't return an error, but the values wasn't updated on the server.
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
						rsyncWrite = counter % 50

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
