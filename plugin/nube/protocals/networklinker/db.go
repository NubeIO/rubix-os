package main

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/boolean"
)

func (inst *Instance) createNetwork(body *model.Network) (*model.Network, error) {
	if body.Name == "" {
		body.Name = strings.Replace(body.AddressUUID, INTERNAL_SEPARATOR, UI_SEPARATOR, 1)
	}
	netUUIDs := strings.Split(body.AddressUUID, INTERNAL_SEPARATOR)
	net1, _ := inst.db.GetNetwork(netUUIDs[0], api.Args{WithDevices: true})
	net2, _ := inst.db.GetNetwork(netUUIDs[1], api.Args{WithDevices: true})
	if inst.networkIsWriter(net1) && inst.networkIsWriter(net2) {
		return nil, errors.New("both networks cannot be \"writers\"")
	}
	body.AddressUUID = fmt.Sprint(net1.UUID, INTERNAL_SEPARATOR, net2.UUID)
	body, err := inst.db.CreateNetwork(body)
	if err != nil {
		return nil, err
	}

	for _, dev1 := range net1.Devices {
		for _, dev2 := range net2.Devices {
			if dev1.Name == dev2.Name {
				d := model.Device{
					Name:        dev1.Name,
					NetworkUUID: body.UUID,
				}
				inst.createDevice(&d, dev1, dev2, net1, net2)
				break
			}
		}
	}
	return body, err
}

func (inst *Instance) createDevice(body *model.Device, dev1 *model.Device, dev2 *model.Device, net1 *model.Network, net2 *model.Network) (*model.Device, error) {
	if net1 == nil || net2 == nil {
		linkNet, _ := inst.db.GetNetwork(body.NetworkUUID, api.Args{})
		netUUIDs := strings.Split(linkNet.AddressUUID, INTERNAL_SEPARATOR)
		net1, _ = inst.db.GetNetwork(netUUIDs[0], api.Args{WithDevices: true})
		net2, _ = inst.db.GetNetwork(netUUIDs[1], api.Args{WithDevices: true})
	}
	if dev1 == nil || dev2 == nil {
		devUUIDs := strings.Split(*body.AddressUUID, INTERNAL_SEPARATOR)
		for _, dev := range net1.Devices {
			if dev.UUID == devUUIDs[0] || dev.UUID == devUUIDs[1] {
				dev1 = dev
				break
			}
		}
		for _, dev := range net2.Devices {
			if dev.UUID == devUUIDs[0] || dev.UUID == devUUIDs[1] {
				dev2 = dev
				break
			}
		}
		if dev1 == nil || dev2 == nil {
			return nil, errors.New("device does not belong to correct network")
		}
	}
	addr := fmt.Sprint(dev1.UUID, INTERNAL_SEPARATOR, dev2.UUID)
	body.AddressUUID = &addr
	device, err := inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	err = inst.syncDevicePoints(device, dev1, dev2, net1, net2)
	return device, err
}

// Add points to link device for any points from dev1 and dev2.
//  Match points via name, else add as singlular point
//  Checks for any missing points
func (inst *Instance) syncDevicePoints(devLink *model.Device, dev1 *model.Device, dev2 *model.Device, net1 *model.Network, net2 *model.Network) error {
	devUUIDs := strings.Split(*devLink.AddressUUID, INTERNAL_SEPARATOR)
	if dev1 == nil {
		dev1, _ = inst.db.GetDevice(devUUIDs[0], api.Args{})
	}
	if dev2 == nil {
		dev2, _ = inst.db.GetDevice(devUUIDs[1], api.Args{})
	}
	if net1 == nil {
		net1, _ = inst.db.GetNetwork(dev1.NetworkUUID, api.Args{})
	}
	if net2 == nil {
		net2, _ = inst.db.GetNetwork(dev2.NetworkUUID, api.Args{})
	}

	existingLinkedPointUUIDs := make([][]string, len(devLink.Points))
	for i, point := range devLink.Points {
		existingLinkedPointUUIDs[i] = strings.Split(*point.AddressUUID, INTERNAL_SEPARATOR)
	}

	points1, _ := inst.db.GetPointsByDeviceUUID(dev1.UUID, api.Args{})
	points2, _ := inst.db.GetPointsByDeviceUUID(dev2.UUID, api.Args{})

	pointsFound2 := make([]bool, len(points2))
	for i := range pointsFound2 {
		pointsFound2[i] = false
	}
	// add all matched points AND dev1 points
	for _, p1 := range points1 {
		var p2uuid *string = nil
		var net2P *model.Network = nil
		if checkExistingPointLink(p1, existingLinkedPointUUIDs) {
			continue
		}
		// look for match points for dev1 in dev2
		for i2, p2 := range points2 {
			if p1.Name == p2.Name {
				p2uuid = &p2.UUID
				net2P = net2
				pointsFound2[i2] = true
				break
			}
		}
		// matched OR dev1-only points
		inst.createPointAndUpdate(devLink.UUID, p1.Name, &p1.UUID, p2uuid, net1, net2P)
	}
	// add all leftover dev2 points
	for i, p2 := range points2 {
		if !pointsFound2[i] {
			if !checkExistingPointLink(p2, existingLinkedPointUUIDs) {
				inst.createPointAndUpdate(devLink.UUID, p2.Name, &p2.UUID, nil, net2, nil)
			}
		}
	}
	return nil
}

