package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/flow-framework/utils/writemode"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"math"
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

type Operation struct {
	UnitId       uint8  `json:"unit_id"`     // device addr
	ObjectType   string `json:"object_type"` // read_coil
	op           uint
	Addr         uint16  `json:"addr"`
	ZeroMode     bool    `json:"zero_mode"`
	Length       uint16  `json:"length"`
	IsCoil       bool    `json:"is_coil"`
	IsHoldingReg bool    `json:"is_holding_register"`
	WriteValue   float64 `json:"write_value"`
	Encoding     string  `json:"object_encoding"` // BEB_LEW
	coil         uint16
	u16          uint16
	u32          uint32
	f32          float32
	u64          uint64
	f64          float64
}

func pointWrite(pnt *model.Point) (out float64) {
	out = float.NonNil(pnt.WriteValue)
	log.Infof("modbus-write: pointWrite() ObjectType: %s  Addr: %d WriteValue: %v\n", pnt.ObjectType, integer.NonNil(pnt.AddressID), out)
	// if pnt.Priority != nil {
	//	if (*pnt.Priority).P16 != nil {
	//		out = *pnt.Priority.P16
	//		log.Infof("modbus-write: pointWrite() ObjectType: %s  Addr: %d WriteValue: %v\n", pnt.ObjectType, utils.NonNil(pnt.AddressID), out)
	//	}
	// }
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

func pointAddress(pnt *model.Point, zeroMode bool) uint16 {
	address := integer.NonNil(pnt.AddressID)
	// zeroMode will subtract 1 from the register address, so address 1 will be address 0 if set to true
	if !zeroMode {
		return uint16(address) - 1
	} else {
		return uint16(address)
	}
}

func (inst *Instance) networkRequest(mbClient smod.ModbusClient, pnt *model.Point, doWrite bool) (response interface{}, responseValue float64, err error) {
	mbClient.Debug = true
	objectEncoding := pnt.ObjectEncoding                          // beb_lew
	objectType := nstring.NewString(pnt.ObjectType).ToSnakeCase() // eg: readCoil, read_coil, writeCoil
	dataType := nstring.NewString(pnt.DataType).ToSnakeCase()     // eg: int16, uint16
	address := pointAddress(pnt, mbClient.DeviceZeroMode)         // register address
	length := integer.NonNil(pnt.AddressLength)                   // modbus register length

	switch objectEncoding {
	case string(model.ByteOrderLebBew):
		mbClient.SetEncoding(smod.LittleEndian, smod.HighWordFirst)
	case string(model.ByteOrderLebLew):
		mbClient.SetEncoding(smod.LittleEndian, smod.LowWordFirst)
	case string(model.ByteOrderBebLew):
		mbClient.SetEncoding(smod.BigEndian, smod.LowWordFirst)
	case string(model.ByteOrderBebBew):
		mbClient.SetEncoding(smod.BigEndian, smod.HighWordFirst)
	default:
		mbClient.SetEncoding(smod.BigEndian, smod.LowWordFirst)
	}
	if length <= 0 { // make sure length is > 0
		length = 1
	}
	var writeValue float64
	if doWrite {
		writeValue = pointWrite(pnt)
	}

	if doWrite {
		inst.modbusDebugMsg("modbus-write: ObjectType: %s  Addr: %d WriteValue: %v\n", objectType, address, writeValue)
	} else {
		inst.modbusDebugMsg("modbus-read: ObjectType: %s  Addr: %d", objectType, address)
	}

	switch objectType {
	// COILS
	case string(model.ObjTypeReadCoil):
		return mbClient.ReadCoils(address, uint16(length))
	case string(model.ObjTypeWriteCoil):
		return mbClient.WriteCoil(address, writeCoilPayload(writeValue))
		// READ DISCRETE INPUTS
	case string(model.ObjTypeReadDiscreteInput):
		return mbClient.ReadDiscreteInputs(address, uint16(length))
		// READ HOLDINGS
	case string(model.ObjTypeReadHolding):
		if dataType == string(model.TypeUint16) || dataType == string(model.TypeInt16) {
			return mbClient.ReadHoldingRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeUint32) || dataType == string(model.TypeInt32) {
			return mbClient.ReadHoldingRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeUint64) || dataType == string(model.TypeInt64) {
			return mbClient.ReadHoldingRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeFloat32) {
			return mbClient.ReadFloat32(address, smod.HoldingRegister)
		} else if dataType == string(model.TypeFloat64) {
			return mbClient.ReadFloat32(address, smod.HoldingRegister)
		}
		// READ INPUT REGISTERS
	case string(model.ObjTypeReadRegister):
		if dataType == string(model.TypeUint16) || dataType == string(model.TypeInt16) {
			return mbClient.ReadInputRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeUint32) || dataType == string(model.TypeInt32) {
			return mbClient.ReadInputRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeUint64) || dataType == string(model.TypeInt64) {
			return mbClient.ReadInputRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeFloat32) {
			return mbClient.ReadFloat32(address, smod.InputRegister)
		} else if dataType == string(model.TypeFloat64) {
			return mbClient.ReadFloat32(address, smod.InputRegister)
		}
		// WRITE HOLDINGS
	case string(model.ObjTypeWriteHolding):
		if dataType == string(model.TypeUint16) || dataType == string(model.TypeInt16) {
			return mbClient.WriteSingleRegister(address, uint16(writeValue))
		} else if dataType == string(model.TypeUint32) || dataType == string(model.TypeInt32) {
			return mbClient.WriteDoubleRegister(address, uint32(writeValue))
		} else if dataType == string(model.TypeUint64) || dataType == string(model.TypeInt64) {
			return mbClient.WriteQuadRegister(address, uint64(writeValue))
		} else if dataType == string(model.TypeFloat32) {
			return mbClient.WriteFloat32(address, writeValue)
		} else if dataType == string(model.TypeFloat64) {
			return mbClient.WriteFloat32(address, writeValue)
		}

	}

	return nil, 0, nil
}

func (inst *Instance) networkWrite(mbClient smod.ModbusClient, pnt *model.Point) (response interface{}, responseValue float64, err error) {
	if pnt.WriteValue == nil {
		return nil, 0, errors.New("modbus-write: point has no WriteValue")
	}
	mbClient.Debug = true
	objectEncoding := pnt.ObjectEncoding                          // beb_lew
	objectType := nstring.NewString(pnt.ObjectType).ToSnakeCase() // eg: readCoil, read_coil, writeCoil
	dataType := nstring.NewString(pnt.DataType).ToSnakeCase()     // eg: int16, uint16
	address := pointAddress(pnt, mbClient.DeviceZeroMode)         // register address
	length := integer.NonNil(pnt.AddressLength)                   // modbus register length

	switch objectEncoding {
	case string(model.ByteOrderLebBew):
		mbClient.SetEncoding(smod.LittleEndian, smod.HighWordFirst)
	case string(model.ByteOrderLebLew):
		mbClient.SetEncoding(smod.LittleEndian, smod.LowWordFirst)
	case string(model.ByteOrderBebLew):
		mbClient.SetEncoding(smod.BigEndian, smod.LowWordFirst)
	case string(model.ByteOrderBebBew):
		mbClient.SetEncoding(smod.BigEndian, smod.HighWordFirst)
	default:
		mbClient.SetEncoding(smod.BigEndian, smod.LowWordFirst)
	}
	if length <= 0 { // make sure length is > 0
		length = 1
	}

	writeValue := *pnt.WriteValue

	inst.modbusDebugMsg(fmt.Sprintf("modbus-write: ObjectType: %s  Addr: %d WriteValue: %v\n", objectType, address, writeValue))

	switch objectType {
	// WRITE COILS
	case string(model.ObjTypeWriteCoil), string(model.ObjTypeWriteCoils):
		return mbClient.WriteCoil(address, writeCoilPayload(writeValue))

	// WRITE HOLDINGS
	case string(model.ObjTypeWriteHolding), string(model.ObjTypeWriteHoldings):
		if dataType == string(model.TypeUint16) || dataType == string(model.TypeInt16) {
			return mbClient.WriteSingleRegister(address, uint16(writeValue))
		} else if dataType == string(model.TypeUint32) || dataType == string(model.TypeInt32) {
			return mbClient.WriteDoubleRegister(address, uint32(writeValue))
		} else if dataType == string(model.TypeUint64) || dataType == string(model.TypeInt64) {
			return mbClient.WriteQuadRegister(address, uint64(writeValue))
		} else if dataType == string(model.TypeFloat32) {
			return mbClient.WriteFloat32(address, writeValue)
		} else if dataType == string(model.TypeFloat64) {
			return mbClient.WriteFloat32(address, writeValue)
		}
	}

	return nil, 0, errors.New("modbus-write: dataType is not recognized")
}

func (inst *Instance) networkRead(mbClient smod.ModbusClient, pnt *model.Point) (response interface{}, responseValue float64, err error) {
	mbClient.Debug = true
	objectEncoding := pnt.ObjectEncoding                          // beb_lew
	objectType := nstring.NewString(pnt.ObjectType).ToSnakeCase() // eg: readCoil, read_coil, writeCoil
	dataType := nstring.NewString(pnt.DataType).ToSnakeCase()     // eg: int16, uint16
	address := pointAddress(pnt, mbClient.DeviceZeroMode)         // register address
	length := integer.NonNil(pnt.AddressLength)                   // modbus register length

	switch objectEncoding {
	case string(model.ByteOrderLebBew):
		mbClient.SetEncoding(smod.LittleEndian, smod.HighWordFirst)
	case string(model.ByteOrderLebLew):
		mbClient.SetEncoding(smod.LittleEndian, smod.LowWordFirst)
	case string(model.ByteOrderBebLew):
		mbClient.SetEncoding(smod.BigEndian, smod.LowWordFirst)
	case string(model.ByteOrderBebBew):
		mbClient.SetEncoding(smod.BigEndian, smod.HighWordFirst)
	default:
		mbClient.SetEncoding(smod.BigEndian, smod.LowWordFirst)
	}
	if length <= 0 { // make sure length is > 0
		length = 1
	}

	inst.modbusDebugMsg(fmt.Sprintf("modbus-read: ObjectType: %s  Addr: %d", objectType, address))

	switch objectType {
	// COILS
	case string(model.ObjTypeReadCoil), string(model.ObjTypeReadCoils), string(model.ObjTypeWriteCoil), string(model.ObjTypeWriteCoils):
		return mbClient.ReadCoils(address, uint16(length))

	// READ DISCRETE INPUTS
	case string(model.ObjTypeReadDiscreteInput), string(model.ObjTypeReadDiscreteInputs):
		return mbClient.ReadDiscreteInputs(address, uint16(length))

	// READ INPUT REGISTERS
	case string(model.ObjTypeReadRegister), string(model.ObjTypeReadRegisters):
		if dataType == string(model.TypeUint16) || dataType == string(model.TypeInt16) {
			return mbClient.ReadInputRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeUint32) || dataType == string(model.TypeInt32) {
			return mbClient.ReadInputRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeUint64) || dataType == string(model.TypeInt64) {
			return mbClient.ReadInputRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeFloat32) {
			return mbClient.ReadFloat32(address, smod.InputRegister)
		} else if dataType == string(model.TypeFloat64) {
			return mbClient.ReadFloat32(address, smod.InputRegister)
		}

	// READ HOLDINGS
	case string(model.ObjTypeReadHolding), string(model.ObjTypeReadHoldings), string(model.ObjTypeWriteHolding), string(model.ObjTypeWriteHoldings):
		if dataType == string(model.TypeUint16) || dataType == string(model.TypeInt16) {
			return mbClient.ReadHoldingRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeUint32) || dataType == string(model.TypeInt32) {
			return mbClient.ReadHoldingRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeUint64) || dataType == string(model.TypeInt64) {
			return mbClient.ReadHoldingRegisters(address, uint16(length), dataType)
		} else if dataType == string(model.TypeFloat32) {
			return mbClient.ReadFloat32(address, smod.HoldingRegister)
		} else if dataType == string(model.TypeFloat64) {
			return mbClient.ReadFloat32(address, smod.HoldingRegister)
		}

	}

	return nil, 0, errors.New("modbus-read: dataType is not recognized")
}

func SetPriorityArrayModeBasedOnWriteMode(pnt *model.Point) bool {
	switch pnt.WriteMode {
	case model.ReadOnce, model.ReadOnly:
		pnt.PointPriorityArrayMode = model.ReadOnlyNoPriorityArrayRequired
		return true
	case model.WriteOnce, model.WriteOnceReadOnce, model.WriteAlways, model.WriteOnceThenRead, model.WriteAndMaintain:
		pnt.PointPriorityArrayMode = model.PriorityArrayToWriteValue
		return true
	}
	return false
}

func isWriteable(writeMode model.WriteMode, objectType string) bool {
	if isWriteableObjectType(objectType) && writemode.IsWriteable(writeMode) {
		return true
	} else {
		return false
	}
}

func isWriteableObjectType(objectType string) bool {
	switch objectType {
	case string(model.ObjTypeWriteCoil), string(model.ObjTypeWriteCoils):
		return true
	case string(model.ObjTypeWriteHolding), string(model.ObjTypeWriteHoldings):
		return true
	case string(model.ObjTypeWriteInt16), string(model.ObjTypeWriteUint16):
		return true
	case string(model.ObjTypeWriteFloat32):
		return true
	}
	return false
}

func checkForBooleanType(ObjectType, DataType string) (isTypeBool bool) {
	isTypeBool = false
	if DataType == string(model.TypeDigital) {
		isTypeBool = true
	}
	switch ObjectType {
	case string(model.ObjTypeReadCoil), string(model.ObjTypeReadCoils), string(model.ObjTypeReadDiscreteInput), string(model.ObjTypeReadDiscreteInputs), string(model.ObjTypeWriteCoil), string(model.ObjTypeWriteCoils):
		isTypeBool = true
	}
	return
}

func checkForOutputType(ObjectType string) (isOutput bool) {
	isOutput = false
	switch ObjectType {
	case string(model.ObjTypeWriteCoil), string(model.ObjTypeWriteCoils), string(model.ObjTypeWriteHolding), string(model.ObjTypeWriteHoldings), string(model.ObjTypeWriteInt16), string(model.ObjTypeWriteUint16), string(model.ObjTypeWriteFloat32), string(model.ObjTypeWriteFloat64):
		isOutput = true
	}
	return
}

func getBitsFromFloat64(value float64) (bitArray []bool, err error) {
	if math.Mod(value, 1) != 0 {
		err = errors.New("cannot get bits from floats")
		return
	}
	if value < 0 {
		err = errors.New("cannot get bits from negative numbers")
		return
	}
	buf := make([]byte, binary.MaxVarintLen64)
	length := binary.PutUvarint(buf, uint64(value))
	fmt.Println("getBitsFromFloat64() length: ", length)
	fmt.Println("getBitsFromFloat64() uint64(value): ", uint64(value))
	fmt.Println("getBitsFromFloat64() raw buf: ", buf)
	fmt.Println(fmt.Sprintf("getBitsFromFloat64() binary: %b", uint64(value)))
	for j := 0; j < length; j++ {
		bits := buf[j]
		for i := 0; bits > 0; i, bits = i+1, bits>>1 {
			if bits&1 == 1 {
				bitArray = append(bitArray, true)
			} else if bits&1 == 0 {
				bitArray = append(bitArray, false)
			}
		}
	}
	fmt.Println("getBitsFromFloat64() bitArray: ", bitArray)
	return
}

func getBitFromFloat64(value float64, reqIndex int) (indexValue bool, err error) {
	if math.Mod(value, 1) != 0 {
		err = errors.New("cannot get bits from floats")
		return
	}
	if value < 0 {
		err = errors.New("cannot get bits from negative numbers")
		return
	}
	buf := make([]byte, binary.MaxVarintLen64)
	length := binary.PutUvarint(buf, uint64(value))
	currentIndex := 0
	for j := 0; j < length; j++ {
		bits := buf[j]
		for i := 0; bits > 0; i, bits = i+1, bits>>1 {
			if bits&1 == 1 && currentIndex == reqIndex {
				return true, nil
			} else if bits&1 == 0 && currentIndex == reqIndex {
				return false, nil
			}
			currentIndex++
		}
	}
	return
}
