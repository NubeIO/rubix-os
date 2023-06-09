package pollqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"time"
)

// REFS:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

func (pm *NetworkPollManager) RebuildPollingQueue() error {
	// TODO: STOP ANY OTHER QUEUE LOADERS
	pm.pollQueueDebugMsg("RebuildPollingQueue()")
	wasRunning := pm.PluginQueueUnloader != nil
	pm.EmptyQueue()
	var arg argspkg.Args
	arg.WithDevices = true
	arg.WithPoints = true
	net, err := pm.DBHandlerRef.GetNetwork(pm.FFNetworkUUID, arg)
	if err != nil || net.Devices == nil || len(net.Devices) == 0 {
		pm.pollQueueDebugMsg("RebuildPollingQueue() couldn't find any devices for the network %s", pm.FFNetworkUUID)
		return errors.New(fmt.Sprintf("NetworkPollManager.RebuildPollingQueue: couldn't find any devices for the network %s", pm.FFNetworkUUID))
	}
	devs := net.Devices
	for _, dev := range devs { // DEVICES
		if dev.NetworkUUID == pm.FFNetworkUUID && boolean.IsTrue(dev.Enable) {
			for _, pnt := range dev.Points { // POINTS
				if pnt.DeviceUUID == dev.UUID && boolean.IsTrue(pnt.Enable) {
					pp := NewPollingPoint(pnt.UUID, pnt.DeviceUUID, dev.NetworkUUID, pm.FFPluginUUID)
					pp.PollPriority = pnt.PollPriority
					pm.pollQueueDebugMsg(fmt.Sprintf("RebuildPollingQueue() pp: %+v", pp))
					pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
					pm.PollQueue.AddPollingPoint(pp)
				} else {
					pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.RebuildPollingQueue: Point (%s) is not enabled", pnt.UUID))
				}
			}
		} else {
			pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.RebuildPollingQueue: Device (%s) is not enabled", dev.UUID))
		}
	}
	heap.Init(pm.PollQueue.PriorityQueue)
	if wasRunning {
		pm.StartQueueUnloader()
	}
	// TODO: START ANY OTHER REQUIRED QUEUE LOADERS/OPTIMIZERS
	// pm.PrintPollQueuePointUUIDs()
	return nil
}

