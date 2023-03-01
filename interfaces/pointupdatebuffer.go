package interfaces

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type PointUpdateBuffer struct {
	UUID  string       `json:"uuid"`
	Body  *model.Point `json:"body"`
	Point *model.Point `json:"point"`
}
