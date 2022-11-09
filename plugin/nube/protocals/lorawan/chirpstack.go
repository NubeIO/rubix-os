package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/projects/galvintmv/chirpstackrest"
)

func (inst *Instance) GetChirpstackToken(user, pass string) (*chirpstackrest.ChirpstackToken, error) {
	inst.lorawanDebugMsg("GetChirpstackToken()")
	// host := inst.config.Job.ChirpstackHost
	host := "0.0.0.0"
	if host == "" {
		host = "0.0.0.0"
	}
	// port := inst.config.Job.ChirpstackPort
	port := 8080
	if port <= 0 {
		port = 8080
	}
	rest := chirpstackrest.NewNoAuth(host, int(port))
	token, err := rest.GetChirpstackToken(user, pass)
	if err != nil {
		inst.lorawanErrorMsg(err)
	}
	if err != nil {
		return nil, errors.New("could not get chirpstack token")
	}
	return token, nil
}

func (inst *Instance) GetChirpstackDeviceProfileUUID(chirpstackToken string) (string, error) {
	inst.lorawanDebugMsg("GetChirpstackDeviceProfileUUID()")
	// host := inst.config.Job.ChirpstackHost
	host := "0.0.0.0"
	if host == "" {
		host = "0.0.0.0"
	}
	// port := inst.config.Job.ChirpstackPort
	port := 8080
	if port <= 0 {
		port = 8080
	}
	rest := chirpstackrest.NewNoAuth(host, int(port))
	deviceProfileArray, err := rest.GetChirpstackDeviceProfileUUID(chirpstackToken)
	if err != nil {
		inst.lorawanErrorMsg(err)
		return "", errors.New("could not get chirpstack device profiles")
	}
	for _, deviceProfile := range deviceProfileArray {
		fmt.Println(deviceProfile.Name)
		if deviceProfile.Name == "R-IO-OTAA" {
			return deviceProfile.ID, nil
		}
	}
	return "", errors.New("could not find 'R-IO-OTAA' device profile UUID")
}

func (inst *Instance) AddChirpstackDevice(chirpstackAppNum, modbusAddress int, deviceName, lorawanDeviceEUI, chirpstackDeviceProfileUUID, token string) error {
	inst.lorawanDebugMsg("AddChirpstackDevice()")
	// host := inst.config.Job.ChirpstackHost
	host := "0.0.0.0"
	if host == "" {
		host = "0.0.0.0"
	}
	// port := inst.config.Job.ChirpstackPort
	port := 8080
	if port <= 0 {
		port = 8080
	}
	rest := chirpstackrest.NewNoAuth(host, int(port))
	err := rest.AddChirpstackDevice(chirpstackAppNum, modbusAddress, deviceName, lorawanDeviceEUI, chirpstackDeviceProfileUUID, token)
	if err != nil {
		return err
	}
	return nil
}

func (inst *Instance) ActivateChirpstackDevice(applicationKey, lorawanDeviceEUI, token, lorawanNetworkKey string) error {
	inst.lorawanDebugMsg("ActivateChirpstackDevice(): ", lorawanDeviceEUI)
	// host := inst.config.Job.ChirpstackHost
	host := "0.0.0.0"
	if host == "" {
		host = "0.0.0.0"
	}
	// port := inst.config.Job.ChirpstackPort
	port := 8080
	if port <= 0 {
		port = 8080
	}
	rest := chirpstackrest.NewNoAuth(host, int(port))
	err := rest.ActivateChirpstackDevice(applicationKey, lorawanDeviceEUI, token, lorawanNetworkKey)
	if err != nil {
		return err
	}
	return nil
}
