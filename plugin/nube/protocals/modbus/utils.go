package main

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/simonvetter/modbus"
	log "github.com/sirupsen/logrus"
	"strings"
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

var err error

type Operation struct {
	UnitId       uint8  `json:"unit_id"`     //device addr
	ObjectType   string `json:"object_type"` //readCoil
	op           uint
	Addr         uint16  `json:"addr"`
	ZeroMode     bool    `json:"zero_mode"`
	Length       uint16  `json:"length"`
	IsCoil       bool    `json:"is_coil"`
	IsHoldingReg bool    `json:"is_holding_register"`
	WriteValue   float64 `json:"write_value"`
	Encoding     string  `json:"encoding"` //BEB_LEW
	coil         bool
	u16          uint16
	u32          uint32
	f32          float32
	u64          uint64
	f64          float64
}

func setRequest(body Operation) (Operation, error) {
	r := body.ObjectType
	if r == model.ObjectTypes.ReadCoil || r == model.ObjectTypes.ReadDiscreteInput {
		body.Length = 1
	}
	if r == model.ObjectTypes.ReadCoil || r == model.ObjectTypes.ReadCoils || r == model.ObjectTypes.WriteCoil || r == model.ObjectTypes.WriteCoils {
		body.IsCoil = true
	}
	return body, nil
}

func parseRequest(body Operation) (Operation, error) {
	set, _ := setRequest(body)
	ops := utils.NewString(body.ObjectType).ToCamelCase() //eg: readCoil, read_coil, writeCoil
	ops = utils.LcFirst(ops)
	switch ops {
	case model.ObjectTypes.ReadCoil, model.ObjectTypes.ReadCoils, model.ObjectTypes.ReadDiscreteInput, model.ObjectTypes.ReadDiscreteInputs:
		set.op = readBool
		return set, err
	case model.ObjectTypes.WriteCoil, model.ObjectTypes.WriteCoils:
		if ops == model.ObjectTypes.WriteCoil || ops == model.ObjectTypes.WriteCoils {
			set.IsCoil = true
		}
		if body.WriteValue > 0 {
			set.coil = true
		} else {
			set.coil = false
		}
		set.op = writeCoil
		return set, err
	case model.ObjectTypes.ReadFloat32, model.ObjectTypes.ReadSingleFloat32:
		if body.IsHoldingReg {
			set.IsHoldingReg = true
		} else {
			set.IsHoldingReg = false
		}
		set.op = readFloat32
		return set, err
	case model.ObjectTypes.WriteSingleFloat32:
		set.IsHoldingReg = true
		set.op = writeFloat32
		set.f32 = float32(body.WriteValue)
		return set, err
	}
	return set, errors.New("req not found")
}

//zeroMode will subtract 1 from the register address, so address 1 will be address 0 if set to true
func zeroMode(addr uint16, mode bool) uint16 {
	if mode {
		if addr <= 0 {
			return 0
		} else {
			return addr - 1
		}
	} else {
		return addr
	}

}

