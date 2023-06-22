package database

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/NubeIO/rubix-os/utils/integer"
	"github.com/NubeIO/rubix-os/utils/priorityarray"
	"strings"
)

// updatePriority it updates priority array of point model
// it attaches the point model fields values for updating it on its parent function
func (d *GormDatabase) updatePriority(pointModel *model.Point, priority *map[string]*float64) (
	*model.Point, *map[string]*float64, *float64, *float64, bool) {
	isPriorityChanged := false
	var presentValue *float64
	var writeValue *float64
	priorityMap := priority
	presentValueFromPriority := pointModel.PointPriorityArrayMode != model.ReadOnlyNoPriorityArrayRequired &&
		pointModel.PointPriorityArrayMode != model.PriorityArrayToWriteValue
	// These values are not required for model.ReadOnlyNoPriorityArrayRequired
	if pointModel.PointPriorityArrayMode == model.ReadOnlyNoPriorityArrayRequired {
		pointModel.CurrentPriority = nil
		pointModel.WriteValueOriginal = nil

		pointModel.Priority.P1 = nil
		pointModel.Priority.P2 = nil
		pointModel.Priority.P3 = nil
		pointModel.Priority.P4 = nil
		pointModel.Priority.P5 = nil
		pointModel.Priority.P6 = nil
		pointModel.Priority.P7 = nil
		pointModel.Priority.P8 = nil
		pointModel.Priority.P9 = nil
		pointModel.Priority.P10 = nil
		pointModel.Priority.P11 = nil
		pointModel.Priority.P12 = nil
		pointModel.Priority.P13 = nil
		pointModel.Priority.P14 = nil
		pointModel.Priority.P15 = nil
		pointModel.Priority.P16 = nil
		d.DB.Model(&model.Priority{}).Where("point_uuid = ?", pointModel.UUID).Updates(&pointModel.Priority)
	}

	if priority != nil {
		pm, highestValue, currentPriority, doesPriorityExist, ipc :=
			priorityarray.ParsePriority(pointModel.Priority, priority, boolean.IsTrue(pointModel.IsTypeBool))
		priorityMap = pm
		isPriorityChanged = ipc
		if doesPriorityExist {
			if currentPriority == nil && highestValue == nil && !float.IsNil(pointModel.Fallback) {
				pointModel.Priority.P16 = float.New(*pointModel.Fallback)
				if boolean.IsTrue(pointModel.IsTypeBool) {
					pointModel.Priority.P16 = float.EvalAsBoolOnlyOneIsTrue(pointModel.Priority.P16)
				}
				priorityMapTemp := map[string]*float64{"_16": pointModel.Fallback}
				priorityMap = &priorityMapTemp
				currentPriority = integer.New(16)
				highestValue = float.New(*pointModel.Priority.P16)
			}
			if priorityMap != nil {
				pointModel.CurrentPriority = currentPriority
				pointModel.WriteValueOriginal = highestValue
				presentValue = highestValue
				writeValue = highestValue
			}
		}
		priorityMapToPatch_ := priorityMapToPatch(priorityMap)
		d.DB.Model(&pointModel.Priority).Where("point_uuid = ?", pointModel.UUID).Updates(&priorityMapToPatch_)
	}
	if !presentValueFromPriority {
		// presentValue will be OriginalValue if PointPriorityArrayMode is PriorityArrayToWriteValue or
		// ReadOnlyNoPriorityArrayRequired
		presentValue = pointModel.OriginalValue
	}
	return pointModel, priorityMap, presentValue, writeValue, isPriorityChanged
}

func priorityMapToPatch(priorityMap *map[string]*float64) map[string]interface{} {
	priorityMapToPatch_ := map[string]interface{}{}
	if priorityMap != nil {
		for k, v := range *priorityMap {
			priorityMapToPatch_[fmt.Sprintf("P%s", strings.Replace(k, "_", "", -1))] = v
		}
	}
	return priorityMapToPatch_
}

func (d *GormDatabase) GetPointsForPostgresSync() ([]*interfaces.PointForPostgresSync, error) {
	var pointsForPostgresModel []*interfaces.PointForPostgresSync
	query := d.DB.Table("points").
		Select("points.source_uuid AS uuid, points.name, networks.source_uuid AS network_uuid, networks.name AS network_name, " +
			"networks.global_uuid, devices.source_uuid AS device_uuid, devices.name AS device_name, hosts.uuid AS host_uuid, hosts.name AS host_name, " +
			"groups.uuid AS group_uuid, groups.name AS group_name, locations.uuid AS location_uuid, locations.name AS " +
			"location_name").
		Joins("INNER JOIN devices ON devices.uuid = points.device_uuid").
		Joins("INNER JOIN networks ON networks.uuid = devices.network_uuid").
		Joins("INNER JOIN hosts ON hosts.uuid = networks.host_uuid").
		Joins("INNER JOIN groups ON groups.uuid = hosts.group_uuid").
		Joins("INNER JOIN locations ON locations.uuid = groups.location_uuid").
		Scan(&pointsForPostgresModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return pointsForPostgresModel, nil
}

func (d *GormDatabase) GetPointsTagsForPostgresSync() ([]*interfaces.PointTagForPostgresSync, error) {
	var pointTagsForPostgresModel []*interfaces.PointTagForPostgresSync
	query := d.DB.Table("points_tags").
		Select("points.source_uuid AS point_uuid, points_tags.tag_tag AS tag").
		Joins("INNER JOIN points ON points.uuid = points_tags.point_uuid").
		Where("IFNULL(points.source_uuid,'') != ''").
		Scan(&pointTagsForPostgresModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return pointTagsForPostgresModel, nil
}
