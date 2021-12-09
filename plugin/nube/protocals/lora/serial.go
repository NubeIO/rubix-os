package main

import (
	"bufio"
	"errors"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lora/decoder"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type SerialSetting struct {
	SerialPort     string
	Enable         bool
	BaudRate       int
	StopBits       serial.StopBits
	Parity         serial.Parity
	DataBits       int
	Timeout        int
	ActivePortList []string
	Connected      bool
	Error          bool
	I              Instance
}

func (i *Instance) SerialOpen() error {
	s := new(SerialSetting)
	var arg api.Args
	net, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if err != nil {
		return err
	}
	if net.SerialPort == nil || net.SerialBaudRate == nil {
		return errors.New("serial_port & serial_baud_rate required to open")
	}
	s.SerialPort = *net.SerialPort
	s.BaudRate = int(*net.SerialBaudRate)
	go func() error {
		sc := New(s)
		_, err = sc.NewSerialConnection()
		if err != nil {
			log.Errorf("lora: issue on SerialOpenAndRead: %v\n", err)
		}
		sc.Loop()
		return nil
	}()
	return nil

}

func (i *Instance) SerialClose() error {
	err := Disconnect()
	if err != nil {
		return err
	}
	return nil
}

func New(s *SerialSetting) *SerialSetting {
	if s.SerialPort == "" {
		s.SerialPort = "/dev/ttyACM0"
	}
	if s.BaudRate == 0 {
		s.BaudRate = 38400
	}
	return &SerialSetting{
		SerialPort: s.SerialPort,
		BaudRate:   s.BaudRate,
	}
}

var Port serial.Port

func (s *SerialSetting) NewSerialConnection() (connected bool, err error) {
	portName := s.SerialPort
	baudRate := s.BaudRate
	parity := s.Parity
	stopBits := s.StopBits
	dataBits := s.DataBits
	if s.Connected {
		log.Info("Existing serial port connection by this app is open So! close existing connection")
		err := Disconnect()
		if err != nil {
			log.Info(err)
			s.Error = true
			return false, err
		}
	}
	log.Info("LORA: connecting to port:", portName)
	m := &serial.Mode{
		BaudRate: baudRate,
		Parity:   parity,
		DataBits: dataBits,
		StopBits: stopBits,
	}

	ports, err := serial.GetPortsList()
	s.ActivePortList = ports

	port, err := serial.Open(portName, m)
	if err != nil {
		s.Error = true
		log.Error("LORA: error on open port", " ", err)
		return false, err
	}
	Port = port
	s.Connected = true
	log.Info("LORA: Connected to serial port: ", " ", portName, " ", "connected: ", " ", s.Connected)
	return s.Connected, nil
}

func (s *SerialSetting) Loop() {
	if s.Error || !s.Connected || Port == nil {
		return
	}
	count := 0
	scanner := bufio.NewScanner(Port)
	for scanner.Scan() {
		var data = scanner.Text()
		if decoder.CheckPayloadLength(data) {
			count = count + 1
			commonData, fullData := decoder.DecodePayload(data)
			s.I.updateDevicePointValues(commonData, fullData)
		} else {
			log.Printf("LORA: serial messsage size %d", len(data))
		}
	}
}
func Disconnect() error {
	if Port != nil {
		err := Port.Close()
		if err != nil {
			log.Error("LORA: err on trying to close the port")
			return err
		}
	}
	return nil
}
