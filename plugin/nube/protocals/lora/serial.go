package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
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
	log.Info("LORA: try and connect to:", portName)
	m := &serial.Mode{
		BaudRate: baudRate,
		Parity:   parity,
		DataBits: dataBits,
		StopBits: stopBits,
	}
	ports, err := serial.GetPortsList()
	log.Info("LORA: ports: ", ports)
	portNameFound := ""
	for _, port := range ports {
		if port == portName {
			portNameFound = portName
		}
	}
	if portNameFound == "" {
		log.Errorf("LORA: port not found: %v\n", s.SerialPort)
		return false, errors.New("LORA: port not found")
	}
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
	fmt.Println(scanner)
	for scanner.Scan() {
		var data = scanner.Text()
		if decoder.CheckPayloadLength(data) {
			count = count + 1
			commonData, fullData := decoder.DecodePayload(data)
			s.I.publishSensor(commonData, fullData)
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
