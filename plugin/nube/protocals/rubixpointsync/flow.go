package main

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
)

func (inst *Instance) GetRequiredFFNetworks(requiredNetworks []string) ([]*model.Network, error) {
	var networksArray []*model.Network
	if requiredNetworks == nil || len(requiredNetworks) == 0 {
		inst.rubixpointsyncDebugMsg("no network name was specified to sync")
		requiredNetworks = []string{"system"}
	}
	for _, requiredNetwork := range requiredNetworks {
		net, _ := inst.db.GetNetworkByName(requiredNetwork, args.Args{WithDevices: true, WithPoints: true})
		networksArray = append(networksArray, net)
	}
	return networksArray, nil
}
