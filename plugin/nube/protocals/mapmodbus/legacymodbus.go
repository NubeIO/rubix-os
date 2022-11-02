package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/mapmodbus/legacymodbusrest"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *Instance) GetLegacyModbusNetworksAndDevices() (*[]legacymodbusrest.ModbusNet, error) {
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1919
	}
	rest := legacymodbusrest.NewNoAuth(host, int(port))
	modbusNet, err := rest.GetLegacyModbusNetworksAndDevices()
	if err != nil || modbusNet == nil {
		inst.mapmodbusErrorMsg("no legacy modbus network found. err: ", err)
		return nil, errors.New(fmt.Sprint("no legacy modbus network found. err:", err))
	}
	return modbusNet, nil
}

func (inst *Instance) GetLegacyModbusDeviceAndPoints(modbusDevUUID string) (*legacymodbusrest.ModbusDev, error) {
	inst.mapmodbusErrorMsg("GetLegacyModbusDeviceAndPoints() modbusDevUUID: ", modbusDevUUID)
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1616
	}
	rest := legacymodbusrest.NewNoAuth(host, int(port))
	modbusDevsArray, err := rest.GetLegacyModbusDeviceAndPoints(modbusDevUUID)
	if err != nil {
		return nil, errors.New("no legacy modbus devices found")
	}
	return modbusDevsArray, nil
}

func (inst *Instance) ConvertLegacyModbusPropsToFFModbusPoints(functionCode, dataTypeLegacy string) (objectType, dataType string, writeMode model.WriteMode, dataLength int, writeable bool) {
	writeable = false
	writeMode = model.ReadOnly
	objectType = ""
	dataLength = 1
	switch functionCode {
	case "READ_COILS":
		objectType = "read_coil"
		writeMode = model.ReadOnly
		writeable = false
	case "READ_DISCRETE_INPUTS":
		objectType = "read_discrete_input"
		writeMode = model.ReadOnly
		writeable = false
	case "READ_HOLDING_REGISTERS":
		objectType = "read_holding"
		writeMode = model.ReadOnly
		writeable = false
	case "READ_INPUT_REGISTERS":
		objectType = "read_register"
		writeMode = model.ReadOnly
		writeable = false
	case "WRITE_COIL":
		objectType = "write_coil"
		writeMode = model.WriteAndMaintain
		writeable = true
	case "WRITE_REGISTER":
		objectType = "write_holding"
		writeMode = model.WriteAndMaintain
		writeable = true
	case "WRITE_COILS":
		objectType = "write_coil"
		writeMode = model.WriteAndMaintain
		writeable = true
	case "WRITE_REGISTERS":
		objectType = "write_holding"
		writeMode = model.WriteAndMaintain
		writeable = true
	default:
		inst.mapmodbusErrorMsg("unrecognized modbus object type / function code")
	}

	dataType = ""
	switch dataTypeLegacy {
	case "RAW":
		dataType = "int16"
		dataLength = 1
	case "INT16":
		dataType = "int16"
		dataLength = 1
	case "UINT16":
		dataType = "uint16"
		dataLength = 1
	case "INT32":
		dataType = "int32"
		dataLength = 2
	case "UINT32":
		dataType = "uint32"
		dataLength = 2
	case "FLOAT":
		dataType = "float32"
		dataLength = 2
	case "DOUBLE":
		dataType = "int32"
		dataLength = 2
	case "DIGITAL":
		dataType = "digital"
		dataLength = 1
	default:
		inst.mapmodbusErrorMsg("unrecognized modbus data type")
	}
	return objectType, dataType, writeMode, dataLength, writeable
}
