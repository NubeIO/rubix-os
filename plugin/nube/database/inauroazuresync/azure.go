package main

import (
	"context"
	"github.com/amenzhinsky/iothub/iotdevice"
	iotmqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
)

func (inst *Instance) newAzureMQTTClientByHostUUID(hostUUID string) (*iotdevice.Client, error) {
	// USE MQTT FROM https://github.com/amenzhinsky/iothub LIBRARY
	inst.inauroazuresyncDebugMsg("newAzureMQTTClientByHostUUID() hostUUID: ", hostUUID)
	sasString, _, _, _, err := inst.getAzureDetailsFromConfigByHostUUID(hostUUID)
	inst.inauroazuresyncDebugMsg("newAzureMQTTClient() sasString: ", sasString)
	if err != nil {
		return nil, err
	}
	c, err := iotdevice.NewFromConnectionString(iotmqtt.New(), sasString)
	if err != nil {
		return nil, err
	}
	// connect to the iothub
	if err = c.Connect(context.Background()); err != nil {
		return nil, err
	}
	inst.inauroazuresyncDebugMsg("newAzureMQTTClientByHostUUID() hostUUID: ", hostUUID, "  CONNECTED")
	return c, nil
}

func (inst *Instance) newAzureMQTTClientByGatewayDetails(gatewayDetails GatewayDetails) (*iotdevice.Client, error) {
	// USE MQTT FROM https://github.com/amenzhinsky/iothub LIBRARY
	inst.inauroazuresyncDebugMsg("newAzureMQTTClientByGatewayDetails() gateway: ", gatewayDetails.AzureDeviceId)
	sasString, err := inst.getSAS(gatewayDetails.AzureDeviceId, gatewayDetails.AzureDeviceKey)
	if err != nil {
		return nil, err
	}
	c, err := iotdevice.NewFromConnectionString(iotmqtt.New(), sasString)
	if err != nil {
		return nil, err
	}
	// connect to the iothub
	if err = c.Connect(context.Background()); err != nil {
		return nil, err
	}
	inst.inauroazuresyncDebugMsg("newAzureMQTTClientByGatewayDetails() gateway: ", gatewayDetails.AzureDeviceId, "  CONNECTED")
	return c, nil
}
