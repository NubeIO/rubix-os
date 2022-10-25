package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/projects/galvintmv/chirpstackrest"
)

func (inst *Instance) GetChirpstackToken() (*chirpstackrest.ChirpstackToken, error) {
	inst.tmvDebugMsg("GetChirpstackToken()")
	host := inst.config.Job.ChirpstackHost
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.ChirpstackPort
	if port == 0 {
		port = 8080
	}
	rest := chirpstackrest.NewNoAuth(host, int(port))
	token, err := rest.GetChirpstackToken()
	if err != nil {
		inst.tmvErrorMsg(err)
	}
	if err != nil {
		return nil, errors.New("could not get chirpstack token")
	}
	return token, nil
}

func (inst *Instance) GetChirpstackDeviceProfileUUID(chirpstackToken string) (string, error) {
	inst.tmvDebugMsg("GetChirpstackDeviceProfileUUID()")
	host := inst.config.Job.ChirpstackHost
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.ChirpstackPort
	if port == 0 {
		port = 8080
	}
	rest := chirpstackrest.NewNoAuth(host, int(port))
	deviceProfileArray, err := rest.GetChirpstackDeviceProfileUUID(chirpstackToken)
	if err != nil {
		inst.tmvErrorMsg(err)
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

func (inst *Instance) AddChirpstackDevice(chirpstackAppNum, modbusAddress int, deviceName, lorawanDeviceEUI, chirpstackDeviceProfileUUID string) (*interface{}, error) {
	inst.tmvDebugMsg("AddChirpstackDevice()")
	host := inst.config.Job.ChirpstackHost
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.ChirpstackPort
	if port == 0 {
		port = 8080
	}
	rest := chirpstackrest.NewNoAuth(host, int(port))
	response, err := rest.AddChirpstackDevice(chirpstackAppNum, modbusAddress, deviceName, lorawanDeviceEUI, chirpstackDeviceProfileUUID)
	if err != nil {
		inst.tmvErrorMsg(err)
	}
	if err != nil {
		return nil, errors.New("could not get chirpstack token")
	}
	return response, nil
}
