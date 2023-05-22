package system

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-dhcpd/dhcpd"
	"github.com/NubeIO/lib-networking/networking"
)

type DHCPPortExists struct {
	IsDHCP          bool   `json:"is_dhcp"`
	InterfaceExists bool   `json:"interface_exists"`
	Error           string `json:"error"`
}

func (inst *System) DHCPPortExists(body NetworkingBody) (*DHCPPortExists, error) {
	resp := &DHCPPortExists{}
	var foundPort bool
	isDHCP, err := inst.dhcp.Exists(body.PortName)
	if err != nil {
		resp.Error = err.Error()
		return nil, err
	}
	ifaces, err := networking.New().GetInterfacesNames()
	if err != nil {
		resp.Error = err.Error()
		return nil, err
	}
	for _, name := range ifaces.Names {
		if body.PortName == name {
			foundPort = true
		}
	}
	resp.IsDHCP = isDHCP
	resp.InterfaceExists = foundPort
	return resp, nil
}

func (inst *System) DHCPSetAsAuto(body NetworkingBody) (*Message, error) {
	ok, err := inst.dhcp.SetAsAuto(body.PortName)
	if err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("was not able :%s to auto", body.PortName)
	if ok {
		msg = fmt.Sprintf("was able to set interface :%s to auto", body.PortName)
	} else {
		return nil, errors.New(fmt.Sprintf("was not able :%s to auto", body.PortName))
	}
	return &Message{
		Message: msg,
	}, nil
}

func (inst *System) DHCPSetStaticIP(body *dhcpd.SetStaticIP) (string, error) {
	return inst.dhcp.SetStaticIP(body)
}
