package interfaces

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type State string

const (
	Created  State = "Created"
	Updating State = "Updating"
)

type PointUpdateBuffer struct {
	UUID  string       `json:"uuid"`
	Body  *model.Point `json:"body"`
	Point *model.Point `json:"point"`
	State State        `json:"state"`
}
