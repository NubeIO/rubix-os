package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/plugin/nube/projects/galvintmv/chirpstackrest"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/NubeIO/rubix-os/utils/integer"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// THE FOLLOWING GROUP OF FUNCTIONS ARE THE PLUGIN RESPONSES TO API CALLS FOR PLUGIN POINT, DEVICE, NETWORK (CRUD)
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	inst.tmvDebugMsg("addNetwork(): ", body.Name)
	network, err = inst.db.CreateNetwork(body)
	if err != nil {
		inst.tmvErrorMsg("addNetwork(): failed to create tmv network: ", body.Name)
		return nil, errors.New("failed to create tmv network")
	}

	if boolean.IsFalse(network.Enable) {
		err = inst.networkUpdateErr(network, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError)
		err = inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.NetworkError, true)
	}
	return network, nil
}

func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	if body == nil {
		inst.tmvDebugMsg("addDevice(): nil device object")
		return nil, errors.New("empty device body, no device created")
	}
	inst.tmvDebugMsg("addDevice(): ", body.Name)
	device, err = inst.db.CreateDevice(body)
	if device == nil || err != nil {
		inst.tmvDebugMsg("addDevice(): failed to create tmv device: ", body.Name)
		return nil, errors.New("failed to create tmv device")
	}

	inst.tmvDebugMsg("addDevice(): ", body.UUID)

	if boolean.IsFalse(device.Enable) {
		err = inst.deviceUpdateErr(device, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
		err = inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	}

	// NOTHING TO DO ON DEVICE CREATED
	return device, nil
}

func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	if body == nil {
		inst.tmvDebugMsg("addPoint(): nil point object")
		return nil, errors.New("empty point body, no point created")
	}
	inst.tmvDebugMsg("addPoint(): ", body.Name)

	point, err = inst.db.CreatePoint(body, true)
	if point == nil || err != nil {
		inst.tmvDebugMsg("addPoint(): failed to create tmv point: ", body.Name)
		return nil, errors.New("failed to create tmv point")
	}
	inst.tmvDebugMsg(fmt.Sprintf("addPoint(): %+v\n", point))

	if boolean.IsFalse(point.Enable) {
		err = inst.pointUpdateErr(point, "point disabled", model.MessageLevel.Warning, model.CommonFaultCode.PointError)
	}
	return point, nil

}

func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	inst.tmvDebugMsg("updateNetwork(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("updateNetwork():  nil network object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Warning
		body.CommonFault.MessageCode = model.CommonFaultCode.NetworkError
		body.CommonFault.Message = "network disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	} else {
		body.CommonFault.InFault = false
		body.CommonFault.MessageLevel = model.MessageLevel.Info
		body.CommonFault.MessageCode = model.CommonFaultCode.Ok
		body.CommonFault.Message = ""
		body.CommonFault.LastOk = time.Now().UTC()
	}

	network, err = inst.db.UpdateNetwork(body.UUID, body)
	if err != nil || network == nil {
		return nil, err
	}

	if boolean.IsFalse(network.Enable) {
		// DO POLLING DISABLE ACTIONS
		inst.db.SetErrorsForAllDevicesOnNetwork(network.UUID, "network disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError, true)
	}

	network, err = inst.db.UpdateNetwork(body.UUID, network)
	if err != nil || network == nil {
		return nil, err
	}
	return network, nil
}

func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	inst.tmvDebugMsg("updateDevice(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("updateDevice(): nil device object")
		return
	}

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Warning
		body.CommonFault.MessageCode = model.CommonFaultCode.DeviceError
		body.CommonFault.Message = "device disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	} else {
		body.CommonFault.InFault = false
		body.CommonFault.MessageLevel = model.MessageLevel.Info
		body.CommonFault.MessageCode = model.CommonFaultCode.Ok
		body.CommonFault.Message = ""
		body.CommonFault.LastOk = time.Now().UTC()
	}

	device, err = inst.db.UpdateDevice(body.UUID, body)
	if err != nil || device == nil {
		return nil, err
	}

	if boolean.IsFalse(device.Enable) {
		// DO POLLING DISABLE ACTIONS FOR DEVICE
		inst.db.SetErrorsForAllPointsOnDevice(device.UUID, "device disabled", model.MessageLevel.Warning, model.CommonFaultCode.DeviceError)
	} else {
		// TODO: Currently on every device update, all device points are removed, and re-added.
		device.CommonFault.InFault = false
		device.CommonFault.MessageLevel = model.MessageLevel.Info
		device.CommonFault.MessageCode = model.CommonFaultCode.Ok
		device.CommonFault.Message = ""
		device.CommonFault.LastOk = time.Now().UTC()
	}

	device, err = inst.db.UpdateDevice(device.UUID, device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	inst.tmvDebugMsg("updatePoint(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("updatePoint(): nil point object")
		return
	}

	inst.tmvDebugMsg(fmt.Sprintf("updatePoint() body: %+v\n", body))
	inst.tmvDebugMsg(fmt.Sprintf("updatePoint() priority: %+v\n", body.Priority))

	if boolean.IsFalse(body.Enable) {
		body.CommonFault.InFault = true
		body.CommonFault.MessageLevel = model.MessageLevel.Fail
		body.CommonFault.MessageCode = model.CommonFaultCode.PointError
		body.CommonFault.Message = "point disabled"
		body.CommonFault.LastFail = time.Now().UTC()
	}
	body.CommonFault.InFault = false
	body.CommonFault.MessageLevel = model.MessageLevel.Info
	body.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
	body.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
	body.CommonFault.LastOk = time.Now().UTC()
	point, err = inst.db.UpdatePoint(body.UUID, body)
	if err != nil || point == nil {
		inst.tmvDebugMsg("updatePoint(): bad response from UpdatePoint() err:", err)
		return nil, err
	}
	return point, nil
}

func (inst *Instance) writePoint(pntUUID string, body *model.PointWriter) (point *model.Point, err error) {
	// TODO: check for PointWriteByName calls that might not flow through the plugin.

	point = nil
	inst.tmvDebugMsg("writePoint(): ", pntUUID)
	if body == nil {
		inst.tmvDebugMsg("writePoint(): nil point object")
		return
	}

	inst.tmvDebugMsg(fmt.Sprintf("writePoint() body: %+v", body))
	inst.tmvDebugMsg(fmt.Sprintf("writePoint() priority: %+v", body.Priority))

	point, _, _, _, err = inst.db.PointWrite(pntUUID, body)
	if err != nil {
		inst.tmvDebugMsg("writePoint(): bad response from WritePoint(), ", err)
		return nil, err
	}

	return point, nil
}