func (pm *NetworkPollManager) PollingPointCompleteNotification(pp *PollingPoint, writeSuccess, readSuccess bool, pollTimeSecs float64, pointUpdate, resetToConfiguredPriority bool, retryType PollRetryType, actualPoll, pollingWasNotRequired, justToReAdd bool) {
	if !justToReAdd {
		pm.pollQueuePollingMsg(fmt.Sprintf("POLLING COMPLETE: Point UUID: %s, writeSuccess: %t, readSuccess: %t, pointUpdate: %t, actualPoll: %t, pollingWasNotRequired: %t, justToReAdd: %t, retryType: %s, pollTime: %f", pp.FFPointUUID, writeSuccess, readSuccess, pointUpdate, actualPoll, pollingWasNotRequired, justToReAdd, retryType, pollTimeSecs))
	}

	if !actualPoll && !justToReAdd { // This posts the next polling point immediately (faster than the MaxPollRate) because no poll was actually made.
		pm.postNextPointCallback()
	}

	if !justToReAdd {
		_, success := pm.PollQueue.OutstandingPollingPoints.RemovePollingPointByPointUUID(pp.FFPointUUID)
		if !success {
			// It's ok to get this error message when adding/re-adding points.
			pm.pollQueueErrorMsg("NetworkPollManager.PollingPointCompleteNotification(): couldn't find polling point in OutstandingPollingPoints, %s", pp.FFPointUUID)
		}
	}

	if !pointUpdate {
		pm.PollCompleteStatsUpdate(pp, pollTimeSecs) // This will update the relevant PollManager statistics.
	}

	point, err := pm.DBHandlerRef.GetPoint(pp.FFPointUUID, argspkg.Args{WithPriority: true})
	if point == nil || err != nil {
		pm.pollQueueErrorMsg("NetworkPollManager.PollingPointCompleteNotification(): couldn't find point %s", pp.FFPointUUID)
		return
	}
	// TODO: potentially only required on writeSuccess (but possibility of lockup on a bad point)
	// Reset poll priority to set value (in cases where pp has been escalated to ASAP).
	if resetToConfiguredPriority {
		pp.PollPriority = point.PollPriority
	}

	val, ok := pm.PollQueue.PointsUpdatedWhilePolling[point.UUID]
	if ok {
		delete(pm.PollQueue.PointsUpdatedWhilePolling, point.UUID)
		if val == true { // point needs an ASAP write
			pp.PollPriority = model.PRIORITY_ASAP
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
			return
		}
	}

	// pm.pollQueuePollingMsg(fmt.Sprintf("PollingPointCompleteNotification: point %+v", point))
	// pm.PrintPointDebugInfo(point)

	// If the device was deleted while this point was being polled, discard the PollingPoint
	if !pointUpdate && !pm.PollQueue.CheckIfActiveDevicesListIncludes(point.DeviceUUID) {
		return
	}

	switch point.WriteMode {
	case model.ReadOnce: // ReadOnce          If read_successful then don't re-add.
		point.WritePollRequired = boolean.NewFalse()
		if retryType == NEVER_RETRY || ((readSuccess || pollingWasNotRequired) && retryType == NORMAL_RETRY) {
			point.ReadPollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if (boolean.IsTrue(point.ReadPollRequired) && !readSuccess && retryType == NORMAL_RETRY) || retryType == IMMEDIATE_RETRY {
			point.ReadPollRequired = boolean.NewTrue()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)
		} else if retryType == DELAYED_RETRY {
			point.ReadPollRequired = boolean.NewTrue()
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		}

	case model.ReadOnly: // ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.WritePollRequired = boolean.NewFalse()
		point.ReadPollRequired = boolean.NewTrue()
		if ((readSuccess || pollingWasNotRequired) && retryType == NORMAL_RETRY) || retryType == DELAYED_RETRY {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if (!readSuccess && retryType == NORMAL_RETRY) || retryType == IMMEDIATE_RETRY {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.  TODO: This might conflict with pausing polling on PortUnavailable
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		} else if retryType == NEVER_RETRY {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		}

	case model.WriteOnce: // WriteOnce         If write_successful then don't re-add.
		point.ReadPollRequired = boolean.NewFalse()
		if ((writeSuccess || pollingWasNotRequired) && retryType == NORMAL_RETRY) || retryType == NEVER_RETRY {
			point.WritePollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if (boolean.IsTrue(point.WritePollRequired) && !writeSuccess && retryType == NORMAL_RETRY) || retryType == IMMEDIATE_RETRY {
			point.WritePollRequired = boolean.NewTrue() // TODO: this might cause these points to write more than once.
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		} else if retryType == DELAYED_RETRY {
			point.WritePollRequired = boolean.NewTrue()
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		}

	case model.WriteOnceReadOnce: // WriteOnceReadOnce     If write_successful and read_success then don't re-add.
		if (boolean.IsTrue(point.WritePollRequired) && writeSuccess && retryType == NORMAL_RETRY) || retryType == NEVER_RETRY {
			point.WritePollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if pointUpdate || (boolean.IsTrue(point.WritePollRequired) && !writeSuccess && retryType == NORMAL_RETRY) || retryType == IMMEDIATE_RETRY {
			point.WritePollRequired = boolean.NewTrue()
			if pointUpdate {
				point.ReadPollRequired = boolean.NewTrue()
			}
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
			break
		} else if retryType == DELAYED_RETRY {
			point.WritePollRequired = boolean.NewTrue()
			point.ReadPollRequired = boolean.NewTrue()
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
			break
		}
		if readSuccess && retryType == NORMAL_RETRY || retryType == NEVER_RETRY {
			point.ReadPollRequired = boolean.NewFalse()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if boolean.IsTrue(point.ReadPollRequired) && !readSuccess && retryType == NORMAL_RETRY || retryType == IMMEDIATE_RETRY {
			point.ReadPollRequired = boolean.NewTrue()
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		}

	case model.WriteAlways: // WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = boolean.NewFalse()
		point.WritePollRequired = boolean.NewTrue()
		if (writeSuccess && retryType == NORMAL_RETRY) || retryType == DELAYED_RETRY {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
			break
		} else if (!writeSuccess && retryType == NORMAL_RETRY) || retryType == IMMEDIATE_RETRY {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
			break
		} else if retryType == NEVER_RETRY {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		}

	case model.WriteOnceThenRead: // WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = boolean.NewTrue()
		if retryType == NEVER_RETRY {
			if writeSuccess {
				point.WritePollRequired = boolean.NewFalse()
			}
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		} else if pointUpdate || (boolean.IsTrue(point.WritePollRequired) && !writeSuccess && retryType == NORMAL_RETRY) || retryType == IMMEDIATE_RETRY {
			if writeSuccess {
				point.WritePollRequired = boolean.NewFalse()
			}
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
			break
		} else if (boolean.IsTrue(point.WritePollRequired) && writeSuccess && retryType == NORMAL_RETRY) || retryType == DELAYED_RETRY {
			if writeSuccess {
				point.WritePollRequired = boolean.NewFalse()
			}
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
			break
		}
		if readSuccess && retryType == NORMAL_RETRY || retryType == DELAYED_RETRY {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
			break
		} else if !readSuccess && retryType == NORMAL_RETRY || retryType == IMMEDIATE_RETRY {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
		}

	case model.WriteAndMaintain: // WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = boolean.NewTrue()
		// pm.pollQueueDebugMsg(fmt.Sprintf("WriteAndMaintain point %+v\n", point))
		if (boolean.IsTrue(point.WritePollRequired) && !writeSuccess && retryType == NORMAL_RETRY) || retryType == IMMEDIATE_RETRY {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			// point.WritePollRequired = boolean.NewTrue()
			pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
			pm.PollQueue.AddPollingPoint(pp)
			break
		} else if retryType == DELAYED_RETRY {
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
			break
		} else if retryType == NEVER_RETRY {
			if pp.RepollTimer != nil {
				pp.RepollTimer.Stop()
				pp.RepollTimer = nil
			}
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
			break
		}

		if point.WriteValue != nil {
			noPV := true
			var readValue float64
			if point.OriginalValue != nil {
				noPV = false
				readValue = *point.OriginalValue
			}
			if noPV || readValue != *point.WriteValue {
				if pp.RepollTimer != nil {
					pp.RepollTimer.Stop()
					pp.RepollTimer = nil
				}
				point.WritePollRequired = boolean.NewTrue()
				pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
				pm.PollQueue.AddPollingPoint(pp)                              // re-add to poll queue immediately
			} else {
				point.WritePollRequired = boolean.NewFalse()
				duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
				// This line sets a timer to re-add the point to the poll queue after the PollRate time.
				pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
				addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
				if !addSuccess {
					pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
				}
			}
		} else {
			// If WriteValue is nil we still need to re-add the point to perform a read
			point.WritePollRequired = boolean.NewFalse()
			duration := pm.GetPollRateDuration(point.PollRate, pp.FFDeviceUUID)
			// This line sets a timer to re-add the point to the poll queue after the PollRate time.
			pp.RepollTimer = time.AfterFunc(duration, pm.MakePollingPointRepollCallback(pp, point.WriteMode))
			addSuccess := pm.PollQueue.StandbyPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus PollingPointCompleteNotification(): polling point could not be added to StandbyPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		}
	}

	// pm.pollQueueDebugMsg(fmt.Sprintf("PollingPointCompleteNotification (ABOUT TO DB UPDATE): point  %+v", point))
	// point.PrintPointValues()
	// TODO: WOULD BE GOOD IF THIS COULD BE MOVED TO app.go
	resetFaults := readSuccess || writeSuccess
	if resetFaults {
		point.CommonFault.InFault = false
		point.CommonFault.MessageLevel = model.MessageLevel.Info
		point.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
		point.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
		point.CommonFault.LastOk = time.Now().UTC()
	}
	point, err = pm.DBHandlerRef.UpdatePoint(point.UUID, point)
	// printPointDebugInfo(point)
	pm.pollQueuePollingMsg("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
}

func (pm *NetworkPollManager) MakePollingPointRepollCallback(pp *PollingPoint, writeMode model.WriteMode) func() {
	// log.Info("MakePollingPointRepollCallback()")
	f := func() {
		// pm.pollQueueDebugMsg(fmt.Sprintf("CALL PollingPointRepollCallback func() pp: %+v", pp))
		pp.RepollTimer = nil
		_, removeSuccess := pm.PollQueue.StandbyPollingPoints.RemovePollingPointByPointUUID(pp.FFPointUUID)
		if !removeSuccess {
			pm.pollQueueErrorMsg(fmt.Sprintf("Modbus MakePollingPointRepollCallback(): polling point could not be found in StandbyPollingPoints.  (%s)", pp.FFPointUUID))
		}
		/*
			point, err := pm.DBHandlerRef.GetPoint(pp.FFPointUUID, api.Args{WithPriority: true})
			if point == nil || err != nil {
				pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.PollingPointCompleteNotification(): couldn't find point %s", pp.FFPointUUID))
				return
			}

				pointUpdateReq := false

				switch writeMode {
				case model.ReadOnce:

				case model.ReadOnly: // ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
					point.ReadPollRequired = boolean.NewTrue()
					point.WritePollRequired = boolean.NewFalse()
					pointUpdateReq = true

				case model.WriteOnce: // WriteOnce         If write_successful then don't re-add.

				case model.WriteOnceReadOnce: // WriteOnceReadOnce     If write_successful and read_success then don't re-add.

				case model.WriteAlways: // WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
					point.ReadPollRequired = boolean.NewFalse()
					point.WritePollRequired = boolean.NewTrue()
					pointUpdateReq = true

				case model.WriteOnceThenRead: // WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
					point.ReadPollRequired = boolean.NewTrue()
					point.WritePollRequired = boolean.NewFalse()
					pointUpdateReq = true

				case model.WriteAndMaintain: // WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
					point.ReadPollRequired = boolean.NewTrue()
					point.WritePollRequired = boolean.NewFalse()
					pointUpdateReq = true
				}
		*/

		// Now add the polling point back to the polling queue
		pp.LockupAlertTimer = pm.MakeLockupTimerFunc(pp.PollPriority) // starts a countdown for queue lockup alerts.
		pm.PollQueue.AddPollingPoint(pp)

		/*
			if pointUpdateReq {
				// TODO: WOULD BE GOOD IF THIS COULD BE MOVED TO app.go
				// pm.pollQueueDebugMsg(fmt.Sprintf("pm.DBHandlerRef: %+v", pm.DBHandlerRef))
				point, err = pm.DBHandlerRef.UpdatePoint(point.UUID, point)
				if err != nil || point == nil {
					pm.pollQueueErrorMsg(fmt.Sprintf("point DB UPDATE FAILED Err: %+v", err))
					return
				}
				// pm.pollQueueDebugMsg(fmt.Sprintf("point after DB UPDATE: %+v", point))
				// printPointDebugInfo(point)
			}
		*/
	}
	return f
}

func (pm *NetworkPollManager) MakeLockupTimerFunc(priority model.PollPriority) *time.Timer {
	timeoutDuration := 5 * time.Minute

	switch priority {
	case model.PRIORITY_ASAP:
		timeoutDuration = pm.ASAPPriorityMaxCycleTime

	case model.PRIORITY_HIGH:
		timeoutDuration = pm.HighPriorityMaxCycleTime

	case model.PRIORITY_NORMAL:
		timeoutDuration = pm.NormalPriorityMaxCycleTime

	case model.PRIORITY_LOW:
		timeoutDuration = pm.LowPriorityMaxCycleTime

	}

	f := func() {
		pm.pollQueueDebugMsg("Polling Lockout Timer Expired! Polling Priority: %d,  Polling Network: %s", priority, pm.FFNetworkUUID)
		plugin, err := pm.DBHandlerRef.GetPlugin(pm.FFPluginUUID)
		switch priority {
		case model.PRIORITY_ASAP:
			pm.ASAPPriorityLockupAlert = true
			if plugin != nil && err == nil {
				pm.pollQueueErrorMsg(fmt.Sprintf("%s Plugin: ASAP Priority Poll Queue LOCKUP", plugin.Name))
			} else {
				pm.pollQueueErrorMsg("Unknown Plugin: ASAP Priority Poll Queue LOCKUP")
			}
			// TODO: update network error to show queue lockup

		case model.PRIORITY_HIGH:
			pm.HighPriorityLockupAlert = true
			if plugin != nil && err == nil {
				pm.pollQueueErrorMsg(fmt.Sprintf("%s Plugin: HIGH Priority Poll Queue LOCKUP", plugin.Name))
			} else {
				pm.pollQueueErrorMsg("Unknown Plugin: HIGH Priority Poll Queue LOCKUP")
			}
			// TODO: update network error to show queue lockup

		case model.PRIORITY_NORMAL:
			pm.NormalPriorityLockupAlert = true
			if plugin != nil && err == nil {
				pm.pollQueueErrorMsg(fmt.Sprintf("%s Plugin: NORMAL Priority Poll Queue LOCKUP", plugin.Name))
			} else {
				pm.pollQueueErrorMsg("Unknown Plugin: NORMAL Priority Poll Queue LOCKUP")
			}
			// TODO: update network error to show queue lockup

		case model.PRIORITY_LOW:
			pm.LowPriorityLockupAlert = true
			if plugin != nil && err == nil {
				pm.pollQueueErrorMsg(fmt.Sprintf("%s Plugin: LOW Priority Poll Queue LOCKUP", plugin.Name))
			} else {
				pm.pollQueueErrorMsg("Unknown Plugin: LOW Priority Poll Queue LOCKUP")
			}
			// TODO: update network error to show queue lockup

		}
	}
	return time.AfterFunc(timeoutDuration, f)
}

func (pm *NetworkPollManager) SetPointPollRequiredFlagsBasedOnWriteMode(point *model.Point) {

	if point == nil {
		pm.pollQueueDebugMsg("NetworkPollManager.SetPointPollRequiredFlagsBasedOnWriteMode(): couldn't find point")
		return
	}

	switch point.WriteMode {
	case model.ReadOnce:
		return

	case model.ReadOnly: // ReadOnly          Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewFalse()

	case model.WriteOnce: // WriteOnce         If write_successful then don't re-add.
		return

	case model.WriteOnceReadOnce: // WriteOnceReadOnce     If write_successful and read_success then don't re-add.
		return

	case model.WriteAlways: // WriteAlways       Re-add with ReadPollRequired false, WritePollRequired true. confirm that a successful write ensures the value is set to the write value.
		point.ReadPollRequired = boolean.NewFalse()
		point.WritePollRequired = boolean.NewTrue()

	case model.WriteOnceThenRead: // WriteOnceThenRead     If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewTrue()

	case model.WriteAndMaintain: // WriteAndMaintain    If write_successful: Re-add with ReadPollRequired true, WritePollRequired false.  Need to check that write value matches present value after each read poll.
		point.ReadPollRequired = boolean.NewTrue()
		point.WritePollRequired = boolean.NewTrue()
	}

	pm.DBHandlerRef.UpdatePoint(point.UUID, point)
}