func checkExistingPointLink(point *model.Point, linkUUIDs [][]string) bool {
	for _, uuidPair := range linkUUIDs {
		if point.UUID == uuidPair[0] || (len(uuidPair) > 1 && point.UUID == uuidPair[1]) {
			return true
		}
	}
	return false
}

func (inst *Instance) createPointAndUpdate(devUUID string, name string, uuid1 *string, uuid2 *string, net1 *model.Network, net2 *model.Network) model.Point {
	p := inst.createPoint(devUUID, name, uuid1, uuid2)
	inst.syncPoint(&p, net1, net2)
	return p
}

func (inst *Instance) relinkPointAndUpdate(point *model.Point, uuid1 *string, uuid2 *string, net1 *model.Network, net2 *model.Network) *model.Point {
	*point.AddressUUID = fmt.Sprintf("%s%s%s", *uuid1, INTERNAL_SEPARATOR, *uuid2)
	inst.db.UpdatePoint(point.UUID, point)
	inst.syncPoint(point, net1, net2)
	return point
}

func (inst *Instance) createPoint(devUUID string, name string, uuid1 *string, uuid2 *string) model.Point {
	var addr string
	if uuid2 != nil {
		addr = fmt.Sprintf("%s%s%s", *uuid1, INTERNAL_SEPARATOR, *uuid2)
	} else {
		addr = *uuid1
	}
	p := model.Point{
		DeviceUUID:   devUUID,
		Name:         name,
		AddressUUID:  &addr,
		CommonEnable: model.CommonEnable{Enable: boolean.NewTrue()},
	}
	inst.db.CreatePoint(&p, false)
	return p
}

func (inst *Instance) syncPoint(point *model.Point, net1 *model.Network, net2 *model.Network) *model.Point {
	pUUIDs := strings.Split(*point.AddressUUID, INTERNAL_SEPARATOR)
	selectedUUID := pUUIDs[0]
	if len(pUUIDs) > 1 {
		// select the "publish"/"non-write" network if multple networks
		if inst.networkIsWriter(net1) {
			selectedUUID = pUUIDs[1]
		}
	}
	return inst.syncPointSelected(point, selectedUUID)
}

func (inst *Instance) syncPointSelected(point *model.Point, linkedUUID string) *model.Point {
	origPoint, _ := inst.db.GetPoint(linkedUUID, api.Args{WithPriority: true})
	if origPoint.PresentValue == nil || (point.PresentValue != nil && (*point.PresentValue == *origPoint.PresentValue)) {
		return point
	}
	origPoint.UUID = point.UUID
	origPoint.AddressUUID = point.AddressUUID
	origPoint.DeviceUUID = point.DeviceUUID
	origPoint.Name = point.Name
	origPoint.Enable = point.Enable
	point, _ = inst.db.UpdatePoint(point.UUID, origPoint)
	return point
}

func (inst *Instance) networkIsWriter(net *model.Network) bool {
	return slices.Contains(inst.config.Writers, net.PluginPath)
}

func (inst *Instance) getWriterNetworkAndPoint(linkPointUUID string) (network *model.Network, pointUUID *string) {
	linkPoint, _ := inst.db.GetPoint(linkPointUUID, api.Args{})
	pointUUIDs := strings.Split(*linkPoint.AddressUUID, INTERNAL_SEPARATOR)
	network, _ = inst.db.GetNetworkByDeviceUUID(linkPoint.DeviceUUID, api.Args{})
	netUUIDs := strings.Split(network.AddressUUID, INTERNAL_SEPARATOR)
	network1, _ := inst.db.GetNetwork(netUUIDs[0], api.Args{})
	var newNet *model.Network = nil

	if len(pointUUIDs) == 1 {
		pointUUID = &pointUUIDs[0]
		point, _ := inst.db.GetPoint(*pointUUID, api.Args{})
		newNet, _ = inst.db.GetNetworkByDeviceUUID(point.DeviceUUID, api.Args{})
	} else if inst.networkIsWriter(network1) {
		pointUUID = &pointUUIDs[0]
		newNet = network1
	} else {
		network2, _ := inst.db.GetNetwork(netUUIDs[1], api.Args{})
		pointUUID = &pointUUIDs[1]
		newNet = network2
	}
	return newNet, pointUUID
}
