package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/maplora/legacylorarest"
)

func (inst *Instance) GetLegacyLoRaNetwork() (*legacylorarest.LoRaNet, error) {
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1919
	}
	rest := legacylorarest.NewNoAuth(host, int(port))
	loraNet, err := rest.GetLegacyLoRaNetwork()
	if err != nil || loraNet == nil {
		inst.maploraErrorMsg("no legacy lora network found. err: ", err)
		return nil, errors.New(fmt.Sprint("no legacy lora network found. err:", err))
	}
	return loraNet, nil
}

func (inst *Instance) GetLegacyLoRaDevices() (*[]legacylorarest.LoRaDev, error) {
	host := inst.config.Job.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.Port
	if port == 0 {
		port = 1919
	}
	rest := legacylorarest.NewNoAuth(host, int(port))
	loraDevsArray, err := rest.GetLegacyLoRaDevices()
	if err != nil {
		return nil, errors.New("no legacy lora devices found")
	}
	return loraDevsArray, nil
}

func (inst *Instance) ConvertLegacyLoRaModelToFFModel(legacyModelString string) (string, error) {
	switch legacyModelString {
	case "DROPLET_TH":
		return "TH", nil
	case "DROPLET_THL":
		return "THL", nil
	case "DROPLET_THLM":
		return "THLM", nil
	case "MICRO_EDGE":
		return "MicroEdge", nil
	}
	return "", errors.New("unrecognized lora device model")
}
