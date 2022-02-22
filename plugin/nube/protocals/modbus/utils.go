package main

import (
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

const (
	readBool uint = iota + 1
	readUint16
	readInt16
	readUint32
	readInt32
	readFloat32
	readUint64
	readInt64
	readFloat64
	readCoil
	writeCoil
	writeCoils
	writeUint16
	writeInt16
	writeInt32
	writeUint32
	writeFloat32
	writeInt64
	writeUint64
	writeFloat64
)

func isWrite(t string) bool {
	switch model.ObjectType(t) {
	case model.ObjTypeWriteCoil, model.ObjTypeWriteCoils:
		return true
	case model.ObjTypeWriteHolding, model.ObjTypeWriteHoldings:
		return true
	case model.ObjTypeWriteInt16, model.ObjTypeWriteUint16:
		return true
	case model.ObjTypeWriteFloat32:
		return true
	}
	return false
}

var err error

type Operation struct {
	UnitId       uint8  `json:"unit_id"`     //device addr
	ObjectType   string `json:"object_type"` //read_coil
	op           uint
	Addr         uint16  `json:"addr"`
	ZeroMode     bool    `json:"zero_mode"`
	Length       uint16  `json:"length"`
	IsCoil       bool    `json:"is_coil"`
	IsHoldingReg bool    `json:"is_holding_register"`
	WriteValue   float64 `json:"write_value"`
	Encoding     string  `json:"object_encoding"` //BEB_LEW
	coil         uint16
	u16          uint16
	u32          uint32
	f32          float32
	u64          uint64
	f64          float64
}

func pointWrite(pnt *model.Point) (out float64) {
	if pnt.Priority != nil {
		if (*pnt.Priority).P16 != nil {
			out = *pnt.Priority.P16
			log.Infof("modbus-write: WRITE ObjectType: %s  Addr: %d WriteValue: %v\n", pnt.ObjectType, pnt.AddressID, out)
		}
	}
	return
}

func writeCoilPayload(in float64) (out uint16) {
	if in > 0 {
		out = 0xFF00
	} else {
		out = 0x0000
	}
	return
}

func pointAddress(pnt *model.Point, zeroMode bool) (out uint16, err error) {
	address := utils.IntIsNil(pnt.AddressID)
	//zeroMode will subtract 1 from the register address, so address 1 will be address 0 if set to true
	if !zeroMode {
		if address <= 0 {
			return 0, nil
		} else {
			return uint16(address) - 1, nil
		}
	} else {
		if address <= 0 {
			return 0, nil
		}
		return uint16(address), nil
	}
}

func networkRequest(mbClient smod.ModbusClient, pnt *model.Point) (response interface{}, responseValue float64, err error) {

	objectEncoding := pnt.ObjectEncoding                        //beb_lew
	objectType := utils.NewString(pnt.ObjectType).ToSnakeCase() //eg: readCoil, read_coil, writeCoil
	address, err := pointAddress(pnt, mbClient.DeviceZeroMode)  //register address
	length := utils.IntIsNil(pnt.AddressLength)                 //modbus register length
	writeValue := pointWrite(pnt)
	isOutput := utils.BoolIsNil(pnt.IsOutput)

	switch objectEncoding {
	case string(model.ByteOrderLebBew):
		err = mbClient.SetEncoding(smod.LittleEndian, smod.HighWordFirst)
	case string(model.ByteOrderLebLew):
		err = mbClient.SetEncoding(smod.LittleEndian, smod.LowWordFirst)
	case string(model.ByteOrderBebLew):
		err = mbClient.SetEncoding(smod.BigEndian, smod.LowWordFirst)
	case string(model.ByteOrderBebBew):
		err = mbClient.SetEncoding(smod.BigEndian, smod.HighWordFirst)
	default:
		err = mbClient.SetEncoding(smod.BigEndian, smod.LowWordFirst)
	}
	if length <= 0 { //make sure length is > 0
		length = 1
	}
	if utils.BoolIsNil(pnt.IsOutput) {
		log.Infof("modbus-write: ObjectType: %s  Addr: %d WriteValue: %v\n", objectType, address, writeValue)
	} else {
		log.Infof("modbus-read: ObjectType: %s  Addr: %d", objectType, address)
	}

	switch objectType {
	//COILS
	case string(model.ObjTypeReadCoil):
		return mbClient.ReadCoil(address)
	case string(model.ObjTypeWriteCoil):
		return mbClient.WriteCoil(address, writeCoilPayload(writeValue))
		//FLOAT32
	case string(model.ObjTypeReadFloat32):
		if isOutput {
			return mbClient.ReadFloat32(address, smod.HoldingRegister)
		} else {
			return mbClient.ReadFloat32(address, smod.InputRegister)
		}
	case string(model.ObjTypeWriteFloat32):
		return mbClient.WriteFloat32(address, writeValue)
	}

	return nil, 0, nil
}
