package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	pollqueue "github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/poll-queue"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils"
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

type devCheck struct {
	devUUID string
	client  Client
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

func (i *Instance) getNetworkPollManagerByUUID(netUUID string) (*pollqueue.NetworkPollManager, error) {
	for _, netPollMan := range i.NetworkPollManagers {
		if netPollMan.FFNetworkUUID == netUUID {
			return netPollMan, nil
		}
	}
	return nil, errors.New("modbus getNetworkPollManagerByUUID(): couldn't find NetworkPollManager")
}

var poll poller.Poller

//TODO: currently Polling loops through each network, grabs one point, and polls it.  Could be improved by having a seperate client/go routine for each of the networks.
func (i *Instance) ModbusPolling() error {
	poll = poller.New()
	var counter = 0
	f := func() (bool, error) {
		counter++
		//fmt.Println("\n \n")
		//modbusDebugMsg("LOOP COUNT: ", counter)
		var netArg api.Args
		/*
			nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, netArg)
			if err != nil {
				return false, err
			}
		*/

		if len(i.NetworkPollManagers) == 0 {
			modbusDebugMsg("NO MODBUS NETWORKS FOUND")
		}
		//modbusDebugMsg("i.NetworkPollManagers")
		//modbusDebugMsg("%+v\n", i.NetworkPollManagers)
		for _, netPollMan := range i.NetworkPollManagers { //LOOP THROUGH AND POLL NEXT POINTS IN EACH NETWORK QUEUE
			//modbusDebugMsg("ModbusPolling: netPollMan ", netPollMan.FFNetworkUUID)
			pollStartTime := time.Now()
			//Check that network exists
			//modbusDebugMsg("netPollMan")
			//modbusDebugMsg("%+v\n", netPollMan)
			net, err := i.db.GetNetwork(netPollMan.FFNetworkUUID, netArg)
			//modbusDebugMsg("net")
			//modbusDebugMsg("%+v\n", net)
			//modbusDebugMsg("err")
			//modbusDebugMsg("%+v\n", err)
			if err != nil || net == nil || net.PluginConfId != i.pluginUUID {
				modbusErrorMsg("MODBUS NETWORK NOT FOUND")
				continue
			}
			//modbusDebugMsg(fmt.Sprintf("modbus-poll: POLL START: NAME: %s\n", net.Name))

			if !utils.BoolIsNil(net.Enable) {
				modbusDebugMsg(fmt.Sprintf("NETWORK DISABLED: COUNT %v NAME: %s", counter, net.Name))
				continue
			}
			//netPollMan.PrintPollQueuePointUUIDs()
			pp, callback := netPollMan.GetNextPollingPoint() //callback function is called once polling is completed.
			//pp, _ := netPollMan.GetNextPollingPoint() //TODO: once polling completes, callback should be called
			if pp == nil {
				//modbusDebugMsg("No PollingPoint available in Network ", net.UUID)
				continue
			}
			if pp.FFNetworkUUID != net.UUID {
				modbusErrorMsg("PollingPoint FFNetworkUUID does not match the Network UUID")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			printPollingPointDebugInfo(pp)

			var devArg api.Args
			dev, err := i.db.GetDevice(pp.FFDeviceUUID, devArg)
			if dev == nil || err != nil {
				modbusErrorMsg("could not find deviceID:", pp.FFDeviceUUID)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			if !utils.BoolIsNil(dev.Enable) {
				modbusErrorMsg("device is disabled.")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			if dev.AddressId <= 0 || dev.AddressId >= 255 {
				modbusErrorMsg("address is not valid.  modbus addresses must be between 1 and 254")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			pnt, err := i.db.GetPoint(pp.FFPointUUID, api.Args{WithPriority: true})
			if pnt == nil || err != nil {
				modbusErrorMsg("could not find pointID: ", pp.FFPointUUID)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			//TODO: REPLACE WITH FUNCTION THAT PRINTS IMPORTANT POLLING INFORMATION
			printPointDebugInfo(pnt)

			if pnt.Priority == nil {
				modbusErrorMsg("ModbusPolling: HAD TO ADD PRIORITY ARRAY")
				pnt.Priority = &model.Priority{}
			}

			if !utils.BoolIsNil(pnt.Enable) {
				modbusErrorMsg("point is disabled.")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			modbusDebugMsg(fmt.Sprintf("MODBUS POLL! : Priority: %s, Network: %s Device: %s Point: %s Device-Add: %d Point-Add: %d Point Type: %s, WriteRequired: %t, ReadRequired: %t", pp.PollPriority, net.UUID, dev.UUID, pnt.UUID, dev.AddressId, *pnt.AddressID, pnt.ObjectType, utils.BoolIsNil(pnt.WritePollRequired), utils.BoolIsNil(pnt.ReadPollRequired)))

			if !utils.BoolIsNil(pnt.WritePollRequired) && !utils.BoolIsNil(pnt.ReadPollRequired) {
				modbusDebugMsg("polling not required on this point")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			SetPriorityArrayModeBasedOnWriteMode(pnt) //ensures the point PointPriorityArrayMode is set correctly

			// SETUP MODBUS CLIENT CONNECTION
			var mbClient smod.ModbusClient
			//var dCheck devCheck
			//dCheck.devUUID = dev.UUID
			mbClient, err = i.setClient(net, dev, true)
			if err != nil {
				modbusErrorMsg(fmt.Sprintf("failed to set client error: %v network name:%s", err, net.Name))
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			if net.TransportType == model.TransType.Serial || net.TransportType == model.TransType.LoRa {
				if dev.AddressId >= 1 {
					mbClient.RTUClientHandler.SlaveID = byte(dev.AddressId)
				}
			} else if dev.TransportType == model.TransType.IP {
				url, err := utils.JoinIPPort(utils.URLParts{model.TransType.IP, dev.Host, strconv.Itoa(dev.Port)})
				if err != nil {
					modbusErrorMsg("failed to validate device IP", url)
					netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
					continue
				}
				mbClient.TCPClientHandler.Address = url
				mbClient.TCPClientHandler.SlaveID = byte(dev.AddressId)
			} else {
				modbusDebugMsg(fmt.Sprintf("failed to validate device and network %v %s", err, dev.Name))
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			var responseValue float64
			var response interface{}
			var writeValuePointer *float64
			writeSuccess := false
			if isWriteable(pnt.WriteMode) && utils.BoolIsNil(pnt.WritePollRequired) { //DO WRITE IF REQUIRED
				modbusDebugMsg("modbus write point:")
				modbusDebugMsg("%+v", pnt)
				//pnt.PrintPointValues()
				writeValuePointer = pnt.Priority.GetHighestPriorityValue()
				if writeValuePointer != nil {
					response, responseValue, err = networkWrite(mbClient, pnt)
					if err != nil {
						_, err = i.pointUpdateErr(pnt, err)
						netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
						continue
					}
					writeSuccess = true
					modbusDebugMsg(fmt.Sprintf("modbus-write response: responseValue %f, point UUID: %s, response: %+v", responseValue, pnt.UUID, response))
				} else {
					writeSuccess = true //successful because there is no value to write.  Otherwise the point will short cycle.
					modbusDebugMsg("modbus write point error: no value in priority array to write")
				}
			}

			readSuccess := false
			if utils.BoolIsNil(pnt.ReadPollRequired) { //DO READ IF REQUIRED
				response, responseValue, err = networkRead(mbClient, pnt)
				if err != nil {
					_, err = i.pointUpdateErr(pnt, err)
					netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
					continue
				}
				//check cov
				isChange := !utils.CompareFloatPtr(pnt.PresentValue, &responseValue)
				if isChange {
					if err != nil {
						netPollMan.PollingFinished(pp, pollStartTime, writeSuccess, readSuccess, callback)
						continue
					}
				}
				readSuccess = true
				modbusDebugMsg(fmt.Sprintf("modbus-read response: responseValue %f, point UUID: %s, response: %+v ", responseValue, pnt.UUID, response))
			}

			//update point in DB if required
			//For write_once and write_always type, write value should become present value
			writeValueToPresentVal := (pnt.WriteMode == model.WriteOnce || pnt.WriteMode == model.WriteAlways) && writeSuccess && writeValuePointer != nil

			if readSuccess || writeValueToPresentVal {
				if writeValueToPresentVal {
					responseValue = *writeValuePointer
					//fmt.Println("ModbusPolling: writeOnceWriteValueToPresentVal responseValue: ", responseValue)
					readSuccess = true
				}
				_, err = i.pointUpdate(pnt, responseValue, writeSuccess, readSuccess, true)
			}

			/*
				//JUST FOR TESTING
				pnt, err = i.db.GetPoint(pp.FFPointUUID)
				if pnt == nil || err != nil {
					log.Errorf("modbus: AFTER... could not find pointID : %s\n", pp.FFPointUUID)
				}
			*/

			// This callback function triggers the PollManager to evaluate whether the point should be re-added to the PollQueue (Never, Immediately, or after the Poll Rate Delay)
			netPollMan.PollingFinished(pp, pollStartTime, writeSuccess, readSuccess, callback)

		}
		return false, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	i.pollingCancel = cancel
	go poll.GoPoll(ctx, f)
	return nil
}
