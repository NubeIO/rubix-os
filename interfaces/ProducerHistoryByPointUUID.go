package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type ProducerHistoryByPointUUID struct {
	PointUUID string `json:"point_uuid"`
	model.ProducerHistory
}