func DoOperations(client *modbus.ModbusClient, o Operation) (response interface{}, err error) {
	o.Addr = zeroMode(o.Addr, o.ZeroMode)
	sel := utils.NewString(o.Encoding).ToCamelCase() //eg: LEB_BEW, lebBew
	sel = utils.LcFirst(sel)
	switch sel {
	case model.ObjectEncoding.LebBew:
		err = client.SetEncoding(modbus.LITTLE_ENDIAN, modbus.HIGH_WORD_FIRST)
	case model.ObjectEncoding.LebLew:
		err = client.SetEncoding(modbus.LITTLE_ENDIAN, modbus.LOW_WORD_FIRST)
	case model.ObjectEncoding.BebLew:
		err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.LOW_WORD_FIRST)
	case model.ObjectEncoding.BebBew:
		err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.HIGH_WORD_FIRST)
	default:
		err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.HIGH_WORD_FIRST)
	}
	switch o.op {
	case readBool:
		var res []bool
		if o.IsCoil {
			res, err = client.ReadCoils(o.Addr, o.Length)
			return res, err
		} else {
			res, err = client.ReadDiscreteInputs(o.Addr, o.Length)
		}
		if err != nil {
			log.Errorf("modbus: failed to read coils/discrete inputs: %v\n", err)
		} else {
			return res, err
		}
	case readUint16, readInt16:
		var res []uint16
		if o.IsHoldingReg {
			res, err = client.ReadRegisters(o.Addr, o.Length, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadRegisters(o.Addr, o.Length, modbus.INPUT_REGISTER)
		}
		if err != nil {
			log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				if o.op == readUint16 {
					log.Infof("0x%04x\t%-5v : 0x%04x\t%v\n",
						o.Addr+uint16(idx),
						o.Addr+uint16(idx),
						res[idx], res[idx])
				} else {
					log.Infof("0x%04x\t%-5v : 0x%04x\t%v\n",
						o.Addr+uint16(idx),
						o.Addr+uint16(idx),
						res[idx], int16(res[idx]))
				}
			}
		}
	case readUint32, readInt32:
		var res []uint32
		if o.IsHoldingReg {
			res, err = client.ReadUint32s(o.Addr, o.Length, modbus.HOLDING_REGISTER)
			return res, err
		} else {
			res, err = client.ReadUint32s(o.Addr, o.Length, modbus.INPUT_REGISTER)
		}
		if err != nil {
			log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
		} else {
			return res, err
		}
	case readFloat32:
		var res []float32
		if o.ObjectType == model.ObjectTypes.ReadSingleFloat32 {
			o.Length = 1
		}
		if o.IsHoldingReg {
			res, err = client.ReadFloat32s(o.Addr, o.Length, modbus.HOLDING_REGISTER)
			if o.ObjectType == model.ObjectTypes.ReadSingleFloat32 {
				if err != nil {
					log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
					return nil, err
				}
				return res[0], err
			} else {
				if err != nil {
					log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
					return nil, err
				}
				return res, err
			}
		} else {
			res, err = client.ReadFloat32s(o.Addr, o.Length, modbus.INPUT_REGISTER)
			if o.ObjectType == model.ObjectTypes.ReadSingleFloat32 {
				if err != nil {
					log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
				}
				return res[0], err
			} else {
				if err != nil {
					log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
				}
				return res, err
			}
		}
	case readUint64, readInt64:
		var res []uint64
		if o.IsHoldingReg {
			res, err = client.ReadUint64s(o.Addr, o.Length, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadUint64s(o.Addr, o.Length, modbus.INPUT_REGISTER)
		}
		if err != nil {
			log.Errorf("modbus:  failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				if o.op == readUint64 {
					log.Infof("modbus: 0x%04x\t%-5v : 0x%016x\t%v\n",
						o.Addr+(uint16(idx)*4),
						o.Addr+(uint16(idx)*4),
						res[idx], res[idx])
				} else {
					log.Infof("modbus: 0x%04x\t%-5v : 0x%016x\t%v\n",
						o.Addr+(uint16(idx)*4),
						o.Addr+(uint16(idx)*4),
						res[idx], int64(res[idx]))
				}
			}
		}
	case readFloat64:
		var res []float64
		if o.IsHoldingReg {
			res, err = client.ReadFloat64s(o.Addr, o.Length, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadFloat64s(o.Addr, o.Length, modbus.INPUT_REGISTER)
		}
		if err != nil {
			log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				log.Infof("modbus: 0x%04x\t%-5v : %f\n",
					o.Addr+(uint16(idx)*4),
					o.Addr+(uint16(idx)*4),
					res[idx])
			}
		}
	case writeCoil:
		err = client.WriteCoil(o.Addr, o.coil)
		if err != nil {
			log.Infof("modbus: failed to write %v at coil address 0x%04x: %v\n",
				o.coil, o.Addr, err)
			return nil, err
		} else {
			log.Infof("modbus: wrote %v at coil address 0x%04x\n",
				o.coil, o.Addr)
			return o.coil, err
		}
	case writeUint16:
		err = client.WriteRegister(o.Addr, o.u16)
		if err != nil {
			log.Errorf("modbus: failed to write %v at register address 0x%04x: %v\n",
				o.u16, o.Addr, err)
		} else {
			log.Infof("modbus: wrote %v at register address 0x%04x\n",
				o.u16, o.Addr)
		}
	case writeInt16:
		err = client.WriteRegister(o.Addr, o.u16)
		if err != nil {
			log.Infof("modbus: failed to write %v at register address 0x%04x: %v\n",
				int16(o.u16), o.Addr, err)
		} else {
			log.Infof("modbus: wrote %v at register address 0x%04x\n",
				int16(o.u16), o.Addr)
		}
	case writeUint32:
		err = client.WriteUint32(o.Addr, o.u32)
		if err != nil {
			log.Errorf("modbus: failed to write %v at address 0x%04x: %v\n",
				o.u32, o.Addr, err)
		} else {
			log.Infof("modbus: wrote %v at address 0x%04x\n",
				o.u32, o.Addr)
		}
	case writeInt32:
		err = client.WriteUint32(o.Addr, o.u32)
		if err != nil {
			log.Infof("modbus: failed to write %v at address 0x%04x: %v\n",
				int32(o.u32), o.Addr, err)
		} else {
			log.Infof("modbus: wrote %v at address 0x%04x\n",
				int32(o.u32), o.Addr)
		}
	case writeFloat32:
		err = client.WriteFloat32(o.Addr, o.f32)
		if err != nil {
			log.Errorf("modbus: failed to write %f at address 0x%04x: %v\n",
				o.f32, o.Addr, err)
		} else {
			log.Infof("modbus: wrote %f at address 0x%04x\n", o.f32, o.Addr)
			return o.f32, err
		}
	case writeUint64:
		err = client.WriteUint64(o.Addr, o.u64)
		if err != nil {
			log.Errorf("modbus: failed to write %v at address 0x%04x: %v\n",
				o.u64, o.Addr, err)
		} else {
			log.Infof("modbus: wrote %v at address 0x%04x\n",
				o.u64, o.Addr)
		}
	case writeInt64:
		err = client.WriteUint64(o.Addr, o.u64)
		if err != nil {
			log.Errorf("modbus: failed to write %v at address 0x%04x: %v\n",
				int64(o.u64), o.Addr, err)
		} else {
			log.Infof("modbus: wrote %v at address 0x%04x\n",
				int64(o.u64), o.Addr)
		}
	case writeFloat64:
		err = client.WriteFloat64(o.Addr, o.f64)
		if err != nil {
			log.Errorf("modbus: failed to write %f at address 0x%04x: %v\n",
				o.f64, o.Addr, err)
		} else {
			log.Infof("modbus: wrote %f at address 0x%04x\n",
				o.f64, o.Addr)
		}
	}
	return nil, nil
}

