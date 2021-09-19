package main

import (
	"errors"
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
	UnitId       uint8   `json:"unit_id"`
	Request      string  `json:"request"` //readCoil
	Op           uint    `json:"op"`
	Addr         uint16  `json:"addr"`
	Length       uint16  `json:"length"`
	IsCoil       bool    `json:"is_coil"`
	IsHoldingReg bool    `json:"is_holding_reg"`
	WriteValue   float64 `json:"write_value"`
	coil         bool
	u16          uint16
	u32          uint32
	f32          float32
	u64          uint64
	f64          float64
}

var req = struct {
	readCoil           string
	readCoils          string
	readDiscreteInput  string
	readDiscreteInputs string
	writeCoil          string
	writeCoils         string
	ReadRegister       string
	ReadRegisters      string
}{
	readCoil:           "readCoil",
	readCoils:          "readCoils",
	readDiscreteInput:  "readDiscreteInput",
	readDiscreteInputs: "readDiscreteInputs",
	writeCoil:          "writeCoil",
	writeCoils:         "writeCoils",
	ReadRegister:       "ReadRegister",
	ReadRegisters:      "ReadRegisters",
}

func setRequest(body Operation) (Operation, error) {
	r := body.Request
	if r == req.readCoil || r == req.readDiscreteInput {
		body.Length = 1
	}
	if r == req.readCoil || r == req.readCoils || r == req.writeCoil || r == req.writeCoils {
		body.IsCoil = true
	}
	return body, nil
}

func parseRequest(body Operation) (Operation, error) {
	set, _ := setRequest(body)
	ops := utils.NewString(body.Request).ToCamelCase() //eg: readCoil, read_coil, writeCoil
	ops = utils.LcFirst(ops)
	switch ops {
	case req.readCoil, req.readCoils, req.readDiscreteInput, req.readDiscreteInputs:
		set.Op = readBool
		return set, err
	case req.writeCoil, req.writeCoils:
		if ops == req.writeCoil || ops == req.writeCoils {
			set.IsCoil = true
		}
		if body.WriteValue > 0 {
			set.coil = true
		} else {
			set.coil = false
		}
		set.Op = writeCoil
		return set, err
	}
	return set, errors.New("req not found")
}

func Operations(client *modbus.ModbusClient, o Operation) (response interface{}, err error) {
	switch o.Op {
	case readBool:
		var res []bool
		if o.IsCoil {
			res, err = client.ReadCoils(o.Addr, o.Length)
			return res, err
		} else {
			res, err = client.ReadDiscreteInputs(o.Addr, o.Length)
		}
		if err != nil {
			log.Infof("modbus: failed to read coils/discrete inputs: %v\n", err)
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
				if o.Op == readUint16 {
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
		} else {
			res, err = client.ReadUint32s(o.Addr, o.Length, modbus.INPUT_REGISTER)
		}
		if err != nil {
			log.Errorf("modbus: failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				if o.Op == readUint32 {
					log.Infof("modbus:  0x%04x\t%-5v : 0x%08x\t%v\n",
						o.Addr+(uint16(idx)*2),
						o.Addr+(uint16(idx)*2),
						res[idx], res[idx])
				} else {
					log.Infof("modbus:  0x%04x\t%-5v : 0x%08x\t%v\n",
						o.Addr+(uint16(idx)*2),
						o.Addr+(uint16(idx)*2),
						res[idx], int32(res[idx]))
				}
			}
		}

	case readFloat32:
		var res []float32

		if o.IsHoldingReg {
			res, err = client.ReadFloat32s(o.Addr, o.Length, modbus.HOLDING_REGISTER)
		} else {
			res, err = client.ReadFloat32s(o.Addr, o.Length, modbus.INPUT_REGISTER)
		}
		if err != nil {
			log.Errorf("modbus:  failed to read holding/input registers: %v\n", err)
		} else {
			for idx := range res {
				log.Infof("modbus: 0x%04x\t%-5v : %f\n",
					o.Addr+(uint16(idx)*2),
					o.Addr+(uint16(idx)*2),
					res[idx])
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
				if o.Op == readUint64 {
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
			log.Infof("modbus: wrote %f at address 0x%04x\n",
				o.f32, o.Addr)
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

func ipCheck(target string) bool {
	if strings.HasPrefix(target, "tcp//") {
		return true
	} else {
		return false
	}
}

var namesEncoding = struct {
	lebBew string
	lebLew string
	bebLew string
	bebBew string
}{
	lebBew: "lebBew",
	lebLew: "lebLew",
	bebLew: "bebLew",
	bebBew: "bebBew",
}

func EncodingBuilder(selection string, client *modbus.ModbusClient) error {
	sel := utils.NewString(selection).ToCamelCase() //eg: LEB_BEW, lebBew
	sel = utils.LcFirst(sel)
	switch sel {
	case namesEncoding.lebBew:
		err = client.SetEncoding(modbus.LITTLE_ENDIAN, modbus.HIGH_WORD_FIRST)
	case namesEncoding.lebLew:
		err = client.SetEncoding(modbus.LITTLE_ENDIAN, modbus.LOW_WORD_FIRST)
	case namesEncoding.bebLew:
		err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.LOW_WORD_FIRST)
	case namesEncoding.bebBew:
		err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.HIGH_WORD_FIRST)
	default:
		log.Errorf("modbus:  unknown endianness setting '%s' (should either be big or little)\n", selection)
	}
	return err
}
