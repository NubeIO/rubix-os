package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type MqttPoint struct {
	NetworkName string             `json:"network_name,omitempty"`
	DeviceName  string             `json:"device_name,omitempty"`
	PointName   string             `json:"point_name,omitempty"`
	PointUUID   string             `json:"point_uuid,omitempty"`
	Priority    *model.PointWriter `json:"priority,omitempty"`
}
