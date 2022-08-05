package main

import (
	"fmt"
	"github.com/NubeDev/modbus"
	log "github.com/sirupsen/logrus"
	"net"
)

type ModbusServer struct {
	Listener net.Listener
	Server   *modbus.TCPServer
}

const size = 0x10000

var discretes [size]bool
var coils [size]bool
var inputRegisters [size]uint16
var holdingRegisters [size]uint16

func (inst *Instance) getIP() string {
	p := inst.config.Port
	if inst.config.Port == 0 {
		p = 10505
	}
	ip := inst.config.Ip
	if inst.config.Port == 0 {
		ip = "0.0.0.0"
	}
	return fmt.Sprintf("%s:%d", ip, p)
}

func (inst *Instance) serverInit() (*ModbusServer, error) {
	host := inst.getIP()
	listener, err := net.Listen("tcp", host)
	if err != nil {
		inst.log(true, "modbus-server:", "serverInit():", "start server error listener: ", err)
		return nil, err

	}
	inst.log(false, "modbus-server:", "serverInit():", "start server with IP", host)
	device := modbus.NewTCPServer(listener)
	err = device.Serve(handlerGenerator())
	if err != nil {
		inst.log(true, "modbus-server:", "serverInit():", "start server error handlerGenerator: ", err)
		return nil, err
	}
	server := &ModbusServer{
		Listener: listener,
		Server:   device,
	}
	inst.ModbusServer = server
	return inst.ModbusServer, nil
}

func handlerGenerator() modbus.ProtocolHandler {
	return &modbus.SimpleHandler{
		ReadDiscreteInputs: func(address, quantity uint16) ([]bool, error) {
			log.Infoln("modbus-server: ReadDiscreteInputs from %v, quantity %v\n", address, quantity)
			return discretes[address : address+quantity], nil
		},
		WriteDiscreteInputs: func(address uint16, values []bool) error {
			log.Printf("modbus-server: WriteDiscreteInputs from %v, quantity %v\n", address, len(values))
			for i, v := range values {
				discretes[address+uint16(i)] = v
			}
			return nil
		},

		ReadCoils: func(address, quantity uint16) ([]bool, error) {
			log.Printf("modbus-server:  ReadCoils from %v, quantity %v\n", address, quantity)
			return coils[address : address+quantity], nil
		},
		WriteCoils: func(address uint16, values []bool) error {
			log.Printf("modbus-server:  WriteCoils from %v, quantity %v\n", address, len(values))
			for i, v := range values {
				coils[address+uint16(i)] = v
				log.Println(i, v)
			}
			return nil
		},

		ReadInputRegisters: func(address, quantity uint16) ([]uint16, error) {
			log.Printf("modbus-server:  ReadInputRegisters from %v, quantity %v\n", address, quantity)
			return inputRegisters[address : address+quantity], nil
		},
		WriteInputRegisters: func(address uint16, values []uint16) error {
			log.Printf("modbus-server:  WriteInputRegisters from %v, quantity %v\n", address, len(values))
			for i, v := range values {
				inputRegisters[address+uint16(i)] = v
			}
			return nil
		},

		ReadHoldingRegisters: func(address, quantity uint16) ([]uint16, error) {
			log.Printf("modbus-server:  ReadHoldingRegisters from %v, quantity %v\n", address, quantity)
			return holdingRegisters[address : address+quantity], nil
		},
		WriteHoldingRegisters: func(address uint16, values []uint16) error {
			log.Printf("modbus-server:  WriteHoldingRegisters from %v, quantity %v\n", address, len(values))
			for i, v := range values {
				holdingRegisters[address+uint16(i)] = v
			}
			return nil
		},

		OnErrorImp: func(req modbus.PDU, errRep modbus.PDU) {
			log.Errorf("modbus-server: error received: %v from req: %v\n", errRep, req)
		},
	}
}
