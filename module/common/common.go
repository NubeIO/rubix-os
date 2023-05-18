package common

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

const (
	NetworksURL    = "/networks"
	DevicesURL     = "/devices"
	PointsURL      = "/points"
	PointsWriteURL = "/points/write/:uuid"
)

func GetFlowNetworkNames(fns []*model.FlowNetwork) []string {
	fnsNames := make([]string, 0)
	for _, fn := range fns {
		fnsNames = append(fnsNames, fn.Name)
	}
	return fnsNames
}
