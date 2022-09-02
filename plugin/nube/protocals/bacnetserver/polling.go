package main

import (
	"context"
	"errors"
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
		if bacnetStarted {
			inst.bacnetErrorMsg("-----------------------------------------")
			inst.bacnetErrorMsg("BACNET JUST STARTED")
			inst.bacnetErrorMsg("-----------------------------------------")
		}
		var err error
		nets, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{WithDevices: true})
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
					if len(net.Devices) == 0 {
						time.Sleep(2 * time.Second) // Delay to prevent unnecessary looping
					}
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
						if len(dev.Points) == 0 {
							time.Sleep(2 * time.Second) // Delay to prevent unnecessary looping
						}
						for _, pnt := range dev.Points { // POINTS
							time.Sleep(pointDelay) // DELAY between points
							if counter == 1 || bacnetStarted {
								inst.bacnetDebugMsg("FIRST START, WRITE ALL POINTS TO BACNET SERVER")
								pnt, err = inst.SyncFFPointWithBACnetServerPoint(pnt, dev.UUID, net.UUID, true)
							} else {
								pnt, err = inst.SyncFFPointWithBACnetServerPoint(pnt, dev.UUID, net.UUID, false)
							}
							if err != nil {
								inst.bacnetErrorMsg(err)
								continue // next point
							}
						}

						rsyncWrite = counter % 5

						timeEnd := time.Now()
						diff := timeEnd.Sub(timeStart)
						out := time.Time{}.Add(diff)
						inst.bacnetDebugMsg(fmt.Sprintf("poll-loop: NETWORK-NAME:%s POLL-DURATION: %s  POLL-COUNT: %d\n", net.Name, out.Format("15:04:05.000"), counter))
					}
				}
			}
		}
		bacnetStarted = false
		return false, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	inst.pollingCancel = cancel
	go poll.GoPoll(ctx, f)
	return nil
}

