package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/bacnet/btypes/priority"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"strconv"
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
	// pointDelay = 100 * time.Millisecond
	pointDelay = 5000 * time.Millisecond
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
						for _, pnt := range dev.Points { // POINTS
							time.Sleep(pointDelay) // DELAY between points
							pnt, err = inst.db.GetPoint(pnt.UUID, api.Args{WithPriority: true})
							if err != nil {
								inst.bacnetErrorMsg("BACnetServerPolling(): Point not found")
								continue // next point
							}
							inst.bacnetDebugMsg("BACnetServerPolling() pnt.ObjectType: ", pnt.ObjectType)
							inst.bacnetDebugMsg("BACnetServerPolling(): pnt.WritePollRequired: ", boolean.IsTrue(pnt.WritePollRequired))
							if boolean.IsFalse(pnt.Enable) {
								inst.bacnetErrorMsg("BACnetServerPolling(): Point is disabled; skipped")
								continue // next point
							}
							if pnt.Priority == nil {
								inst.bacnetErrorMsg("BACnetServerPolling(): Point doesn't have a priority array; skipped")
								err = inst.pointUpdateErr(pnt, "no priority array found", model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
								continue // next point
							}
							if !isWriteableObjectType(pnt.ObjectType) { // Do these actions fort AI, BI
								// For AI/BI there is no priority array we just write the values to match our FF point WriteValue
								writeVal := float.NonNil(pnt.WriteValue)
								err = inst.doWrite(pnt, net.UUID, dev.UUID, writeVal)
								if err != nil {
									err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
									continue // next point
								}
								pnt, err = inst.pointUpdate(pnt, writeVal, true, true)
								continue // next point
							} else { // Do these actions fort AV, AO, BV, BO
								// We need to read the priority array of our FF point, and the BACnet Server point, then update the FF point and BACnet Server point

								// Get Priority array of FF
								currentFFPointPriorityMap := priorityarray.ConvertToMap(*pnt.Priority)

								// Get Priority Array of BACnet Server Point
								inst.bacnetDebugMsg("BACnetServerPolling() Read Priority")
								var currentBACServPriority *priority.Float32
								currentBACServPriority, err = inst.doReadPriority(pnt, net.UUID, dev.UUID)
								if err != nil {
									inst.bacnetErrorMsg("BACnetServerPolling(): doReadPriority error:", err)
									inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
									continue // next point
								}
								if currentBACServPriority == nil {
									inst.bacnetErrorMsg("BACnetServerPolling(): BACnet Server Point returned an empty priority array; skipped")
									inst.pointUpdateErr(pnt, "BACnet Server Point returned an empty priority array", model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
									continue // next point
								}
								currentBACServPriorityMap := ConvertPriorityToMap(*currentBACServPriority)
								for key, val := range currentBACServPriorityMap {
									inst.bacnetDebugMsg("BACnetServerPolling() currentBACServPriorityMap key: ", key, "val", val)
								}

								// var pointOperationError string
								// loop through the FF Point priority array and check for differences
								priorityArraysMatch := true
								for key, FFPointVal := range currentFFPointPriorityMap {
									inst.bacnetDebugMsg("BACnetServerPolling() currentPriorityMap key: ", nstring.NewString(key).RemoveSpecialCharacter(), "val", FFPointVal)
									if BACServVal, ok := currentBACServPriorityMap[key]; ok { // If the matching priority exists on the BACnet Server point,
										if !FFPointAndBACnetServerPointAreEqual(FFPointVal, BACServVal) { // Points are NOT equal
											priorityArraysMatch = false
											if boolean.IsFalse(pnt.WritePollRequired) { // No Writes Required, so match FF point to BACnet Server point
												currentFFPointPriorityMap[key] = float.New(float64(float.NonNil32(BACServVal)))
											} else { // Write is required, so set the BACnet Server point to match FF Point values
												if FFPointVal == nil { // Do priority release on nil
													var priorityAsInt int
													priorityAsInt, err = strconv.Atoi(nstring.NewString(key).RemoveSpecialCharacter())
													if err != nil {
														inst.bacnetErrorMsg("BACnetServerPolling(): cannot parse priority", key)
														// pointOperationError = "cannot parse priority"
														continue // next priority
													} else {
														err = inst.doRelease(pnt, net.UUID, dev.UUID, uint8(priorityAsInt))
														if err != nil {
															inst.bacnetErrorMsg("BACnetServerPolling(): doWrite error:", err)
															// pointOperationError = err.Error()
															continue // next priority
														}
													}
												} else { // Otherwise write the value
													err = inst.doWrite(pnt, net.UUID, dev.UUID, float.NonNil(FFPointVal))
													if err != nil {
														inst.bacnetErrorMsg("BACnetServerPolling(): doWrite error:", err)
														// pointOperationError = err.Error()
													}
												}
											}
										}
									} else {
										inst.bacnetErrorMsg("BACnetServerPolling(): BACnet Server Point is missing priority: ", key)
									}
								}
								if !priorityArraysMatch {
									inst.bacnetDebugMsg("BACnetServerPolling() PRIORITY ARRAYS ARE DIFFERENT")
									// err = inst.pointUpdateErr(pnt, pointOperationError, model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
									// Update point
									var readFloat float64
									readFloat, err = inst.doReadValue(pnt, dev.NetworkUUID, dev.UUID)
									if err != nil {
										inst.bacnetErrorMsg("BACnetServerPolling(): doReadValue error:", err)
									} else {
										pnt, err = inst.pointUpdate(pnt, readFloat, true, true)
										if err != nil {
											inst.bacnetErrorMsg("BACnetServerPolling(): pointUpdate() error:", err)
										}
									}
								} else {
									inst.bacnetDebugMsg("BACnetServerPolling() PRIORITY ARRAYS ARE IDENTICAL")
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
