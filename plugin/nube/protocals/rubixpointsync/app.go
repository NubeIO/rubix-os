package main

import "github.com/NubeIO/flow-framework/plugin/nube/protocals/rubixpointsync/rubixrest"

func (inst *Instance) SyncRubixToFF() (bool, error) {
	inst.rubixpointsyncDebugMsg("SyncRubixToFF()")

	rubixNets, err := inst.GetRubixNetworks()
	if err != nil {
		inst.rubixpointsyncErrorMsg(err)
	}

	ffNetworksArray, err := inst.GetRequiredFFNetworks(inst.config.Job.Networks)
	if err != nil {
		inst.rubixpointsyncErrorMsg(err)
	}

	rescanRubix := false
	if ffNetworksArray != nil {
		for _, net := range ffNetworksArray {
			if net == nil || net.Devices == nil {
				continue
			}
			for _, dev := range net.Devices {
				if dev == nil || dev.Points == nil {
					continue
				}
				for _, pnt := range dev.Points {
					if pnt == nil {
						continue
					}
					if rescanRubix {
						rubixNets, err = inst.GetRubixNetworks()
						if err != nil {
							inst.rubixpointsyncErrorMsg(err)
						}
					}
					netExists, devExists, pointExists, netUUID, devUUID, pntUUID, _ := inst.RubixPointExistsInNetworkArray(rubixNets, net.Name, dev.Name, pnt.Name)
					var newRubixPnt *rubixrest.RubixPnt
					newRubixPnt = nil
					if !pointExists && inst.config.Job.GenerateRubixPoints {
						if (inst.config.Job.RequireNetworkMatch && !netExists) || netUUID == "" {
							newRubixNet, err := inst.CreateNewRubixNetwork(net.Name)
							rescanRubix = true
							if err != nil {
								inst.rubixpointsyncErrorMsg("bad response from CreateNewRubixNetwork(), ", err)
								continue
							}
							netUUID = newRubixNet.UUID
						}
						if !devExists {
							newRubixDev, err := inst.CreateNewRubixDevice(dev.Name, netUUID)
							rescanRubix = true
							if err != nil {
								inst.rubixpointsyncErrorMsg("bad response from CreateNewRubixDevice(), ", err)
								continue
							}
							devUUID = newRubixDev.UUID

						}

						newRubixPnt, err = inst.CreateNewRubixPoint(pnt.Name, devUUID)
						rescanRubix = true
						if err != nil {
							inst.rubixpointsyncErrorMsg("bad response from CreateNewRubixPoint(), ", err)
							continue
						}
					}
					if newRubixPnt != nil {
						pntUUID = newRubixPnt.UUID
					}
					_, err = inst.WriteRubixPointByUUID(pntUUID, pnt.PresentValue)
					if err != nil {
						inst.rubixpointsyncErrorMsg("writePoint(): bad response from WriteRubixPoint(), ", err)
					}
				}
			}
		}
	}
	return true, nil
}

func (inst *Instance) SyncSingleRubixPointWithFF(netName, devName, pntName string, value *float64) error {
	inst.rubixpointsyncDebugMsg("SyncSingleRubixPointWithFF()")

	rubixNets, err := inst.GetRubixNetworks()
	if err != nil {
		inst.rubixpointsyncErrorMsg(err)
	}

	netExists, devExists, pointExists, _, _, pntUUID, _ := inst.RubixPointExistsInNetworkArray(rubixNets, netName, devName, pntName)

	if netExists && devExists && pointExists {
		_, err = inst.WriteRubixPointByUUID(pntUUID, value)
		if err != nil {
			inst.rubixpointsyncErrorMsg("SyncSingleRubixPointWithFF(): bad response from WriteRubixPointByUUID(), ", err)
		}
	}
	return nil
}
