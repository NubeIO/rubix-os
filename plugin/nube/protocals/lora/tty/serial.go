package tty

import (
	"bufio"
	"fmt"
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

func (p *SerialSetting) NewSerialConnection() {
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
		}
	}
	log.Info("SERIAL: try and connect to:", portName)
	m := &serial.Mode{
		BaudRate: baudRate,
		Parity:   parity,
		DataBits: dataBits,
		StopBits: stopBits,
	}
	ports, err := serial.GetPortsList()
	log.Info("SERIAL: ports: ", ports)
	p.ActivePortList = ports
	port, err := serial.Open(portName, m)
	Port = port
	if err != nil {
		p.Error = true
		log.Fatal("SERIAL: ", err)
	}
	p.Connected = true
	log.Info("SERIAL: Connected to serial port: ", portName, " ", "connected: ", p.Connected)

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
			log.Printf("loop count %d", count)
			commonData, fullData := decoder.DecodePayload(data)
			payload.PublishSensor(commonData, fullData)
		} else {
			log.Printf("lora serial messsage size %d", len(data))
		}
	}
}
func Disconnect() error {
	fmt.Println(Port.Close(), 99999456)
	if Port != nil {
		err := Port.Close();if err != nil {
			log.Error("SERIAL: err on trying to close the port")
			return err
		}
	}
	return nil
}
