package system

import (
	"errors"
	"fmt"
	"os/exec"
)

type NetworkingBody struct {
	PortName string `json:"port_name"`
}

func (inst *System) RestartNetworking() (*Message, error) {
	cmd := exec.Command("systemctl", "restart", "networking.service")
	output, err := cmd.Output()
	cleanCommand(string(output), cmd, err, debug)
	if err != nil {
		return nil, err
	}
	return &Message{
		Message: "restarted ok",
	}, err
}

func (inst *System) InterfaceUpDown(port NetworkingBody) (*Message, error) {
	_, err := inst.interfaceUpDown(port.PortName, false) // set down
	if err != nil {
		return nil, err
	}
	_, err = inst.interfaceUpDown(port.PortName, true) // set up
	if err != nil {
		return nil, err
	}
	return &Message{fmt.Sprintf("restart network: %s", port.PortName)}, nil

}

func (inst *System) InterfaceUp(port NetworkingBody) (*Message, error) {
	return inst.interfaceUpDown(port.PortName, true)
}

func (inst *System) InterfaceDown(port NetworkingBody) (*Message, error) {
	return inst.interfaceUpDown(port.PortName, false)
}

// interfaceUpDown ifconfig enp4s0 up
func (inst *System) interfaceUpDown(port string, up bool) (*Message, error) {
	names, err := nets.GetInterfacesNames()
	if err != nil {
		return nil, err
	}
	var found bool
	for _, s := range names.Names {
		if port == s {
			found = true
		}
	}
	if !found {
		return nil, errors.New(fmt.Sprintf("port %s was not found", port))
	}
	cmd := exec.Command("ifconfig", port, "down")
	msg := "disbaled"
	if up {
		msg = "enabled"
		cmd = exec.Command("ifconfig", port, "up")
	}
	output, err := cmd.Output()
	cleanCommand(string(output), cmd, err, debug)
	if err != nil {
		return nil, err
	}
	return &Message{
		Message: fmt.Sprintf("port %s is now %s", port, msg),
	}, err
}
