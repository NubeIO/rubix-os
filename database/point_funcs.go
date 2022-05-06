package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (d *GormDatabase) GetPointByName(networkName, deviceName, pointName string) (*model.Point, error) {
	var pointModel *model.Point
	net, err := d.GetNetworkByName(networkName, api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return nil, errors.New("failed to find a network with that name")
	}
	deviceExist := false
	pointExist := false
	for _, device := range net.Devices {
		if device.Name == deviceName {
			deviceExist = true
			for _, p := range device.Points {
				if p.Name == pointName {
					pointExist = true
					pointModel = p
					break
				}
			}
		}
	}
	if !deviceExist {
		return nil, errors.New("failed to find a device with that name")
	}
	if !pointExist {
		return nil, errors.New("found device but failed to find a point with that name")
	}
	return pointModel, nil
}

// PointWriteByName TODO: functions calling  d.PointWrite(point.UUID, body, fromPlugin) should be routed via plugin!!
func (d *GormDatabase) PointWriteByName(networkName, deviceName, pointName string, body *model.PointWriter, fromPlugin bool) (*model.Point, error) {
	point, err := d.GetPointByName(networkName, deviceName, pointName)
	if err != nil {
		return nil, err
	}
	write, err := d.PointWrite(point.UUID, body, fromPlugin)
	if err != nil {
		return nil, err
	}
	return write, nil
}

func (d *GormDatabase) GetOnePointByArgs(args api.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.First(&pointModel).Error; err != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

// updatePriority it updates priority array of point model
// it attaches the point model fields values for updating it on it's parent function
func (d *GormDatabase) updatePriority(pointModel *model.Point, priority *map[string]*float64) (*model.Point, *map[string]*float64, *float64) {
	var presentValue *float64
	priorityMap := priority
	presentValueFromPriority := pointModel.PointPriorityArrayMode != model.ReadOnlyNoPriorityArrayRequired && pointModel.PointPriorityArrayMode != model.PriorityArrayToWriteValue
	// These values are not required for model.ReadOnlyNoPriorityArrayRequired
	if pointModel.PointPriorityArrayMode == model.ReadOnlyNoPriorityArrayRequired {
		pointModel.CurrentPriority = nil
		pointModel.WriteValue = nil
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
		// override priorityMap
		priorityMap, highestValue, currentPriority, doesPriorityExist := priorityarray.ParsePriority(pointModel.Priority, priority)
		if doesPriorityExist {
			if priorityMap != nil {
				pointModel.CurrentPriority = currentPriority
				pointModel.WriteValueOriginal = highestValue
				writeValue, err := pointEval(highestValue, pointModel.MathOnWriteValue)
				if err != nil {
					log.Errorln("point.db parsePriority() error on run point MathOnWriteValue error:", err)
				} else {
					pointModel.WriteValue = writeValue
				}
				presentValue = highestValue
			} else if !utils.FloatIsNilCheck(pointModel.Fallback) || currentPriority == nil {
				pointModel.Priority.P16 = utils.NewFloat64(*pointModel.Fallback)
				priorityMapTemp := map[string]*float64{"_16": pointModel.Fallback}
				priorityMap = &priorityMapTemp
				pointModel.CurrentPriority = utils.NewInt(16)
				pointModel.WriteValueOriginal = utils.NewFloat64(*pointModel.Priority.P16)
				writeValue, err := pointEval(pointModel.Priority.P16, pointModel.MathOnWriteValue)
				if err != nil {
					log.Errorln("point.db parsePriority() error on run point MathOnWriteValue error:", err)
				} else {
					pointModel.WriteValue = writeValue
				}
				presentValue = pointModel.Fallback // only update presentValue if required by PointPriorityArrayMode
			}
		}
		priorityMapToPatch := d.priorityMapToPatch(priorityMap)
		d.DB.Model(&model.Priority{}).Where("point_uuid = ?", pointModel.UUID).Updates(&priorityMapToPatch)
	}
	if !presentValueFromPriority {
		// presentValue will be OriginalValue if PointPriorityArrayMode is PriorityArrayToWriteValue or
		// ReadOnlyNoPriorityArrayRequired
		presentValue = pointModel.OriginalValue
	}
	return pointModel, priorityMap, presentValue
}

func (d *GormDatabase) priorityMapToPatch(priorityMap *map[string]*float64) map[string]interface{} {
	priorityMapToPatch := map[string]interface{}{}
	if priorityMap != nil {
		for k, v := range *priorityMap {
			priorityMapToPatch[fmt.Sprintf("P%s", strings.Replace(k, "_", "", -1))] = v
		}
	}
	return priorityMapToPatch
}