func getPointAddr(s string) (objType, addr string) {
	mArr := utils.NewArray()
	ss := strings.Split(s, "-")
	for _, e := range ss {
		if e != "" {
			mArr.Add(e)
		}
	}
	o := mArr.Get(0)
	a := mArr.Get(1)
	return o.(string), a.(string)
}

func performBoolScan(client *modbus.ModbusClient, isCoil bool, start uint32, count uint32) (uint, string) {
	var err error
	var addr uint32
	var val bool
	var countFound uint
	var regType string
	if isCoil {
		regType = "coil"
	} else {
		regType = "discrete input"
	}
	fmt.Printf("starting %s scan\n", regType)
	fmt.Println(start, count)
	for addr = start; addr <= count; addr++ {
		if isCoil {
			val, err = client.ReadCoil(uint16(addr))
		} else {
			val, err = client.ReadDiscreteInput(uint16(addr))
		}
		if err == modbus.ErrIllegalDataAddress || err == modbus.ErrIllegalFunction {
			// the register does not exist
			continue
		} else if err != nil {
			fmt.Printf("failed to read %s at address 0x%04x: %v\n",
				regType, addr, err)
		} else {
			// we found a coil: display its address and value
			fmt.Printf("0x%04x\t%-5v : %v\n", addr, addr, val)
			countFound++
			fmt.Println(countFound)
		}
	}
	fmt.Printf("found %v %ss\n", countFound, regType)
	return countFound, regType
}
