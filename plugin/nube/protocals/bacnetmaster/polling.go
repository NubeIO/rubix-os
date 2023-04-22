package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/services/pollqueue"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/writemode"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
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
	deviceDelay = 100 * time.Millisecond
	pointDelay = 100 * time.Millisecond
	if networkType == model.TransType.LoRa {
		deviceDelay = 80 * time.Millisecond
		pointDelay = 6000 * time.Millisecond
	}
	return
}

func (inst *Instance) getNetworkPollManagerByUUID(netUUID string) (*pollqueue.NetworkPollManager, error) {
	for _, netPollMan := range inst.NetworkPollManagers {
		if netPollMan.FFNetworkUUID == netUUID {
			return netPollMan, nil
		}
	}
	return nil, errors.New("bacnet getNetworkPollManagerByUUID(): couldn't find NetworkPollManager")
}

var poll poller.Poller

// BACnetMasterPolling TODO: currently Polling loops through each network, grabs one point, and polls it.  Could be improved by having a seperate client/go routine for each of the networks.
func (inst *Instance) BACnetMasterPolling() error {
	poll = poller.New()
	var counter = 0
	f := func() (bool, error) {
		counter++
		// fmt.Println("\n \n")
		inst.bacnetDebugMsg("LOOP COUNT: ", counter)

		var netArg api.Args

		if len(inst.NetworkPollManagers) == 0 {
			inst.bacnetDebugMsg("NO BACNET NETWORKS FOUND")
			time.Sleep(15000 * time.Millisecond)
		}
		// inst.bacnetDebugMsg("inst.NetworkPollManagers")
		// inst.bacnetDebugMsg("%+v\n", inst.NetworkPollManagers)
		for _, netPollMan := range inst.NetworkPollManagers { // LOOP THROUGH AND POLL NEXT POINTS IN EACH NETWORK QUEUE
			// inst.bacnetDebugMsg("BACnetPolling: netPollMan ", netPollMan.FFNetworkUUID)
			if netPollMan.PortUnavailableTimeout != nil {
				inst.bacnetDebugMsg("bacnet port unavailable. polling paused.")
				continue
			}
			pollStartTime := time.Now()
			// Check that network exists
			// inst.bacnetDebugMsg("netPollMan")
			// inst.bacnetDebugMsg("%+v\n", netPollMan)
			net, err := inst.db.GetNetwork(netPollMan.FFNetworkUUID, netArg)
			// inst.bacnetDebugMsg("net")
			// inst.bacnetDebugMsg("%+v\n", net)
			// inst.bacnetDebugMsg("err")
			// inst.bacnetDebugMsg("%+v\n", err)
			if err != nil || net == nil || net.PluginConfId != inst.pluginUUID {
				inst.bacnetDebugMsg("BACNET NETWORK NOT FOUND")
				continue
			}
			// inst.bacnetDebugMsg(fmt.Sprintf("bacnet-poll: POLL START: NAME: %s\n", net.Name))

			if !boolean.IsTrue(net.Enable) {
				inst.bacnetDebugMsg(fmt.Sprintf("NETWORK DISABLED: NAME: %s", net.Name))
				continue
			}

			pp, callback := netPollMan.GetNextPollingPoint() // callback function is called once polling is completed.
			if pp == nil {
				// inst.bacnetDebugMsg("No PollingPoint available in Network ", net.UUID)
				continue
			}

			if pp.FFNetworkUUID != net.UUID {
				inst.bacnetErrorMsg("PollingPoint FFNetworkUUID does not match the Network UUID")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, false, pollqueue.DELAYED_RETRY, callback)
				continue
			}
			netPollMan.PrintPollQueuePointUUIDs()
			netPollMan.PrintPollingPointDebugInfo(pp)

			var devArg api.Args
			dev, err := inst.db.GetDevice(pp.FFDeviceUUID, devArg)
			if dev == nil || err != nil {
				inst.bacnetErrorMsg("could not find deviceID:", pp.FFDeviceUUID)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.DELAYED_RETRY, callback)
				continue
			}
			if !boolean.IsTrue(dev.Enable) {
				inst.bacnetErrorMsg("device is disabled.")
				inst.db.SetErrorsForAllPointsOnDevice(dev.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.NEVER_RETRY, callback)
				continue
			}
			pnt, err := inst.db.GetPoint(pp.FFPointUUID, api.Args{WithPriority: true})
			if pnt == nil || err != nil {
				inst.bacnetErrorMsg("could not find pointID: ", pp.FFPointUUID)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.DELAYED_RETRY, callback)
				continue
			}

			inst.printPointDebugInfo(pnt)

			if pnt.Priority == nil {
				inst.bacnetErrorMsg("BACnetMasterPolling: HAD TO ADD PRIORITY ARRAY")
				pnt.Priority = &model.Priority{}
			}

			if !boolean.IsTrue(pnt.Enable) {
				inst.bacnetErrorMsg("point is disabled.")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.NEVER_RETRY, callback)
				continue
			}

			if pnt.AddressID == nil {
				inst.bacnetPollingMsg(fmt.Sprintf("NEXT POLL DRAWN! : Network: %s, Device: %s, Point: %s, Priority: %s, Device-Add: %d, Point-Add: %d, Point Type: %s, WriteRequired: %t, ReadRequired: %t", net.Name, dev.Name, pnt.Name, pnt.PollPriority, dev.AddressId, pnt.AddressID, pnt.ObjectType, boolean.IsTrue(pnt.WritePollRequired), boolean.IsTrue(pnt.ReadPollRequired)))
			} else {
				inst.bacnetPollingMsg(fmt.Sprintf("NEXT POLL DRAWN! : Network: %s, Device: %s, Point: %s, Priority: %s, Device-Add: %d, Point-Add: %d, Point Type: %s, WriteRequired: %t, ReadRequired: %t", net.Name, dev.Name, pnt.Name, pnt.PollPriority, dev.AddressId, *pnt.AddressID, pnt.ObjectType, boolean.IsTrue(pnt.WritePollRequired), boolean.IsTrue(pnt.ReadPollRequired)))
			}

			writemode.SetPriorityArrayModeBasedOnWriteMode(pnt) // ensures the point PointPriorityArrayMode is set correctly

			currentBACServPriority, highestPriorityValue, readSuccess, writeSuccess, err := inst.doRead(pnt, net.UUID, dev.UUID)
			if err != nil {
				err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, true, false, pollqueue.IMMEDIATE_RETRY, callback)
				continue
			}

			// this resets IsTypeBool to correct setting (fixes issues from points created before this was fixed in updatePoint()
			isTypeBool := checkForBooleanType(pnt.ObjectType, pnt.DataType)
			pnt.IsTypeBool = nils.NewBool(isTypeBool)

			if currentBACServPriority != nil {
				currentBACServPriorityMap := ConvertPriorityToMap(*currentBACServPriority)
				pnt.CommonFault.InFault = false
				pnt.CommonFault.MessageLevel = model.MessageLevel.Info
				pnt.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
				pnt.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
				pnt.CommonFault.LastOk = time.Now().UTC()
				_, err = inst.pointUpdateFromPriorityArray(pnt, currentBACServPriorityMap, highestPriorityValue)
				if err != nil {
					inst.bacnetDebugMsg("BACnetMasterPolling() pointUpdateFromPriorityArray err: ", err)
				}
			} else {
				inst.pointUpdate(pnt, highestPriorityValue, readSuccess)
				if pnt.ObjectId != nil {
					inst.bacnetDebugMsg("updated bacnet point present value.  point: ", pnt.Name, " bacnet-id: ", pnt.ObjectType, *pnt.ObjectId)
				}
			}

			if !readSuccess && !writeSuccess {
				// inst.bacnetErrorMsg("no read/write performed: ", err)
				if err != nil {
					err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointError)
				}
				netPollMan.PollingFinished(pp, pollStartTime, false, false, true, false, pollqueue.IMMEDIATE_RETRY, callback)
				continue
			}

			/*
				//JUST FOR TESTING
				pnt, err = inst.db.GetPoint(pp.FFPointUUID)
				if pnt == nil || err != nil {
					log.Errorf("bacnet: AFTER... could not find pointID : %s\n", pp.FFPointUUID)
				}
			*/

			// This callback function triggers the PollManager to evaluate whether the point should be re-added to the PollQueue (Never, Immediately, or after the Poll Rate Delay)
			if pnt.ObjectId != nil {
				inst.bacnetDebugMsg("BACnetMasterPolling() success.  point: ", pnt.Name, " bacnet-id: ", pnt.ObjectType, *pnt.ObjectId)
			}
			netPollMan.PollingFinished(pp, pollStartTime, writeSuccess, readSuccess, true, false, pollqueue.NORMAL_RETRY, callback)
		}
		return false, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	inst.pollingCancel = cancel
	go poll.GoPoll(ctx, f)
	return nil
}
