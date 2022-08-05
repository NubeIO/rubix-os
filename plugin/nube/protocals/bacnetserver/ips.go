package main

import "github.com/NubeIO/lib-networking/networking"

var nets = networking.New()

func (inst *Instance) getNetworkByIface(iface string) (networking.NetworkInterfaces, error) {
	return nets.GetNetworkByIface(iface)
}

func (inst *Instance) getIp(iface string) (network string) {
	byIface, err := inst.getNetworkByIface(iface)
	if err != nil {
		return ""
	}
	return byIface.IP
}