func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	inst.tmvDebugMsg("deleteNetwork(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("deleteNetwork(): nil network object")
		return
	}

	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	inst.tmvDebugMsg("deleteDevice(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("deleteDevice(): nil device object")
		return
	}

	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	inst.tmvDebugMsg("deletePoint(): ", body.UUID)
	if body == nil {
		inst.tmvDebugMsg("deletePoint(): nil point object")
		return
	}

	ok, err = inst.db.DeletePoint(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// THE FOLLOWING FUNCTIONS ARE CALLED FROM WITHIN THE PLUGIN

func (inst *Instance) runSetupSteps() error {
	err := inst.createAndActivateChirpstackDevices()
	if err != nil {
		inst.tmvErrorMsg("runSetupSteps() createAndActivateChirpstackDevices() err:", err.Error())
	}
	err = inst.updatePointNames()
	if err != nil {
		inst.tmvErrorMsg("runSetupSteps() updatePointNames() err:", err.Error())
	}
	err = inst.createModbusNetworkDevicesAndPoints()
	if err != nil {
		inst.tmvErrorMsg("runSetupSteps() createModbusNetworkDevicesAndPoints() err:", err.Error())
	}
	err = inst.updateIOModuleRTC()
	if err != nil {
		inst.tmvErrorMsg("runSetupSteps() updateIOModuleRTC() err:", err.Error())
	}

	// CANCEL THE SETUP STEPS AFTER 12 HOURS
	if time.Unix(inst.startTime, 0).Add(12 * time.Hour).Before(time.Now()) {
		inst.tmvDebugMsg("STOP SETUP STEPS (AFTER 12 HOURS)")
		err = cron.RemoveByTag("runSetupSteps")
		if err != nil {
			inst.tmvErrorMsg("runSetupSteps() STOP SETUP STEPS err:", err.Error())
		}

		// WRITE THIS TO THE CONFIG YML ONCE BINOD SHOWS ME HOW
		inst.config.Job.EnableConfigSteps = false
		/*
			pluginConf, err := inst.db.GetPlugin(inst.pluginUUID)
			if err != nil {
				inst.tmvErrorMsg("checkComissioningPoints() GET PLUGIN CONF err:", err.Error())
			}

			inst.config.Job.EnableConfigSteps = false
			pluginConf.Config, err = yaml.Marshal(inst.config)
			if err != nil {
				inst.tmvErrorMsg("checkComissioningPoints() PARSE PLUGIN CONF TO BYTES err:", err.Error())
			}

			fmt.Println(fmt.Sprintf("PLUGIN CONFIG: %+v ", pluginConf))
			err = inst.db.DB.UpdatePluginConf(pluginConf)
			if err != nil {
				inst.tmvErrorMsg("checkComissioningPoints() UPDATE PLUGIN CONF err:", err.Error())
			}
		*/

	}
	return err
}

func (inst *Instance) checkComissioningPoints() error {
	// CREATE/UPDATE THE FLOW TEMPERATURE MODBUS POINTS (set them all disabled)
	inst.tmvDebugMsg("checkComissioningPoints() CREATE/UPDATE THE FLOW TEMPERATURE MODBUS POINTS (set them all disabled)")

	nets, err := inst.db.GetNetworksByPluginName("modbus", args.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}

	for _, net := range nets {
		inst.tmvDebugMsg("checkComissioningPoints() Net: ", net.Name)
		for _, dev := range net.Devices {
			inst.tmvDebugMsg("checkComissioningPoints() Dev: ", dev.Name)
			// modbus device exists now, so create the required points
			var foundFlowTempPoint *model.Point
			if dev.Points != nil {
				for _, modbusPoint := range dev.Points {
					if modbusPoint.Name == "FLOW_TEMPERATURE" {
						inst.tmvDebugMsg("checkComissioningPoints() FLOW_TEMPERATURE point found")
						foundFlowTempPoint = modbusPoint
						pointUpdateReq := false
						if boolean.NonNil(foundFlowTempPoint.Enable) != false {
							pointUpdateReq = true
							foundFlowTempPoint.Enable = boolean.NewFalse()
						}
						if integer.NonNil(foundFlowTempPoint.AddressID) != 1015 {
							pointUpdateReq = true
							foundFlowTempPoint.AddressID = integer.New(1015)
						}
						if foundFlowTempPoint.ObjectType != string(model.ObjTypeReadHolding) {
							pointUpdateReq = true
							foundFlowTempPoint.ObjectType = string(model.ObjTypeReadHolding)
						}
						if foundFlowTempPoint.WriteMode != model.ReadOnly {
							pointUpdateReq = true
							foundFlowTempPoint.WriteMode = model.ReadOnly
						}
						if foundFlowTempPoint.DataType != string(model.TypeFloat64) {
							pointUpdateReq = true
							foundFlowTempPoint.DataType = string(model.TypeFloat64)
						}
						if foundFlowTempPoint.PollRate != model.RATE_FAST {
							pointUpdateReq = true
							foundFlowTempPoint.PollRate = model.RATE_FAST
						}
						if foundFlowTempPoint.PollPriority != model.PRIORITY_HIGH {
							pointUpdateReq = true
							foundFlowTempPoint.PollPriority = model.PRIORITY_HIGH
						}
						if foundFlowTempPoint.HistoryType != model.HistoryTypeCovAndInterval {
							pointUpdateReq = true
							foundFlowTempPoint.HistoryType = model.HistoryTypeCovAndInterval
						}
						if foundFlowTempPoint.HistoryInterval == nil || *foundFlowTempPoint.HistoryInterval != 60 {
							pointUpdateReq = true
							foundFlowTempPoint.HistoryInterval = integer.New(60)
						}
						if pointUpdateReq {
							inst.tmvDebugMsg("checkComissioningPoints() updating FLOW_TEMPERATURE point")
							foundFlowTempPoint, err = inst.db.UpdatePoint(foundFlowTempPoint.UUID, foundFlowTempPoint)
							if err != nil {
								inst.tmvErrorMsg("checkComissioningPoints() FLOW_TEMPERATURE Point update err: ", err)
							}
							break
						}
					}
				}
				if foundFlowTempPoint == nil {
					foundFlowTempPoint = &model.Point{}
					foundFlowTempPoint.DeviceUUID = dev.UUID
					foundFlowTempPoint.Name = "FLOW_TEMPERATURE"
					foundFlowTempPoint.Enable = boolean.NewFalse()
					foundFlowTempPoint.AddressID = integer.New(1015)
					foundFlowTempPoint.ObjectType = string(model.ObjTypeReadHolding)
					foundFlowTempPoint.WriteMode = model.ReadOnly
					foundFlowTempPoint.DataType = string(model.TypeFloat64)
					foundFlowTempPoint.PollRate = model.RATE_FAST
					foundFlowTempPoint.PollPriority = model.PRIORITY_HIGH
					foundFlowTempPoint.Fallback = nil
					foundFlowTempPoint.Priority = &model.Priority{}
					foundFlowTempPoint.WritePollRequired = boolean.NewFalse()
					foundFlowTempPoint.HistoryType = model.HistoryTypeCovAndInterval
					foundFlowTempPoint.HistoryInterval = integer.New(60)
					foundFlowTempPoint, err = inst.db.CreatePoint(foundFlowTempPoint, true)
					if err != nil {
						inst.tmvErrorMsg("checkComissioningPoints() FLOW_TEMPERATURE Point create err: ", err)
					}
					time.Sleep(50 * time.Millisecond)
				}
			}
		}
	}

	// SET THE COMMISSIONING POINTS BACK TO NORMAL (no reads or writes) AFTER 12 HOURS
	if time.Unix(inst.startTime, 0).Add(12 * time.Hour).Before(time.Now()) {
		inst.DisableCommissioningPoints()

		err = cron.RemoveByTag("checkComissioningPoints")
		if err != nil {
			inst.tmvErrorMsg("checkComissioningPoints() STOP COMMISSIONING MODE err:", err.Error())
		}

		inst.config.Job.EnableCommissioning = false
		/*
			pluginConf, err := inst.db.GetPlugin(inst.pluginUUID)
			if err != nil {
				inst.tmvErrorMsg("checkComissioningPoints() GET PLUGIN CONF err:", err.Error())
			}

			inst.config.Job.EnableCommissioning = false
			pluginConf.Config, err = yaml.Marshal(inst.config)
			if err != nil {
				inst.tmvErrorMsg("checkComissioningPoints() PARSE PLUGIN CONF TO BYTES err:", err.Error())
			}

			fmt.Println(fmt.Sprintf("PLUGIN CONFIG: %+v ", pluginConf))
			err = inst.db.DB.UpdatePluginConf(pluginConf)
			if err != nil {
				inst.tmvErrorMsg("checkComissioningPoints() UPDATE PLUGIN CONF err:", err.Error())
			}
		*/
	}
	return err
}

func (inst *Instance) DisableCommissioningPoints() error {
	inst.tmvDebugMsg("DISABLE COMMISSIONING MODE: SET THE FLOW_TEMPERATURE MODBUS POINTS BACK TO NORMAL (disabled)")

	nets, err := inst.db.GetNetworksByPluginName("modbus", args.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}
	for _, net := range nets {
		inst.tmvDebugMsg("DisableCommissioningPoints() Net: ", net.Name)
		for _, dev := range net.Devices {
			for _, pnt := range dev.Points {
				if pnt.Name == "FLOW_TEMPERATURE" {
					inst.tmvDebugMsg("DisableCommissioningPoints() device: ", dev.Name)
					pnt.Enable = boolean.NewFalse()
					_, err = inst.db.UpdatePoint(pnt.UUID, pnt)
					if err != nil {
						inst.tmvErrorMsg("DisableCommissioningPoints() DISABLE FLOW_TEMPERATURE UpdatePoint() error: ", err)
					}
					time.Sleep(1 * time.Second)
				}
			}
		}
	}
	return nil
}

func (inst *Instance) updatePointNames() error {
	inst.tmvDebugMsg("updatePointNames()")
	nets, err := inst.db.GetNetworksByPluginName("lorawan", args.Args{WithDevices: true, WithPoints: true})
	// nets, err := inst.db.GetNetworksByPluginName("system", api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}
	for _, net := range nets {
		inst.tmvDebugMsg("updatePointNames() Net: ", net.Name)
		for _, dev := range net.Devices {
			for _, pnt := range dev.Points {
				newPointName := ""
				switch pnt.Name {
				case "digital-1":
					newPointName = "APP_FAULT"
				case "digital-2":
					newPointName = "FLOW_STATUS"
				case "temp-3":
					newPointName = "FLOW_TEMPERATURE"
				case "digital-4":
					newPointName = "OVER_TEMPERATURE_WARN"
				case "digital-5":
					newPointName = "OVER_TEMPERATURE_ALERT"
				case "int_8-6":
					newPointName = "DAILY_TEMP_TEST_1"
				case "int_8-7":
					newPointName = "DAILY_TEMP_TEST_2"
				case "int_8-8":
					newPointName = "DAILY_TEMP_TEST_3"
				case "int_8-9":
					newPointName = "MONTHLY_MEAN_TEMP_TEST"
				case "uint_32-10":
					newPointName = "TOTAL_FLOW_ACCUMULATION"
				case "uint_16-11":
					newPointName = "ONE_DAY_FLOW_ACCUMULATION"
				case "digital-12":
					newPointName = "ONE_DAY_LOW_FLOW_ALERT"
				case "uint_8-13":
					newPointName = "DAYS_OF_LOW_FLOW"
				case "digital-14":
					newPointName = "FIVE_DAY_LOW_FLOW_ALERT"
				case "digital-15":
					newPointName = "MONTHLY_HOT_FLUSH_STATUS"
				case "uint_8-16":
					newPointName = "OVER_TEMPERATURE_WARN_COUNT"
				case "uint_8-17":
					newPointName = "OVER_TEMPERATURE_ALERT_COUNT"
				case "digital-18":
					newPointName = "SOLENOID_STATUS"
				case "digital-19":
					newPointName = "ENABLE"
				case "temp-20":
					newPointName = "TEMPERATURE_SP"
				case "temp-21":
					newPointName = "OVER_TEMPERATURE_OFFSET"
				case "uint_16-22":
					newPointName = "LOW_FLOW_THRESHOLD"
				case "temp-23":
					newPointName = "HOT_FLUSH_SP"
				case "uint_16-24":
					newPointName = "HOT_FLUSH_DELAY"
				case "digital-25":
					newPointName = "RESET_ALL"
				case "digital-26":
					newPointName = "SOLENOID_ALLOW"
				case "uint_16-27":
					newPointName = "OVERTEMP_ALERT_DURATION_SP"
				case "temp-28":
					newPointName = "TEMP_CALIBRATION_OFFSET"
				}
				inst.tmvDebugMsg("Point Name: ", pnt.Name)
				if newPointName != "" {
					inst.tmvDebugMsg("NEW  Name: ", newPointName)
					pnt.Name = newPointName
					pnt, err = inst.db.UpdatePoint(pnt.UUID, pnt)
				}
			}
		}
	}
	return nil
}

func (inst *Instance) createModbusNetworkDevicesAndPoints() error {
	inst.tmvDebugMsg("createModbusNetworkDevicesAndPoints()")
	jsonFile, err := os.Open(inst.config.Job.TMVJSONFilePath)
	if err != nil {
		inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() file open err: ", err)
		return err
	}
	inst.tmvDebugMsg("createModbusNetworkDevicesAndPoints():  Successfully Opened ", inst.config.Job.TMVJSONFilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	// inst.tmvDebugMsg("createModbusDevicesAndPoints():  byteValue:", byteValue)

	var tmvDevices TMVDevices
	json.Unmarshal(byteValue, &tmvDevices.Devices)

	/*
		for ind, tmvDevice := range tmvDevices.Devices {
			inst.tmvDebugMsg(fmt.Sprintf("createModbusDevicesAndPoints() device %d: %+v", ind, tmvDevice))

		}
	*/
	modbusNet, err := inst.createModbusNetworkIfItNeeded("TMV")
	if err != nil || modbusNet == nil {
		inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() cannot get modbus network")
		return errors.New("createModbusNetworkDevicesAndPoints() cannot get modbus network")
	}

	nets, err := inst.db.GetNetworksByPluginName("lorawan", args.Args{WithDevices: true, WithPoints: true})
	// nets, err := inst.db.GetNetworksByPluginName("system", api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}

	for _, net := range nets {
		inst.tmvDebugMsg("createModbusNetworkDevicesAndPoints() Net: ", net.Name)
		for _, dev := range net.Devices {
			inst.tmvDebugMsg("createModbusNetworkDevicesAndPoints() lorawan network device: ", dev.Name)
			// Check for lorawan devices that are in the device json file
			for _, tmvDevice := range tmvDevices.Devices {
				if tmvDevice.DeviceName == dev.Name {
					inst.tmvDebugMsg("createModbusNetworkDevicesAndPoints() tmvDevice: ", tmvDevice.DeviceName)
					var foundModbusDevice *model.Device
					for _, modbusDevice := range modbusNet.Devices {
						if modbusDevice.Name == tmvDevice.DeviceName {
							foundModbusDevice = modbusDevice
							break
						}
					}
					if foundModbusDevice == nil {
						inst.tmvDebugMsg("createModbusNetworkDevicesAndPoints() no existing ModbusDevice: ")
						newDevice := model.Device{}
						newDevice.Name = tmvDevice.DeviceName
						newDevice.Enable = boolean.NewTrue()
						newDevice.AddressId = tmvDevice.DeviceAddress
						newDevice.ZeroMode = boolean.NewTrue()
						newDevice.NetworkUUID = modbusNet.UUID
						inst.tmvDebugMsg("createModbusNetworkDevicesAndPoints(): ", newDevice.Name)
						foundModbusDevice, err = inst.db.CreateDevice(&newDevice)
						if foundModbusDevice == nil || err != nil {
							inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints(): failed to create tmv device: ", newDevice.Name)
							continue
						}
						time.Sleep(50 * time.Millisecond)
					}
					// modbus device exists now, so create the required points
					var foundEnablePoint *model.Point
					var foundSetpointPoint *model.Point
					var foundResetPoint *model.Point
					var foundSolenoidAllowPoint *model.Point
					var foundCalibrationPoint *model.Point
					var foundRTCPoint *model.Point
					var foundRTCTZOffsetPoint *model.Point
					if foundModbusDevice.Points != nil {
						for _, modbusPoint := range foundModbusDevice.Points {
							switch modbusPoint.Name {
							case "ENABLE":
								foundEnablePoint = modbusPoint
								pointUpdateReq := false
								if boolean.NonNil(foundEnablePoint.Enable) != true {
									pointUpdateReq = true
									foundEnablePoint.Enable = boolean.NewTrue()
								}
								if integer.NonNil(foundEnablePoint.AddressID) != 1001 {
									pointUpdateReq = true
									foundEnablePoint.AddressID = integer.New(1001)
								}
								if foundEnablePoint.ObjectType != string(model.ObjTypeWriteCoil) {
									pointUpdateReq = true
									foundEnablePoint.ObjectType = string(model.ObjTypeWriteCoil)
								}
								if foundEnablePoint.WriteMode != model.WriteOnce {
									pointUpdateReq = true
									foundEnablePoint.WriteMode = model.WriteOnce
								}
								if foundEnablePoint.DataType != string(model.TypeDigital) {
									pointUpdateReq = true
									foundEnablePoint.DataType = string(model.TypeDigital)
								}
								if foundEnablePoint.PollRate != model.RATE_SLOW {
									pointUpdateReq = true
									foundEnablePoint.PollRate = model.RATE_SLOW
								}
								if foundEnablePoint.PollPriority != model.PRIORITY_HIGH {
									pointUpdateReq = true
									foundEnablePoint.PollPriority = model.PRIORITY_HIGH
								}
								if foundEnablePoint.HistoryType != model.HistoryTypeCovAndInterval {
									pointUpdateReq = true
									foundEnablePoint.HistoryType = model.HistoryTypeCovAndInterval
								}
								if foundEnablePoint.HistoryInterval == nil || *foundEnablePoint.HistoryInterval != 60 {
									pointUpdateReq = true
									foundEnablePoint.HistoryInterval = integer.New(60)
								}
								if pointUpdateReq {
									foundEnablePoint, err = inst.db.UpdatePoint(foundEnablePoint.UUID, foundEnablePoint)
									if err != nil {
										inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() EnablePoint update err: ", err)
									}
								}
							case "TEMPERATURE_SP":
								foundSetpointPoint = modbusPoint
								pointUpdateReq := false
								if boolean.NonNil(foundSetpointPoint.Enable) != true {
									pointUpdateReq = true
									foundSetpointPoint.Enable = boolean.NewTrue()
								}
								if integer.NonNil(foundSetpointPoint.AddressID) != 1001 {
									pointUpdateReq = true
									foundSetpointPoint.AddressID = integer.New(1001)
								}
								if foundSetpointPoint.ObjectType != string(model.ObjTypeWriteHolding) {
									pointUpdateReq = true
									foundSetpointPoint.ObjectType = string(model.ObjTypeWriteHolding)
								}
								if foundSetpointPoint.WriteMode != model.WriteOnce {
									pointUpdateReq = true
									foundSetpointPoint.WriteMode = model.WriteOnce
								}
								if foundSetpointPoint.DataType != string(model.TypeFloat64) {
									pointUpdateReq = true
									foundSetpointPoint.DataType = string(model.TypeFloat64)
								}
								if foundSetpointPoint.PollRate != model.RATE_SLOW {
									pointUpdateReq = true
									foundSetpointPoint.PollRate = model.RATE_SLOW
								}
								if foundSetpointPoint.PollPriority != model.PRIORITY_HIGH {
									pointUpdateReq = true
									foundSetpointPoint.PollPriority = model.PRIORITY_HIGH
								}
								if foundSetpointPoint.HistoryType != model.HistoryTypeCovAndInterval {
									pointUpdateReq = true
									foundSetpointPoint.HistoryType = model.HistoryTypeCovAndInterval
								}
								if foundSetpointPoint.HistoryInterval == nil || *foundSetpointPoint.HistoryInterval != 60 {
									pointUpdateReq = true
									foundSetpointPoint.HistoryInterval = integer.New(60)
								}
								if pointUpdateReq {
									foundSetpointPoint, err = inst.db.UpdatePoint(foundSetpointPoint.UUID, foundSetpointPoint)
									if err != nil {
										inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() SetpointPoint update err: ", err)
									}
								}
							case "SOLENOID_ALLOW":
								foundSolenoidAllowPoint = modbusPoint
								pointUpdateReq := false
								if boolean.NonNil(foundSolenoidAllowPoint.Enable) != true {
									pointUpdateReq = true
									foundSolenoidAllowPoint.Enable = boolean.NewTrue()
								}
								if integer.NonNil(foundSolenoidAllowPoint.AddressID) != 1003 {
									pointUpdateReq = true
									foundSolenoidAllowPoint.AddressID = integer.New(1003)
								}
								if foundSolenoidAllowPoint.ObjectType != string(model.ObjTypeWriteCoil) {
									pointUpdateReq = true
									foundSolenoidAllowPoint.ObjectType = string(model.ObjTypeWriteCoil)
								}
								if foundSolenoidAllowPoint.WriteMode != model.WriteOnce {
									pointUpdateReq = true
									foundSolenoidAllowPoint.WriteMode = model.WriteOnce
								}
								if foundSolenoidAllowPoint.DataType != string(model.TypeDigital) {
									pointUpdateReq = true
									foundSolenoidAllowPoint.DataType = string(model.TypeDigital)
								}
								if foundSolenoidAllowPoint.PollRate != model.RATE_SLOW {
									pointUpdateReq = true
									foundSolenoidAllowPoint.PollRate = model.RATE_SLOW
								}
								if foundSolenoidAllowPoint.PollPriority != model.PRIORITY_HIGH {
									pointUpdateReq = true
									foundSolenoidAllowPoint.PollPriority = model.PRIORITY_HIGH
								}
								if foundSolenoidAllowPoint.HistoryType != model.HistoryTypeCovAndInterval {
									pointUpdateReq = true
									foundSolenoidAllowPoint.HistoryType = model.HistoryTypeCovAndInterval
								}
								if foundSolenoidAllowPoint.HistoryInterval == nil || *foundSolenoidAllowPoint.HistoryInterval != 60 {
									pointUpdateReq = true
									foundSolenoidAllowPoint.HistoryInterval = integer.New(60)
								}
								if pointUpdateReq {
									foundSolenoidAllowPoint, err = inst.db.UpdatePoint(foundSolenoidAllowPoint.UUID, foundSolenoidAllowPoint)
									if err != nil {
										inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() SolenoidAllowPoint update err: ", err)
									}
								}
							case "TEMP_CALIBRATION_OFFSET":
								foundCalibrationPoint = modbusPoint
								pointUpdateReq := false
								if boolean.NonNil(foundCalibrationPoint.Enable) != true {
									pointUpdateReq = true
									foundCalibrationPoint.Enable = boolean.NewTrue()
								}
								if integer.NonNil(foundCalibrationPoint.AddressID) != 1013 {
									pointUpdateReq = true
									foundCalibrationPoint.AddressID = integer.New(1013)
								}
								if foundCalibrationPoint.ObjectType != string(model.ObjTypeWriteHolding) {
									pointUpdateReq = true
									foundCalibrationPoint.ObjectType = string(model.ObjTypeWriteHolding)
								}
								if foundCalibrationPoint.WriteMode != model.WriteOnce {
									pointUpdateReq = true
									foundCalibrationPoint.WriteMode = model.WriteOnce
								}
								if foundCalibrationPoint.DataType != string(model.TypeFloat64) {
									pointUpdateReq = true
									foundCalibrationPoint.DataType = string(model.TypeFloat64)
								}
								if foundCalibrationPoint.PollRate != model.RATE_SLOW {
									pointUpdateReq = true
									foundCalibrationPoint.PollRate = model.RATE_SLOW
								}
								if foundCalibrationPoint.PollPriority != model.PRIORITY_HIGH {
									pointUpdateReq = true
									foundCalibrationPoint.PollPriority = model.PRIORITY_HIGH
								}
								if foundCalibrationPoint.HistoryType != model.HistoryTypeCovAndInterval {
									pointUpdateReq = true
									foundCalibrationPoint.HistoryType = model.HistoryTypeCovAndInterval
								}
								if foundCalibrationPoint.HistoryInterval == nil || *foundCalibrationPoint.HistoryInterval != 60 {
									pointUpdateReq = true
									foundCalibrationPoint.HistoryInterval = integer.New(60)
								}
								if pointUpdateReq {
									foundCalibrationPoint, err = inst.db.UpdatePoint(foundCalibrationPoint.UUID, foundCalibrationPoint)
									if err != nil {
										inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() CalibrationPoint update err: ", err)
									}
								}
							case "RESET_ALL":
								foundResetPoint = modbusPoint
								pointUpdateReq := false
								if boolean.NonNil(foundResetPoint.Enable) != true {
									pointUpdateReq = true
									foundResetPoint.Enable = boolean.NewTrue()
								}
								if integer.NonNil(foundResetPoint.AddressID) != 1002 {
									pointUpdateReq = true
									foundResetPoint.AddressID = integer.New(1002)
								}
								if foundResetPoint.ObjectType != string(model.ObjTypeWriteCoil) {
									pointUpdateReq = true
									foundResetPoint.ObjectType = string(model.ObjTypeWriteCoil)
								}
								if foundResetPoint.WriteMode != model.WriteOnce {
									pointUpdateReq = true
									foundResetPoint.WriteMode = model.WriteOnce
								}
								if foundResetPoint.DataType != string(model.TypeDigital) {
									pointUpdateReq = true
									foundResetPoint.DataType = string(model.TypeDigital)
								}
								if foundResetPoint.PollRate != model.RATE_SLOW {
									pointUpdateReq = true
									foundResetPoint.PollRate = model.RATE_SLOW
								}
								if foundResetPoint.PollPriority != model.PRIORITY_HIGH {
									pointUpdateReq = true
									foundResetPoint.PollPriority = model.PRIORITY_HIGH
								}
								if foundResetPoint.HistoryType != model.HistoryTypeCovAndInterval {
									pointUpdateReq = true
									foundResetPoint.HistoryType = model.HistoryTypeCovAndInterval
								}
								if foundResetPoint.HistoryInterval == nil || *foundResetPoint.HistoryInterval != 60 {
									pointUpdateReq = true
									foundResetPoint.HistoryInterval = integer.New(60)
								}
								if pointUpdateReq {
									foundResetPoint, err = inst.db.UpdatePoint(foundResetPoint.UUID, foundResetPoint)
									if err != nil {
										inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() ResetPoint update err: ", err)
									}
								}
							case "RTC":
								foundRTCPoint = modbusPoint
								pointUpdateReq := false
								if boolean.NonNil(foundRTCPoint.Enable) != true {
									pointUpdateReq = true
									foundRTCPoint.Enable = boolean.NewTrue()
								}
								if integer.NonNil(foundRTCPoint.AddressID) != 10007 {
									pointUpdateReq = true
									foundRTCPoint.AddressID = integer.New(10007)
								}
								if foundRTCPoint.ObjectType != string(model.ObjTypeWriteHolding) {
									pointUpdateReq = true
									foundRTCPoint.ObjectType = string(model.ObjTypeWriteHolding)
								}
								if foundRTCPoint.WriteMode != model.WriteOnce {
									pointUpdateReq = true
									foundRTCPoint.WriteMode = model.WriteOnce
								}
								if foundRTCPoint.DataType != string(model.TypeUint32) {
									pointUpdateReq = true
									foundRTCPoint.DataType = string(model.TypeUint32)
								}
								if foundRTCPoint.PollRate != model.RATE_SLOW {
									pointUpdateReq = true
									foundRTCPoint.PollRate = model.RATE_SLOW
								}
								if foundRTCPoint.PollPriority != model.PRIORITY_HIGH {
									pointUpdateReq = true
									foundRTCPoint.PollPriority = model.PRIORITY_HIGH
								}
								if foundRTCPoint.HistoryType != model.HistoryTypeCovAndInterval {
									pointUpdateReq = true
									foundRTCPoint.HistoryType = model.HistoryTypeCovAndInterval
								}
								if foundRTCPoint.HistoryInterval == nil || *foundRTCPoint.HistoryInterval != 60 {
									pointUpdateReq = true
									foundRTCPoint.HistoryInterval = integer.New(60)
								}
								if pointUpdateReq {
									foundRTCPoint, err = inst.db.UpdatePoint(foundRTCPoint.UUID, foundRTCPoint)
									if err != nil {
										inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() RTCPoint update err: ", err)
									}
								}
							case "RTC_TZ_OFFSET":
								foundRTCTZOffsetPoint = modbusPoint
								pointUpdateReq := false
								if boolean.NonNil(foundRTCTZOffsetPoint.Enable) != true {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.Enable = boolean.NewTrue()
								}
								if integer.NonNil(foundRTCTZOffsetPoint.AddressID) != 10009 {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.AddressID = integer.New(10009)
								}
								if foundRTCTZOffsetPoint.ObjectType != string(model.ObjTypeWriteHolding) {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.ObjectType = string(model.ObjTypeWriteHolding)
								}
								if foundRTCTZOffsetPoint.WriteMode != model.WriteOnce {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.WriteMode = model.WriteOnce
								}
								if foundRTCTZOffsetPoint.DataType != string(model.TypeUint32) {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.DataType = string(model.TypeUint32)
								}
								if foundRTCTZOffsetPoint.PollRate != model.RATE_SLOW {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.PollRate = model.RATE_SLOW
								}
								if foundRTCTZOffsetPoint.PollPriority != model.PRIORITY_HIGH {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.PollPriority = model.PRIORITY_HIGH
								}
								if foundRTCTZOffsetPoint.HistoryType != model.HistoryTypeCovAndInterval {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.HistoryType = model.HistoryTypeCovAndInterval
								}
								if foundRTCTZOffsetPoint.HistoryInterval == nil || *foundRTCTZOffsetPoint.HistoryInterval != 60 {
									pointUpdateReq = true
									foundRTCTZOffsetPoint.HistoryInterval = integer.New(60)
								}
								if pointUpdateReq {
									foundRTCTZOffsetPoint, err = inst.db.UpdatePoint(foundRTCTZOffsetPoint.UUID, foundRTCTZOffsetPoint)
									if err != nil {
										inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() RTCPoint update err: ", err)
									}
								}
							}
							time.Sleep(50 * time.Millisecond)
						}
					}
					if foundEnablePoint == nil {
						foundEnablePoint = &model.Point{}
						foundEnablePoint.DeviceUUID = foundModbusDevice.UUID
						foundEnablePoint.Name = "ENABLE"
						foundEnablePoint.Enable = boolean.NewTrue()
						foundEnablePoint.AddressID = integer.New(1001)
						foundEnablePoint.ObjectType = string(model.ObjTypeWriteCoil)
						foundEnablePoint.WriteMode = model.WriteOnce
						foundEnablePoint.DataType = string(model.TypeDigital)
						foundEnablePoint.PollRate = model.RATE_SLOW
						foundEnablePoint.PollPriority = model.PRIORITY_HIGH
						foundEnablePoint.Fallback = float.New(1)
						foundEnablePoint.Priority = &model.Priority{}
						foundEnablePoint.Priority.P3 = float.New(1)
						foundEnablePoint.WritePollRequired = boolean.NewTrue()
						foundEnablePoint.HistoryType = model.HistoryTypeCovAndInterval
						foundEnablePoint.HistoryInterval = integer.New(60)
						foundEnablePoint, err = inst.db.CreatePoint(foundEnablePoint, true)
						if err != nil {
							inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() EnablePoint create err: ", err)
						}
						time.Sleep(50 * time.Millisecond)
					}
					if foundSetpointPoint == nil {
						foundSetpointPoint = &model.Point{}
						foundSetpointPoint.DeviceUUID = foundModbusDevice.UUID
						foundSetpointPoint.Name = "TEMPERATURE_SP"
						foundSetpointPoint.Enable = boolean.NewTrue()
						foundSetpointPoint.AddressID = integer.New(1001)
						foundSetpointPoint.ObjectType = string(model.ObjTypeWriteHolding)
						foundSetpointPoint.WriteMode = model.WriteOnce
						foundSetpointPoint.DataType = string(model.TypeFloat64)
						foundSetpointPoint.PollRate = model.RATE_SLOW
						foundSetpointPoint.PollPriority = model.PRIORITY_HIGH
						foundSetpointPoint.Fallback = float.New(tmvDevice.TemperatureSetpoint)
						foundSetpointPoint.Priority = &model.Priority{}
						foundSetpointPoint.Priority.P3 = float.New(tmvDevice.TemperatureSetpoint)
						foundSetpointPoint.WritePollRequired = boolean.NewTrue()
						foundSetpointPoint.HistoryType = model.HistoryTypeCovAndInterval
						foundSetpointPoint.HistoryInterval = integer.New(60)
						foundSetpointPoint, err = inst.db.CreatePoint(foundSetpointPoint, true)
						if err != nil {
							inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() SetpointPoint create err: ", err)
						}
						time.Sleep(50 * time.Millisecond)
					}
					if foundResetPoint == nil {
						foundResetPoint = &model.Point{}
						foundResetPoint.DeviceUUID = foundModbusDevice.UUID
						foundResetPoint.Name = "RESET_ALL"
						foundResetPoint.Enable = boolean.NewTrue()
						foundResetPoint.AddressID = integer.New(1002)
						foundResetPoint.ObjectType = string(model.ObjTypeWriteCoil)
						foundResetPoint.WriteMode = model.WriteOnce
						foundResetPoint.DataType = string(model.TypeDigital)
						foundResetPoint.PollRate = model.RATE_SLOW
						foundResetPoint.PollPriority = model.PRIORITY_HIGH
						foundResetPoint.Fallback = float.New(0)
						foundResetPoint.Priority = &model.Priority{}
						foundResetPoint.Priority.P3 = float.New(0)
						foundResetPoint.WritePollRequired = boolean.NewFalse()
						foundResetPoint.HistoryType = model.HistoryTypeCovAndInterval
						foundResetPoint.HistoryInterval = integer.New(60)
						foundResetPoint, err = inst.db.CreatePoint(foundResetPoint, true)
						if err != nil {
							inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() ResetPoint create err: ", err)
						}
						time.Sleep(50 * time.Millisecond)
					}
					if foundSolenoidAllowPoint == nil {
						foundSolenoidAllowPoint = &model.Point{}
						foundSolenoidAllowPoint.DeviceUUID = foundModbusDevice.UUID
						foundSolenoidAllowPoint.Name = "SOLENOID_ALLOW"
						foundSolenoidAllowPoint.Enable = boolean.NewTrue()
						foundSolenoidAllowPoint.AddressID = integer.New(1003)
						foundSolenoidAllowPoint.ObjectType = string(model.ObjTypeWriteCoil)
						foundSolenoidAllowPoint.WriteMode = model.WriteOnce
						foundSolenoidAllowPoint.DataType = string(model.TypeDigital)
						foundSolenoidAllowPoint.PollRate = model.RATE_SLOW
						foundSolenoidAllowPoint.PollPriority = model.PRIORITY_HIGH
						fallbackFloat := float64(0)
						if tmvDevice.SolenoidRequired == "Yes" || tmvDevice.SolenoidRequired == "true" {
							fallbackFloat = 1
						}
						foundSolenoidAllowPoint.Fallback = float.New(fallbackFloat)
						foundSolenoidAllowPoint.Priority = &model.Priority{}
						foundSolenoidAllowPoint.Priority.P3 = float.New(fallbackFloat)
						foundSolenoidAllowPoint.WritePollRequired = boolean.NewTrue()
						foundSolenoidAllowPoint.HistoryType = model.HistoryTypeCovAndInterval
						foundSolenoidAllowPoint.HistoryInterval = integer.New(60)
						foundSolenoidAllowPoint, err = inst.db.CreatePoint(foundSolenoidAllowPoint, true)
						if err != nil {
							inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() SolenoidAllowPoint create err: ", err)
						}
						time.Sleep(50 * time.Millisecond)
					}
					if foundCalibrationPoint == nil {
						foundCalibrationPoint = &model.Point{}
						foundCalibrationPoint.DeviceUUID = foundModbusDevice.UUID
						foundCalibrationPoint.Name = "TEMP_CALIBRATION_OFFSET"
						foundCalibrationPoint.Enable = boolean.NewTrue()
						foundCalibrationPoint.AddressID = integer.New(1013)
						foundCalibrationPoint.ObjectType = string(model.ObjTypeWriteHolding)
						foundCalibrationPoint.WriteMode = model.WriteOnce
						foundCalibrationPoint.DataType = string(model.TypeFloat64)
						foundCalibrationPoint.PollRate = model.RATE_SLOW
						foundCalibrationPoint.PollPriority = model.PRIORITY_HIGH
						foundCalibrationPoint.Fallback = float.New(0)
						foundCalibrationPoint.Priority = &model.Priority{}
						foundCalibrationPoint.Priority.P3 = float.New(0)
						foundCalibrationPoint.WritePollRequired = boolean.NewFalse()
						foundCalibrationPoint.HistoryType = model.HistoryTypeCovAndInterval
						foundCalibrationPoint.HistoryInterval = integer.New(60)
						foundCalibrationPoint, err = inst.db.CreatePoint(foundCalibrationPoint, true)
						if err != nil {
							inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() CalibrationPoint create err: ", err)
						}
						time.Sleep(50 * time.Millisecond)
					}
					if foundRTCPoint == nil {
						foundRTCPoint = &model.Point{}
						foundRTCPoint.DeviceUUID = foundModbusDevice.UUID
						foundRTCPoint.Name = "RTC"
						foundRTCPoint.Enable = boolean.NewTrue()
						foundRTCPoint.AddressID = integer.New(10007)
						foundRTCPoint.ObjectType = string(model.ObjTypeWriteHolding)
						foundRTCPoint.WriteMode = model.WriteOnce
						foundRTCPoint.DataType = string(model.TypeUint32)
						foundRTCPoint.PollRate = model.RATE_SLOW
						foundRTCPoint.PollPriority = model.PRIORITY_HIGH
						foundRTCPoint.WritePollRequired = boolean.NewFalse()
						foundRTCPoint.HistoryType = model.HistoryTypeCovAndInterval
						foundRTCPoint.HistoryInterval = integer.New(60)
						foundRTCPoint, err = inst.db.CreatePoint(foundRTCPoint, true)
						if err != nil {
							inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() RTCPoint create err: ", err)
						}
						time.Sleep(50 * time.Millisecond)
					}
					if foundRTCTZOffsetPoint == nil {
						foundRTCTZOffsetPoint = &model.Point{}
						foundRTCTZOffsetPoint.DeviceUUID = foundModbusDevice.UUID
						foundRTCTZOffsetPoint.Name = "RTC_TZ_OFFSET"
						foundRTCTZOffsetPoint.Enable = boolean.NewTrue()
						foundRTCTZOffsetPoint.AddressID = integer.New(10009)
						foundRTCTZOffsetPoint.ObjectType = string(model.ObjTypeWriteHolding)
						foundRTCTZOffsetPoint.WriteMode = model.WriteOnce
						foundRTCTZOffsetPoint.DataType = string(model.TypeUint32)
						foundRTCTZOffsetPoint.PollRate = model.RATE_SLOW
						foundRTCTZOffsetPoint.PollPriority = model.PRIORITY_HIGH
						foundRTCTZOffsetPoint.WritePollRequired = boolean.NewFalse()
						foundRTCTZOffsetPoint.HistoryType = model.HistoryTypeCovAndInterval
						foundRTCTZOffsetPoint.HistoryInterval = integer.New(60)
						foundRTCTZOffsetPoint, err = inst.db.CreatePoint(foundRTCTZOffsetPoint, true)
						if err != nil {
							inst.tmvErrorMsg("createModbusNetworkDevicesAndPoints() RTCTZOffsetPoint create err: ", err)
						}
						time.Sleep(50 * time.Millisecond)
					}
				}
			}
			/*
				for _, pnt := range dev.Points {

				}
			*/
		}
	}
	return nil
}

func (inst *Instance) createModbusNetworkIfItNeeded(reqNetName string) (*model.Network, error) {
	inst.tmvDebugMsg("createModbusNetworkIfItNeeded(): reqNetName")

	modbusNetworks, err := inst.db.GetNetworksByPluginName("modbus", args.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return nil, err
	}
	for _, modbusNet := range modbusNetworks {
		inst.tmvDebugMsg("createModbusDevicesAndPoints() modbusNet: ", modbusNet.Name)
		if modbusNet.Name == reqNetName {
			return modbusNet, nil
		}
	}
	// not found, so create a new FF modbus network
	newModbusNet := model.Network{}
	newModbusNet.PluginPath = "modbus"
	newModbusNet.Name = reqNetName
	newModbusNet.Enable = boolean.NewTrue()
	serialPort := "/data/socat/ptyLORAWAN-1"
	newModbusNet.SerialPort = &serialPort
	newModbusNet.SerialTimeout = integer.New(8)
	newModbusNet.MaxPollRate = float.New(10)
	newModbusNet.TransportType = "serial"
	return inst.db.CreateNetwork(&newModbusNet)
}

func (inst *Instance) createAndActivateChirpstackDevices() error {
	inst.tmvDebugMsg("createAndActivateChirpstackDevices()")
	jsonFile, err := os.Open(inst.config.Job.TMVJSONFilePath)
	if err != nil {
		inst.tmvErrorMsg("createAndActivateChirpstackDevices() file open err: ", err)
		return err
	}
	inst.tmvDebugMsg("createAndActivateChirpstackDevices():  Successfully Opened ", inst.config.Job.TMVJSONFilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var tmvDevices TMVDevices
	json.Unmarshal(byteValue, &tmvDevices.Devices)

	token := &chirpstackrest.ChirpstackToken{Token: inst.config.Job.ChirpstackToken}
	if token != nil && token.Token != "" {
		inst.tmvDebugMsg("createAndActivateChirpstackDevices() token: ", token.Token)
		profileUUID, err := inst.GetChirpstackDeviceProfileUUID(token.Token)
		if profileUUID == "" || err != nil {
			inst.tmvErrorMsg("createAndActivateChirpstackDevices() err: ", err)
			return err
		}
		inst.tmvDebugMsg("createAndActivateChirpstackDevices() device profile UUID: ", profileUUID)

		mapFileContent := ""
		for _, tmvDevice := range tmvDevices.Devices {
			err := inst.AddChirpstackDevice(inst.config.Job.ChirpstackApplicationNumber, tmvDevice.DeviceAddress, tmvDevice.DeviceName, tmvDevice.LoRaWANDeviceEUI, profileUUID, token.Token)
			// err := inst.AddChirpstackDevice(inst.config.Job.ChirpstackApplicationNumber, 666, "Test21", "4E7562654910FFFF", profileUUID, token.Token)
			if err != nil {
				inst.tmvErrorMsg("createAndActivateChirpstackDevices() AddChirpstackDevice() error: ", err)
			}
			time.Sleep(50 * time.Millisecond)
			err = inst.ActivateChirpstackDevice(inst.config.Job.ChirpstackNetworkKey, tmvDevice.LoRaWANDeviceEUI, token.Token, inst.config.Job.ChirpstackNetworkKey)
			if err != nil {
				inst.tmvErrorMsg("createAndActivateChirpstackDevices() ActivateChirpstackDevice() error: ", err)
			}
			time.Sleep(50 * time.Millisecond)
			mapFileContent += "\"" + strconv.Itoa(tmvDevice.DeviceAddress) + ":" + tmvDevice.LoRaWANDeviceEUI + "\"," + "\n"

		}
		// err = os.Remove(inst.config.Job.LorawanBridgeMapFilePath)
		err = os.Remove("/home/pi/lorawan-modbus-bridge/map.txt")
		if err != nil {
			inst.tmvErrorMsg("createAndActivateChirpstackDevices() remove old map file error: ", err)
		}
		// err = os.WriteFile(inst.config.Job.LorawanBridgeMapFilePath, []byte(mapFileContent), 0666)
		err = os.WriteFile("/data/lorawan-modbus-bridge/map.txt", []byte(mapFileContent), 0666)
		if err != nil {
			inst.tmvErrorMsg("createAndActivateChirpstackDevices() create new map file error: ", err)
		}
	}
	return nil
}

func (inst *Instance) updateIOModuleRTC() error {
	inst.tmvDebugMsg("updateIOModuleRTC()")
	nets, err := inst.db.GetNetworksByPluginName("modbus", args.Args{WithDevices: true, WithPoints: true})
	// nets, err := inst.db.GetNetworksByPluginName("system", api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return err
	}
	for _, net := range nets {
		inst.tmvDebugMsg("updateIOModuleRTC() Net: ", net.Name)
		for _, dev := range net.Devices {
			for _, pnt := range dev.Points {
				if pnt.Name == "RTC" {
					inst.tmvDebugMsg("updateIOModuleRTC() rtcPointWrite() device: ", dev.Name)
					inst.rtcPointWrite(pnt)
					time.Sleep(30 * time.Second)
				} else if pnt.Name == "RTC_TZ_OFFSET" {
					inst.tmvDebugMsg("updateIOModuleRTC() rtcTZOffsetPointWrite() device: ", dev.Name)
					inst.rtcTZOffsetPointWrite(pnt)
					time.Sleep(30 * time.Second)
				}
			}
		}
	}
	return nil
}

func (inst *Instance) rtcPointWrite(rtcPoint *model.Point) error {
	now := time.Now().Unix()
	inst.tmvDebugMsg("rtcPointWrite() now: ", now)
	pointUpdateMap := make(map[string]*float64)
	pointUpdateMap["_1"] = float.New(float64(now))
	pointWriter := model.PointWriter{Priority: &pointUpdateMap}

	_, _, _, _, err := inst.db.PointWrite(rtcPoint.UUID, &pointWriter)
	if err != nil {
		inst.tmvErrorMsg("rtcPointWrite() PointWrite() error: ", err)
		return err
	}
	return nil
}

func (inst *Instance) rtcTZOffsetPointWrite(rtcTZOffsetPoint *model.Point) error {
	_, offset := time.Now().Local().Zone()
	inst.tmvDebugMsg("rtcTZOffsetPointWrite() offset: ", offset)
	if rtcTZOffsetPoint.PresentValue != nil && *rtcTZOffsetPoint.PresentValue == float64(offset) {
		return nil
	}
	pointUpdateMap := make(map[string]*float64)
	pointUpdateMap["_1"] = float.New(float64(offset))
	pointWriter := model.PointWriter{Priority: &pointUpdateMap}

	_, _, _, _, err := inst.db.PointWrite(rtcTZOffsetPoint.UUID, &pointWriter)
	if err != nil {
		inst.tmvErrorMsg("rtcTZOffsetPointWrite() PointWrite() error: ", err)
		return err
	}
	return nil
}

func (inst *Instance) pointUpdateErr(point *model.Point, message string, messageLevel string, messageCode string) error {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = messageLevel
	point.CommonFault.MessageCode = messageCode
	point.CommonFault.Message = message
	point.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdatePointErrors(point.UUID, point)
	if err != nil {
		inst.tmvErrorMsg(" pointUpdateErr()", err)
	}
	return err
}

func (inst *Instance) deviceUpdateErr(device *model.Device, message string, messageLevel string, messageCode string) error {
	device.CommonFault.InFault = true
	device.CommonFault.MessageLevel = messageLevel
	device.CommonFault.MessageCode = messageCode
	device.CommonFault.Message = message
	device.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdateDeviceErrors(device.UUID, device)
	if err != nil {
		inst.tmvErrorMsg(" deviceUpdateErr()", err)
	}
	return err
}

func (inst *Instance) networkUpdateErr(network *model.Network, message string, messageLevel string, messageCode string) error {
	network.CommonFault.InFault = true
	network.CommonFault.MessageLevel = messageLevel
	network.CommonFault.MessageCode = messageCode
	network.CommonFault.Message = message
	network.CommonFault.LastFail = time.Now().UTC()
	err := inst.db.UpdateNetworkErrors(network.UUID, network)
	if err != nil {
		inst.tmvErrorMsg(" networkUpdateErr()", err)
	}
	return err
}
