package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
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
					devDelay, pointDelay := delays(net.TransportType)
					// counter++
					for _, dev := range net.Devices { // DEVICES
						time.Sleep(devDelay) // DELAY between devices
						dev, err = inst.db.GetDevice(dev.UUID, api.Args{WithPoints: true})
						if err != nil {
							inst.bacnetErrorMsg("BACnetServerPolling(): Device not found")
							continue
						}
						if boolean.IsFalse(dev.Enable) {
							inst.bacnetDebugMsg("DEVICE DISABLED: NAME: ", dev.Name)
							continue
						}
						for _, pnt := range dev.Points { // POINTS
							time.Sleep(pointDelay) // DELAY between points
							// pnt, err = inst.db.GetPoint(pnt.UUID, api.Args{WithPriority: true})
							pnt, err = inst.db.GetPoint(pnt.UUID, api.Args{WithPriority: true})
							if err != nil {
								inst.bacnetErrorMsg("BACnetServerPolling(): Point not found")
								continue
							}
							inst.bacnetDebugMsg("BACnetServerPolling() pnt.ObjectType: ", pnt.ObjectType)
							inst.bacnetDebugMsg("BACnetServerPolling(): pnt.WritePollRequired: ", boolean.IsTrue(pnt.WritePollRequired))
							if boolean.IsFalse(pnt.Enable) {
								continue
							}
							if pnt.Priority == nil {
								inst.bacnetErrorMsg("BACnetServerPolling(): Point doesn't have a priority array; skipped")
								continue
							}

							currentPriorityMap := priorityarray.ConvertToMap(*pnt.Priority)

							if !isWriteableObjectType(pnt.ObjectType) || boolean.IsTrue(pnt.WritePollRequired) {
								// For these points we don't need to read the priority array because we are forcing the values to match our FF point
								for key, val := range currentPriorityMap {
									inst.bacnetDebugMsg("BACnetServerPolling() currentPriorityMap key: ", key, "val", val)
									err := inst.doWrite(pnt, net.UUID, dev.UUID, val)
									if err != nil {
										err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
										continue
										errors.New("hello")
									}
								}

//PointReleasePriority use this when releasing a priority
func (device *Device) PointReleasePriority(pnt *Point, pri uint8) error {
	if pnt == nil{
		return errors.New("invalid point to PointReleasePriority()")
	}
	if pri > 16 || pri < 1 {
		return errors.New("invalid priority to PointReleasePriority()")
	}
	write := &Write{
		ObjectID:      pnt.ObjectID,
		ObjectType:    pnt.ObjectType,
		Prop:          btypes.PropPresentValue,
		WriteNull:     true,
		WritePriority: pri,
	}
	err := device.Write(write)
	if err != nil {
		return err
	}
	return nil
}

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
