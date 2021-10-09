package main

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)

// wizard
type wizard struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

//wizard make a network/dev/pnt
func (i *Instance) wizard(body wizard) (string, error) {
	ip := "192.168.15.101"
	if body.IP != "" {
		ip = body.IP
	}
	p := 5000
	if body.Port != 0 {
		p = body.Port
	}

	var net model.Network
	net.Name = "edge28"
	net.TransportType = model.TransType.IP
	net.PluginPath = "edge28"

	network, err := i.db.CreateNetwork(&net)
	if err != nil {
		return "", err
	}
	if network.UUID == "" {
		return "", errors.New("failed to create a new network")
	}
	var dev model.Device
	dev.NetworkUUID = network.UUID
	dev.Name = "edge28"
	dev.CommonIP.Host = ip
	dev.CommonIP.Port = p
	dev.PollDelayPointsMS = 5000

	device, err := i.db.CreateDevice(&dev)
	if err != nil {
		return "", err
	}
	if device.UUID == "" {
		return "", errors.New("failed to create a new device")
	}

	var pnt model.Point
	pnt.DeviceUUID = device.UUID

	for _, e := range pointsAll() {
		pName := fmt.Sprintf("edge28_%s", e)
		pnt.Name = pName
		pnt.Description = pName
		pnt.IoID = e
		pnt.IoType = model.IOType.RAW
		point, err := i.db.CreatePoint(&pnt)
		if err != nil {
			return "", err
		}
		if point.UUID == "" {
			return "", errors.New("failed to create a new point")
		}

	}
	return "pass: added network and points", nil
}
