package main

import (
	"fmt"
	"github.com/simonvetter/modbus"
	"go.bug.st/serial"
	"time"
)

func main() {
	var client *modbus.ModbusClient
	var err error

	ports, err := serial.GetPortsList()

	for _, port := range ports {
		fmt.Println(port)
	}

	// for a TCP endpoint
	// (see examples/tls_client.go for TLS usage and options)
	//client, err = modbus.NewClient(&modbus.ClientConfiguration{
	//	URL:     "tcp://192.168.15.202:502",
	//	Timeout: 1 * time.Second,
	//})

	// for an RTU (serial) device/bus
	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:      "rtu:///dev/ttyUSB0",
		Speed:    9600,               // default
		DataBits: 8,                  // default, optional
		Parity:   modbus.PARITY_NONE, // default, optional
		StopBits: 2,                  // default if no parity, optional
		Timeout:  300 * time.Millisecond,
	})

	if err != nil {
		fmt.Println(err)
		// error out if client creation failed
	}
	client.SetUnitId(1)

	// now that the client is created and configured, attempt to connect
	err = client.Open()
	if err != nil {
		// error out if we failed to connect/open the device
		// note: multiple Open() attempts can be made on the same client until
		// the connection succeeds (i.e. err == nil), calling the constructor again
		// is unnecessary.
		// likewise, a client can be opened and closed as many times as needed.
	}

	// read a single 16-bit holding register at address 100
	var res []uint16
	res, err = client.ReadRegisters(0, 4, modbus.HOLDING_REGISTER)
	if err != nil {
		// error out
	} else {
		for idx := range res {
			fmt.Printf("0x%04x\t%-5v : 0x%04x\t%v\n",
				0+uint16(idx),
				0+uint16(idx),
				res[idx], res[idx])

		}
	}
}
