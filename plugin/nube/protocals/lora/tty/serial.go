package tty

import (
	"bufio"
	"errors"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/payload"
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
		BaudRate:  s.BaudRate,
	}
}

var Port serial.Port

func (p *SerialSetting) NewSerialConnection() error {
	portName := p.SerialPort
	baudRate := p.BaudRate
	parity := p.Parity
	stopBits := p.StopBits
	dataBits := p.DataBits
	if p.Connected {
		log.Info("Existing serial port connection by this app is open So! close existing connection")
		err := Disconnect()
		if err != nil {
			log.Info(err)
			p.Error = true
			return err
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
		return errors.New("LORA: port not found")
	}
	p.ActivePortList = ports
	port, err := serial.Open(portName, m);if err != nil {
		p.Error = true
		log.Fatal("LORA: error on open port", err)
		return err
	}
	Port = port
	p.Connected = true
	log.Info("LORA: Connected to serial port: ", portName, " ", "connected: ", p.Connected)
	return nil
}

func (p *SerialSetting) Loop() {
	if p.Error || !p.Connected || Port == nil {
		return
	}
	count := 0
	scanner := bufio.NewScanner(Port)
	for scanner.Scan() {
		var data = scanner.Text()
		if decoder.CheckPayloadLength(data) {
			count = count + 1
			commonData, fullData := decoder.DecodePayload(data)
			payload.PublishSensor(commonData, fullData)
		} else {
			log.Printf("LORA: serial messsage size %d", len(data))
		}
	}
}
func Disconnect() error {
	if Port != nil {
		err := Port.Close();if err != nil {
			log.Error("LORA: err on trying to close the port")
			return err
		}
	}
	return nil
}
