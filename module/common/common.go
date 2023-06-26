package common

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

const (
	NetworksURL    = "networks"
	DevicesURL     = "devices"
	PointsURL      = "points"
	PointsWriteURL = "points/write"
)

type PointWriteResponse struct {
	Point                model.Point `json:"point"`
	IsPresentValueChange bool        `json:"is_present_value_change"`
	IsWriteValueChange   bool        `json:"is_write_value_change"`
	IsPriorityChanged    bool        `json:"is_priority_changed"`
}
