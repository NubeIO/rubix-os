package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/NubeIO/flow-framework/services/pollqueue"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/nurl"
	"github.com/NubeIO/flow-framework/utils/writemode"
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

func (inst *Instance) getNetworkPollManagerByUUID(netUUID string) (*pollqueue.NetworkPollManager, error) {
	for _, netPollMan := range inst.NetworkPollManagers {
		if netPollMan.FFNetworkUUID == netUUID {
			return netPollMan, nil
		}
	}
	return nil, errors.New("modbus getNetworkPollManagerByUUID(): couldn't find NetworkPollManager")
}

var poll poller.Poller

// ModbusPolling TODO: currently Polling loops through each network, grabs one point, and polls it.  Could be improved by having a seperate client/go routine for each of the networks.
func (inst *Instance) ModbusPolling() error {
	poll = poller.New()
	var counter = 0
	f := func() (bool, error) {
		counter++
		// fmt.Println("\n \n")
		inst.modbusDebugMsg("LOOP COUNT: ", counter)
		var netArg api.Args
		/*
			nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, netArg)
			if err != nil {
				return false, err
			}
		*/

		if len(inst.NetworkPollManagers) == 0 {
			inst.modbusDebugMsg("NO MODBUS NETWORKS FOUND")
			time.Sleep(15000 * time.Millisecond)
		}
		// inst.modbusDebugMsg("inst.NetworkPollManagers")
		// inst.modbusDebugMsg("%+v\n", inst.NetworkPollManagers)
		for _, netPollMan := range inst.NetworkPollManagers { // LOOP THROUGH AND POLL NEXT POINTS IN EACH NETWORK QUEUE
			// inst.modbusDebugMsg("ModbusPolling: netPollMan ", netPollMan.FFNetworkUUID)
			if netPollMan.PortUnavailableTimeout != nil {
				inst.modbusDebugMsg("modbus port unavailable. polling paused.")
				continue
			}
			pollStartTime := time.Now()
			// Check that network exists
			// inst.modbusDebugMsg("netPollMan")
			// inst.modbusDebugMsg("%+v\n", netPollMan)
			net, err := inst.db.GetNetwork(netPollMan.FFNetworkUUID, netArg)
			// inst.modbusDebugMsg("net")
			// inst.modbusDebugMsg("%+v\n", net)
			// inst.modbusDebugMsg("err")
			// inst.modbusDebugMsg("%+v\n", err)
			if err != nil || net == nil || net.PluginConfId != inst.pluginUUID {
				inst.modbusDebugMsg("MODBUS NETWORK NOT FOUND")
				continue
			}
			// inst.modbusDebugMsg(fmt.Sprintf("modbus-poll: POLL START: NAME: %s\n", net.Name))

			if !boolean.IsTrue(net.Enable) {
				inst.modbusDebugMsg(fmt.Sprintf("NETWORK DISABLED: NAME: %s", net.Name))
				continue
			}

			pp, callback := netPollMan.GetNextPollingPoint() // callback function is called once polling is completed.
			if pp == nil {
				// inst.modbusDebugMsg("No PollingPoint available in Network ", net.UUID)
				continue
			}

			if pp.FFNetworkUUID != net.UUID {
				inst.modbusErrorMsg("PollingPoint FFNetworkUUID does not match the Network UUID")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, false, pollqueue.DELAYED_RETRY, callback)
				continue
			}
			netPollMan.PrintPollQueuePointUUIDs()
			netPollMan.PrintPollingPointDebugInfo(pp)

			var devArg api.Args
			dev, err := inst.db.GetDevice(pp.FFDeviceUUID, devArg)
			if dev == nil || err != nil {
				inst.modbusErrorMsg("could not find deviceID:", pp.FFDeviceUUID)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.DELAYED_RETRY, callback)
				continue
			}
			if boolean.IsFalse(dev.Enable) {
				inst.modbusErrorMsg("device is disabled.")
				inst.db.SetErrorsForAllPointsOnDevice(dev.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.NEVER_RETRY, callback)
				continue
			}
			if dev.AddressId <= 0 || dev.AddressId >= 255 {
				inst.modbusErrorMsg("address is not valid.  modbus addresses must be between 1 and 254")
				inst.db.SetErrorsForAllPointsOnDevice(dev.UUID, "address out of range", model.MessageLevel.Critical, model.CommonFaultCode.ConfigError)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.NEVER_RETRY, callback)
				continue
			}

			pnt, err := inst.db.GetPoint(pp.FFPointUUID, api.Args{WithPriority: true})
			if pnt == nil || err != nil {
				inst.modbusErrorMsg("could not find pointID: ", pp.FFPointUUID)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.DELAYED_RETRY, callback)
				continue
			}

			inst.printPointDebugInfo(pnt)

			if pnt.Priority == nil {
				inst.modbusErrorMsg("ModbusPolling: HAD TO ADD PRIORITY ARRAY")
				pnt.Priority = &model.Priority{}
			}

			if !boolean.IsTrue(pnt.Enable) {
				inst.modbusErrorMsg("point is disabled.")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.NEVER_RETRY, callback)
				continue
			}

			inst.modbusDebugMsg(fmt.Sprintf("MODBUS POLL! : Priority: %s, Network: %s Device: %s Point: %s Device-Add: %d Point-Add: %d Point Type: %s, WriteRequired: %t, ReadRequired: %t", pp.PollPriority, net.UUID, dev.UUID, pnt.UUID, dev.AddressId, *pnt.AddressID, pnt.ObjectType, boolean.IsTrue(pnt.WritePollRequired), boolean.IsTrue(pnt.ReadPollRequired)))

			if !boolean.IsTrue(pnt.WritePollRequired) && !boolean.IsTrue(pnt.ReadPollRequired) {
				inst.modbusDebugMsg("polling not required on this point")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, true, pollqueue.NORMAL_RETRY, callback)
				continue
			}

			writemode.SetPriorityArrayModeBasedOnWriteMode(pnt) // ensures the point PointPriorityArrayMode is set correctly

			// SETUP MODBUS CLIENT CONNECTION
			var mbClient smod.ModbusClient
			// var dCheck devCheck
			// dCheck.devUUID = dev.UUID
			mbClient, err = inst.setClient(net, dev, true)
			if err != nil {
				inst.modbusErrorMsg(fmt.Sprintf("failed to set client error: %v. network name:%s", err, net.Name))
				if mbClient.PortUnavailable {
					netPollMan.PortUnavailable()
					unpauseFunc := func() {
						netPollMan.PortAvailable()
					}
					netPollMan.PortUnavailableTimeout = time.AfterFunc(10*time.Second, unpauseFunc)
				}
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, false, pollqueue.NORMAL_RETRY, callback)
				continue
			}
			if net.TransportType == model.TransType.Serial || net.TransportType == model.TransType.LoRa {
				if dev.AddressId >= 1 {
					mbClient.RTUClientHandler.SlaveID = byte(dev.AddressId)
				}
			} else if dev.TransportType == model.TransType.IP {
				url, err := nurl.JoinIPPort(nurl.Parts{Host: dev.Host, Port: strconv.Itoa(dev.Port)})
				if err != nil {
					inst.modbusErrorMsg("failed to validate device IP", url)
					netPollMan.PollingFinished(pp, pollStartTime, false, false, false, false, pollqueue.DELAYED_RETRY, callback)
					continue
				}
				mbClient.TCPClientHandler.Address = url
				mbClient.TCPClientHandler.SlaveID = byte(dev.AddressId)
			} else {
				inst.modbusDebugMsg(fmt.Sprintf("failed to validate device and network %v %s", err, dev.Name))
				netPollMan.PollingFinished(pp, pollStartTime, false, false, false, false, pollqueue.DELAYED_RETRY, callback)
				continue
			}

			var responseValue float64
			var response interface{}
			writeSuccess := false
			if writemode.IsWriteable(pnt.WriteMode) && boolean.IsTrue(pnt.WritePollRequired) { // DO WRITE IF REQUIRED
				inst.modbusDebugMsg(fmt.Sprintf("modbus write point: %+v", pnt))
				// pnt.PrintPointValues()
				if pnt.WriteValue != nil {
					response, responseValue, err = inst.networkWrite(mbClient, pnt)
					if err != nil {
						err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointWriteError)
						netPollMan.PollingFinished(pp, pollStartTime, false, false, true, false, pollqueue.DELAYED_RETRY, callback)
						continue
					}
					writeSuccess = true
					inst.modbusDebugMsg(fmt.Sprintf("modbus-write response: responseValue %f, point UUID: %s, response: %+v", responseValue, pnt.UUID, response))
				} else {
					writeSuccess = true // successful because there is no value to write.  Otherwise the point will short cycle.
					inst.modbusDebugMsg("modbus write point error: no value in priority array to write")
				}
			}

			readSuccess := false
			if boolean.IsTrue(pnt.ReadPollRequired) { // DO READ IF REQUIRED
				response, responseValue, err = inst.networkRead(mbClient, pnt)
				if err != nil {
					err = inst.pointUpdateErr(pnt, err.Error(), model.MessageLevel.Fail, model.CommonFaultCode.PointError)
					netPollMan.PollingFinished(pp, pollStartTime, false, false, true, false, pollqueue.DELAYED_RETRY, callback)
					continue
				}
				if boolean.IsTrue(pnt.IsBitwise) && pnt.BitwiseIndex != nil {

				}
				isChange := !float.ComparePtrValues(pnt.PresentValue, &responseValue)
				if isChange {
					if err != nil {
						netPollMan.PollingFinished(pp, pollStartTime, writeSuccess, readSuccess, true, false, pollqueue.NORMAL_RETRY, callback)
						continue
					}
				}
				readSuccess = true
				inst.modbusDebugMsg(fmt.Sprintf("modbus-read response: responseValue %f, point UUID: %s, response: %+v ", responseValue, pnt.UUID, response))
			}

			// update point in DB if required
			// For write_once and write_always type, write value should become present value
			writeValueToPresentVal := (pnt.WriteMode == model.WriteOnce || pnt.WriteMode == model.WriteAlways) && writeSuccess && pnt.WriteValue != nil

			if readSuccess || writeValueToPresentVal {
				if writeValueToPresentVal {
					responseValue = *pnt.WriteValue
					// fmt.Println("ModbusPolling: writeOnceWriteValueToPresentVal responseValue: ", responseValue)
					readSuccess = true
				}
				_, err = inst.pointUpdate(pnt, responseValue, readSuccess, true)
			}

			/*
				//JUST FOR TESTING
				pnt, err = inst.db.GetPoint(pp.FFPointUUID)
				if pnt == nil || err != nil {
					log.Errorf("modbus: AFTER... could not find pointID : %s\n", pp.FFPointUUID)
				}
			*/

			// This callback function triggers the PollManager to evaluate whether the point should be re-added to the PollQueue (Never, Immediately, or after the Poll Rate Delay)
			netPollMan.PollingFinished(pp, pollStartTime, writeSuccess, readSuccess, true, false, pollqueue.NORMAL_RETRY, callback)

		}
		return false, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	inst.pollingCancel = cancel
	go poll.GoPoll(ctx, f)
	return nil
}
