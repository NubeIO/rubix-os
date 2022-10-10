package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *Instance) GetFFNetworks(pluginsArray []string) ([]*model.Network, error) {
	inst.rubixpointsyncDebugMsg("GetFFNetworks()")
	var networksArray []*model.Network
	if pluginsArray == nil || len(pluginsArray) == 0 {
		pluginsArray = []string{"system"}
	}
	for _, plugin := range pluginsArray {
		nets, err := inst.db.GetNetworksByPluginName(plugin, api.Args{WithDevices: true, WithPoints: true})
		if err != nil {
			continue
		}
		for _, net := range nets {
			if net.Devices != nil {
				networksArray = append(networksArray, net)
			}
		}

		/*
			for _, net := range nets {
				inst.rubixpointsyncDebugMsg("GetFFPointValues() Net: ", net.Name)
				for _, dev := range net.Devices {
					for _, pnt := range dev.Points {
						point, _ := inst.db.GetPoint(pnt.UUID, api.Args{WithTags: true})
						// inst.rubixpointsyncDebugMsg(fmt.Sprintf("GetFFPointValues() point: %+v", point))
						if point.PresentValue != nil {
							tagMap := make(map[string]string)
							tagMap["plugin_name"] = "lorawan"
							tagMap["network_name"] = net.Name
							tagMap["network_uuid"] = net.UUID
							tagMap["device_name"] = dev.Name
							tagMap["device_uuid"] = dev.UUID
							tagMap["point_name"] = point.Name
							tagMap["point_uuid"] = point.UUID

							pointHistory := History{
								UUID:      point.UUID,
								Value:     float.NonNil(point.PresentValue),
								Timestamp: nowTimestamp,
								Tags:      tagMap,
							}
							inst.rubixpointsyncDebugMsg(fmt.Sprintf("GetHistoryValues() history: %+v", pointHistory))
							historyArray = append(historyArray, &pointHistory)
						}
					}
				}
			}

		*/
	}
	return networksArray, nil
}
