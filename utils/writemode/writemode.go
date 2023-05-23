package writemode

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/boolean"
)

func SetPriorityArrayModeBasedOnWriteMode(pnt *model.Point) bool {
	switch pnt.WriteMode {
	case model.ReadOnce, model.ReadOnly:
		pnt.PointPriorityArrayMode = model.ReadOnlyNoPriorityArrayRequired
		return true
	case model.WriteOnce, model.WriteOnceReadOnce, model.WriteAlways, model.WriteOnceThenRead, model.WriteAndMaintain:
		pnt.PointPriorityArrayMode = model.PriorityArrayToWriteValue
		return true
	}
	return false
}

func IsWriteable(writeMode model.WriteMode) bool {
	switch writeMode {
	case model.ReadOnce, model.ReadOnly:
		return false
	case model.WriteOnce, model.WriteOnceReadOnce, model.WriteAlways, model.WriteOnceThenRead, model.WriteAndMaintain:
		return true
	default:
		return false
	}
}

func ResetWriteableProperties(point *model.Point) *model.Point {
	point.WriteValueOriginal = nil
	point.WriteValue = nil
	point.WritePriority = nil
	point.CurrentPriority = nil
	point.EnableWriteable = boolean.NewFalse()
	point.WritePollRequired = boolean.NewFalse()
	return point
}
