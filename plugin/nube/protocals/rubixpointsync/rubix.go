package main

import (
	"errors"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/rubixpointsync/rubixrest"
)

func (inst *Instance) GetRubixNetworks() (*[]rubixrest.RubixNet, error) {
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1515
	}
	rest := rubixrest.NewNoAuth(host, int(port))
	rubixNetsArray, err := rest.GetAllPoints()
	if err != nil {
		return nil, errors.New("no rubix points found")
	}
	return rubixNetsArray, nil
}

func (inst *Instance) CreateNewRubixPoint(pointName, deviceUUID string) (*rubixrest.RubixPnt, error) {
	inst.rubixpointsyncDebugMsg("CreateNewRubixPoint()", pointName)
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1515
	}
	rest := rubixrest.NewNoAuth(host, int(port))
	rubixPoint, err := rest.CreateNewRubixPoint(pointName, deviceUUID)
	if err != nil {
		inst.rubixpointsyncErrorMsg("CreateNewRubixPoint() err: ", err)
		return nil, errors.New("could not create rubix point")
	}
	return rubixPoint, nil
}

func (inst *Instance) WriteRubixPointByPathNames(networkName, deviceName, pointName string, writeValue *float64) (*rubixrest.RubixPnt, error) {
	inst.rubixpointsyncDebugMsg("WriteRubixPointByPathNames()", networkName, deviceName, pointName, writeValue)
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1515
	}
	rest := rubixrest.NewNoAuth(host, int(port))
	rubixPoint, err := rest.WriteRubixPointByPathNames(networkName, deviceName, pointName, writeValue)
	if err != nil {
		inst.rubixpointsyncErrorMsg("WriteRubixPointByPathNames() could not write rubix point. err:", err)
		return nil, errors.New("could not write rubix point")
	}
	return rubixPoint, nil
}

func (inst *Instance) WriteRubixPointByUUID(pointUUID string, writeValue *float64) (*rubixrest.RubixPnt, error) {
	inst.rubixpointsyncDebugMsg("WriteRubixPointByUUID()", pointUUID, writeValue)
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1515
	}
	rest := rubixrest.NewNoAuth(host, int(port))
	rubixPoint, err := rest.WriteRubixPointByUUID(pointUUID, writeValue)
	if err != nil {
		inst.rubixpointsyncErrorMsg("WriteRubixPointByPathNames() could not write rubix point. err:", err)
		return nil, errors.New("could not write rubix point")
	}
	return rubixPoint, nil
}

func (inst *Instance) RubixPointExistsInNetworkArray(checkNetworks *[]rubixrest.RubixNet, networkName, deviceName, pointName string) (netExists bool, devExists bool, pntExists bool, netUUID, devUUID, pntUUID string, err error) {
	netExists = false
	devExists = false
	pntExists = false
	if checkNetworks != nil {
		for _, net := range *checkNetworks {
			if inst.config.Job.RequireNetworkMatch && net.Name != networkName {
				continue
			}
			netUUID = net.UUID
			if net.Name == networkName {
				netExists = true
			}
			if net.Devices == nil {
				continue
			}
			for _, dev := range net.Devices {
				if dev == nil || dev.Name != deviceName {
					continue
				}
				devUUID = dev.UUID
				devExists = true
				if dev.Points == nil {
					continue
				}
				for _, pnt := range dev.Points {
					if pnt == nil || pnt.Name != pointName {
						continue
					} else { // Found the point
						pntExists = true
						pntUUID = pnt.UUID
						return netExists, devExists, pntExists, netUUID, devUUID, pntUUID, nil
					}
				}
			}
		}
	}
	return netExists, devExists, pntExists, netUUID, devUUID, pntUUID, errors.New("point couldn't be found")
}

func (inst *Instance) CreateNewRubixDevice(deviceName, networkUUID string) (*rubixrest.RubixDev, error) {
	inst.rubixpointsyncDebugMsg("CreateNewRubixDevice()")
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1515
	}
	rest := rubixrest.NewNoAuth(host, int(port))
	rubixDevice, err := rest.CreateNewRubixDevice(deviceName, networkUUID)
	if err != nil {
		inst.rubixpointsyncErrorMsg("CreateNewRubixDevice() err: ", err)
		return nil, errors.New("could not create rubix device")
	}
	return rubixDevice, nil
}

func (inst *Instance) CreateNewRubixNetwork(netName string) (*rubixrest.RubixNet, error) {
	inst.rubixpointsyncDebugMsg("CreateNewRubixNetwork()")
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1515
	}
	rest := rubixrest.NewNoAuth(host, int(port))
	rubixNetwork, err := rest.CreateNewRubixNetwork(netName)
	if err != nil {
		inst.rubixpointsyncErrorMsg("CreateNewRubixNetwork() err: ", err)
		return nil, errors.New("could not create rubix network")
	}
	return rubixNetwork, nil
}
