package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/simonvetter/modbus"
	"os"
	"strings"
	"time"
)

const (
	readBools uint = iota + 1
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
	setUnitId
	sleep
	repeat
	date
	scanBools
	scanRegisters
	ping
)

var err error

type operation struct {
	op           uint
	addr         uint16
	isCoil       bool
	isHoldingReg bool
	quantity     uint16
	coil         bool
	u16          uint16
	u32          uint32
	f32          float32
	u64          uint64
	f64          float64
	duration     time.Duration
	unitId       uint8
}

func operations(client *modbus.ModbusClient, o operation) (interface{}, error ) {
	switch o.op {
	case readBools:
		var res []bool
		if o.isCoil {
			res, err = client.ReadCoils(o.addr, o.quantity+1)
			return res, err
		} else {
			res, err = client.ReadDiscreteInputs(o.addr, o.quantity+1)
		}
		if err != nil {
			fmt.Printf("failed to read coils/discrete inputs: %v\n", err)
		} else {
			return res, err
		}
	case readUint16, readInt16:
		var res []uint16

		if o.isHoldingReg {
			res, err = client.ReadRegisters(o.addr, o.quantity+1, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadRegisters(o.addr, o.quantity+1, modbus.INPUT_REGISTER)
		}
		if err != nil {
			fmt.Printf("failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				if o.op == readUint16 {
					fmt.Printf("0x%04x\t%-5v : 0x%04x\t%v\n",
						o.addr+uint16(idx),
						o.addr+uint16(idx),
						res[idx], res[idx])
				} else {
					fmt.Printf("0x%04x\t%-5v : 0x%04x\t%v\n",
						o.addr+uint16(idx),
						o.addr+uint16(idx),
						res[idx], int16(res[idx]))
				}
			}
		}
	case readUint32, readInt32:
		var res []uint32

		if o.isHoldingReg {
			res, err = client.ReadUint32s(o.addr, o.quantity+1, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadUint32s(o.addr, o.quantity+1, modbus.INPUT_REGISTER)
		}
		if err != nil {
			fmt.Printf("failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				if o.op == readUint32 {
					fmt.Printf("0x%04x\t%-5v : 0x%08x\t%v\n",
						o.addr+(uint16(idx)*2),
						o.addr+(uint16(idx)*2),
						res[idx], res[idx])
				} else {
					fmt.Printf("0x%04x\t%-5v : 0x%08x\t%v\n",
						o.addr+(uint16(idx)*2),
						o.addr+(uint16(idx)*2),
						res[idx], int32(res[idx]))
				}
			}
		}

	case readFloat32:
		var res []float32

		if o.isHoldingReg {
			res, err = client.ReadFloat32s(o.addr, o.quantity+1, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadFloat32s(o.addr, o.quantity+1, modbus.INPUT_REGISTER)
		}
		if err != nil {
			fmt.Printf("failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				fmt.Printf("0x%04x\t%-5v : %f\n",
					o.addr+(uint16(idx)*2),
					o.addr+(uint16(idx)*2),
					res[idx])
			}
		}

	case readUint64, readInt64:
		var res []uint64

		if o.isHoldingReg {
			res, err = client.ReadUint64s(o.addr, o.quantity+1, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadUint64s(o.addr, o.quantity+1, modbus.INPUT_REGISTER)
		}
		if err != nil {
			fmt.Printf("failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				if o.op == readUint64 {
					fmt.Printf("0x%04x\t%-5v : 0x%016x\t%v\n",
						o.addr+(uint16(idx)*4),
						o.addr+(uint16(idx)*4),
						res[idx], res[idx])
				} else {
					fmt.Printf("0x%04x\t%-5v : 0x%016x\t%v\n",
						o.addr+(uint16(idx)*4),
						o.addr+(uint16(idx)*4),
						res[idx], int64(res[idx]))
				}
			}
		}

	case readFloat64:
		var res []float64
		if o.isHoldingReg {
			res, err = client.ReadFloat64s(o.addr, o.quantity+1, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadFloat64s(o.addr, o.quantity+1, modbus.INPUT_REGISTER)
		}
		if err != nil {
			fmt.Printf("failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				fmt.Printf("0x%04x\t%-5v : %f\n",
					o.addr+(uint16(idx)*4),
					o.addr+(uint16(idx)*4),
					res[idx])
			}
		}

	case writeCoil:
		err = client.WriteCoil(o.addr, o.coil)
		if err != nil {
			fmt.Printf("failed to write %v at coil address 0x%04x: %v\n",
				o.coil, o.addr, err)
		} else {
			fmt.Printf("wrote %v at coil address 0x%04x\n",
				o.coil, o.addr)
		}

	case writeUint16:
		err = client.WriteRegister(o.addr, o.u16)
		if err != nil {
			fmt.Printf("failed to write %v at register address 0x%04x: %v\n",
				o.u16, o.addr, err)
		} else {
			fmt.Printf("wrote %v at register address 0x%04x\n",
				o.u16, o.addr)
		}

	case writeInt16:
		err = client.WriteRegister(o.addr, o.u16)
		if err != nil {
			fmt.Printf("failed to write %v at register address 0x%04x: %v\n",
				int16(o.u16), o.addr, err)
		} else {
			fmt.Printf("wrote %v at register address 0x%04x\n",
				int16(o.u16), o.addr)
		}

	case writeUint32:
		err = client.WriteUint32(o.addr, o.u32)
		if err != nil {
			fmt.Printf("failed to write %v at address 0x%04x: %v\n",
				o.u32, o.addr, err)
		} else {
			fmt.Printf("wrote %v at address 0x%04x\n",
				o.u32, o.addr)
		}

	case writeInt32:
		err = client.WriteUint32(o.addr, o.u32)
		if err != nil {
			fmt.Printf("failed to write %v at address 0x%04x: %v\n",
				int32(o.u32), o.addr, err)
		} else {
			fmt.Printf("wrote %v at address 0x%04x\n",
				int32(o.u32), o.addr)
		}

	case writeFloat32:
		err = client.WriteFloat32(o.addr, o.f32)
		if err != nil {
			fmt.Printf("failed to write %f at address 0x%04x: %v\n",
				o.f32, o.addr, err)
		} else {
			fmt.Printf("wrote %f at address 0x%04x\n",
				o.f32, o.addr)
		}

	case writeUint64:
		err = client.WriteUint64(o.addr, o.u64)
		if err != nil {
			fmt.Printf("failed to write %v at address 0x%04x: %v\n",
				o.u64, o.addr, err)
		} else {
			fmt.Printf("wrote %v at address 0x%04x\n",
				o.u64, o.addr)
		}

	case writeInt64:
		err = client.WriteUint64(o.addr, o.u64)
		if err != nil {
			fmt.Printf("failed to write %v at address 0x%04x: %v\n",
				int64(o.u64), o.addr, err)
		} else {
			fmt.Printf("wrote %v at address 0x%04x\n",
				int64(o.u64), o.addr)
		}

	case writeFloat64:
		err = client.WriteFloat64(o.addr, o.f64)
		if err != nil {
			fmt.Printf("failed to write %f at address 0x%04x: %v\n",
				o.f64, o.addr, err)
		} else {
			fmt.Printf("wrote %f at address 0x%04x\n",
				o.f64, o.addr)
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

func ipCheck(target string) bool {
	if strings.HasPrefix(target, "tcp//") {
		return true
	} else {
		return false
	}
}

func fEndianness(endianness string) modbus.Endianness {
	var e modbus.Endianness
	switch endianness {
	case "big":
		e = modbus.BIG_ENDIAN
	case "little":
		e = modbus.LITTLE_ENDIAN
	default:
		fmt.Printf("unknown endianness setting '%s' (should either be big or little)\n",
			endianness)
	}
	return e
}

func fWordOrder(wordOrder string) modbus.WordOrder {
	var w modbus.WordOrder
	switch wordOrder {
	case "highfirst", "hf":
		w = modbus.HIGH_WORD_FIRST
	case "lowfirst", "lf":
		w = modbus.LOW_WORD_FIRST
	default:
		fmt.Printf("unknown word order setting '%s' (should be one of highfirst, hf, littlefirst, lf)\n",
			w)
		os.Exit(1)
	}
	return w
}
