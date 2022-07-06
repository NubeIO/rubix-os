package main

import (
	"github.com/NubeIO/flow-framework/utils/array"
	"go.bug.st/serial"
)

// listSerialPorts list all serial ports on host
func (inst *Instance) getAllScheduleData() (*array.Array, error) {

	ports, err := serial.GetPortsList()
	p := array.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}
