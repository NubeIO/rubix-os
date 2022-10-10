package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/rubixpointsync/rubixrest"
)

func (inst *Instance) GetRubixNetworks() error {
	inst.rubixpointsyncDebugMsg("GetRubixNetworks()")

	rest := rubixrest.NewNoAuth("192.168.1.30", 1515)
	rubixNetsArray, err := rest.GetAllPoints()
	if err != nil {
		return errors.New("no rubix points found")
	}
	inst.rubixpointsyncDebugMsg("GetRubixNetworks(): API Results rubixNets:")
	if rubixNetsArray != nil {
		inst.rubixpointsyncDebugMsg(fmt.Sprintf("NETWORK ARRAY %+v", *rubixNetsArray))
		for _, net := range *rubixNetsArray {
			inst.rubixpointsyncDebugMsg(fmt.Sprintf("NETWORK %+v", net))
			if net.Devices != nil {
				for _, dev := range net.Devices {
					inst.rubixpointsyncDebugMsg(fmt.Sprintf("DEVICE %+v", dev))
					if dev.Points != nil {
						for _, pnt := range dev.Points {
							inst.rubixpointsyncDebugMsg(fmt.Sprintf("POINT %+v", pnt))
						}
					}
				}
			}
		}
	}
	return nil
}
