package system_model

import (
	"github.com/NubeIO/flow-framework/plugin/defaults"
)

type Point struct {
	ObjectType struct {
		Options  []string `json:"options" default:"[\"analogInput\",\"analogOutput\",\"analogValue\",\"binaryInput\",\"binaryOutput\",\"binaryValue\"]"`
		Type     string   `json:"type" default:"array"`
		Required bool     `json:"required" default:"true"`
	} `json:"object_type"`
}

func GetPointSchema() *Point {
	point := &Point{}
	defaults.Set(point)
	return point
}