func (inst *Instance) SyncFFPointWithBACnetServerPoint(pnt *model.Point, devUUID, netUUID string, forceWrite bool) (*model.Point, error) {
	pnt, err := inst.db.GetPoint(pnt.UUID, api.Args{WithPriority: true})
	if err != nil {
		inst.bacnetErrorMsg("SyncFFPointWithBACnetServerPoint(): Point not found")
		return nil, errors.New("SyncFFPointWithBACnetServerPoint(): Point not found")
	}
	if boolean.IsFalse(pnt.Enable) {
		inst.bacnetErrorMsg("SyncFFPointWithBACnetServerPoint(): Point is disabled; skipped")
		return nil, errors.New("SyncFFPointWithBACnetServerPoint(): Point is disabled; skipped")
	}
	if pnt.Priority == nil {
		inst.bacnetErrorMsg("SyncFFPointWithBACnetServerPoint(): Point doesn't have a priority array; skipped")
		err = inst.pointUpdateErr(pnt, "no priority array found", model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
		return nil, errors.New("SyncFFPointWithBACnetServerPoint(): Point doesn't have a priority array; skipped")
	}
	if !isWriteableObjectType(pnt.ObjectType) { // Do these actions for AI, BI
		// For AI/BI there is no priority array we just write the values to match our FF point WriteValue
		writeVal := float.NonNil(pnt.WriteValue)
		err = inst.doWrite(pnt, netUUID, devUUID, writeVal, 16)
		if err != nil {
			err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
			return nil, err
		}
		pnt, err = inst.pointUpdate(pnt, writeVal, true, true)
		return pnt, nil
	} else { // Do these actions fort AV, AO, BV, BO
		// We need to read the priority array of our FF point, and the BACnet Server point, then update the FF point and BACnet Server point

		// Get Priority array of FF
		currentFFPointPriorityMap := priorityarray.ConvertToMap(*pnt.Priority)

		// Get Priority Array of BACnet Server Point
		var currentBACServPriority *priority.Float32
		currentBACServPriority, err = inst.doReadPriority(pnt, netUUID, devUUID)
		if err != nil {
			inst.bacnetErrorMsg("BACnetServerPolling(): doReadPriority error:", err)
			inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
			return nil, err
		}
		if currentBACServPriority == nil {
			inst.bacnetErrorMsg("BACnetServerPolling(): BACnet Server Point returned an empty priority array; skipped")
			inst.pointUpdateErr(pnt, "BACnet Server Point returned an empty priority array", model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
			return nil, errors.New("BACnetServerPolling(): BACnet Server Point returned an empty priority array; skipped")
		}
		currentBACServPriorityMap := ConvertPriorityToMap(*currentBACServPriority)
		for key, val := range currentBACServPriorityMap {
			if val == nil {
				inst.bacnetDebugMsg("BACnetServerPolling() BACnetServerPoint: key: ", key, "val", val)
			} else {
				inst.bacnetDebugMsg("BACnetServerPolling() BACnetServerPoint: key: ", key, "val", float.NonNil32(val))
			}
		}

		// var pointOperationError string
		// loop through the FF Point priority array and check for differences
		priorityArraysMatch := true
		for key, FFPointVal := range currentFFPointPriorityMap {
			if FFPointVal == nil {
				inst.bacnetDebugMsg("BACnetServerPolling() FFPoint: key: ", key, "val", FFPointVal)
			} else {
				inst.bacnetDebugMsg("BACnetServerPolling() FFPoint: key: ", key, "val", float.NonNil(FFPointVal))
			}
			if BACServVal, ok := currentBACServPriorityMap[key]; ok { // If the matching priority exists on the BACnet Server point,
				if !FFPointAndBACnetServerPointAreEqual(FFPointVal, BACServVal) { // Points are NOT equal
					priorityArraysMatch = false
					if boolean.IsFalse(pnt.WritePollRequired) && !forceWrite { // No Writes Required, so match FF point to BACnet Server point
						if BACServVal == nil {
							currentFFPointPriorityMap[key] = nil
						} else {
							currentFFPointPriorityMap[key] = float.New(float64(float.NonNil32(BACServVal)))
						}
					} else { // Write is required, so set the BACnet Server point to match FF Point values
						var priorityAsInt int
						priorityAsInt, err = strconv.Atoi(nstring.NewString(key).RemoveSpecialCharacter())
						if FFPointVal == nil { // Do priority release on nil
							if err != nil {
								inst.bacnetErrorMsg("BACnetServerPolling(): cannot parse priority", key)
								// pointOperationError = "cannot parse priority"
								continue // next priority
							} else {
								err = inst.doRelease(pnt, netUUID, devUUID, uint8(priorityAsInt))
								if err != nil {
									inst.bacnetErrorMsg("BACnetServerPolling(): doWrite error:", err)
									// pointOperationError = err.Error()
									continue // next priority
								}
							}
						} else { // Otherwise write the value
							err = inst.doWrite(pnt, netUUID, devUUID, float.NonNil(FFPointVal), uint8(priorityAsInt))
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
		if boolean.IsFalse(pnt.WritePollRequired) && !forceWrite {
			inst.bacnetDebugMsg("BACnetServerPolling() WRITE NOT REQUIRED")
		} else {
			inst.bacnetDebugMsg("BACnetServerPolling() WRITE IS REQUIRED")
		}
		if !priorityArraysMatch {
			inst.bacnetDebugMsg("BACnetServerPolling() PRIORITY ARRAYS ARE DIFFERENT")
			// err = inst.pointUpdateErr(pnt, pointOperationError, model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
			// Update point priority array
			pnt, err = priorityarray.ApplyMapToPriorityArray(pnt, &currentFFPointPriorityMap)
			if err != nil {
				inst.bacnetErrorMsg("BACnetServerPolling(): ApplyMapToPriorityArray() error:", err)
			}
			// Update point PresentValue
			var readFloat float64
			readFloat, err = inst.doReadValue(pnt, netUUID, devUUID)
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
	return pnt, nil
}
