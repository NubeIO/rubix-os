package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *Instance) GetDataRateFromDevice(dev *model.Device) (float64, error) {
	if dev.Points == nil {
		return 5, errors.New(fmt.Sprint("no points on device:", dev.Name))
	}
	for _, pnt := range dev.Points {
		if pnt.Name == "data-rate" && pnt.PresentValue != nil {
			return *pnt.PresentValue, nil
		}
	}
	return 5, errors.New(fmt.Sprint("couldn't find `data-rate` point on device:", dev.Name))
}

func (inst *Instance) GetSensorIDFromDeviceDescription(dev *model.Device) (string, error) {
	formattedStringDescription := dev.Description
	if formattedStringDescription != "" {
		return formattedStringDescription, nil
	} else {
		return "", errors.New("no sensor id field found in device description")
	}
}

func printHistoriesByGateway(histories InauroPackagedSensorHistoriesByGateway) {
	for gatewayID, historiesByDevice := range histories {
		fmt.Println("Gateway ID:", gatewayID)
		printHistoriesByDevice(historiesByDevice)
		fmt.Println()
	}
}

func printHistoriesByDevice(histories InauroPackagedSensorHistoriesByDevice) {
	for deviceID, historiesByTimestamp := range histories {
		fmt.Println("Device ID:", deviceID)
		printHistoriesByTimestamp(historiesByTimestamp)
		fmt.Println()
	}
}

func printHistoriesByTimestamp(histories InauroPackagedSensorHistoriesByTimestamp) {
	for timestamp, payload := range histories {
		fmt.Println("Timestamp:", timestamp)
		fmt.Println("Payload:", payload)
		fmt.Println()
	}
}
