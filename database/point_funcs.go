package database

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

const ChuckSize = 5

var pointUpdateBuffers []*interfaces.PointUpdateBuffer

func CreatePointDeepCopy(point model.Point) model.Point {
	var outputPoint model.Point
	out, _ := json.Marshal(point)
	_ = json.Unmarshal(out, &outputPoint)
	return outputPoint
}

func GetPoint(uuid string, args api.Args) *model.Point {
	for _, pub := range pointUpdateBuffers {
		if pub.UUID == uuid && pub.State != interfaces.Initialize {
			point := CreatePointDeepCopy(*pub.Point)
			if !args.WithPriority {
				point.Priority = nil
			}
			if !args.WithTags {
				point.Tags = nil
			}
			if !args.WithMetaTags {
				point.MetaTags = nil
			}
			return &point
		}
	}
	return nil
}

// updatePriority it updates priority array of point model
// it attaches the point model fields values for updating it on its parent function
func (d *GormDatabase) updatePriority(pointModel *model.Point, priority *map[string]*float64, writeOnDB bool) (
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
		if writeOnDB {
			d.DB.Model(&model.Priority{}).Where("point_uuid = ?", pointModel.UUID).Updates(&pointModel.Priority)
		}
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
		priorityMapToPatch := d.priorityMapToPatch(priorityMap)
		if writeOnDB {
			d.DB.Model(&pointModel.Priority).Where("point_uuid = ?", pointModel.UUID).Updates(&priorityMapToPatch)
		}
	}
	if !presentValueFromPriority {
		// presentValue will be OriginalValue if PointPriorityArrayMode is PriorityArrayToWriteValue or
		// ReadOnlyNoPriorityArrayRequired
		presentValue = pointModel.OriginalValue
	}
	return pointModel, priorityMap, presentValue, writeValue, isPriorityChanged
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

func (d *GormDatabase) bufferPointUpdate(uuid string, body *model.Point, point *model.Point) {
	d.pointBuffersMutex.Lock()
	defer d.pointBuffersMutex.Unlock()
	pointUpdateBuffer := &interfaces.PointUpdateBuffer{
		UUID:  uuid,
		Body:  body,
		Point: point,
		State: interfaces.Created,
	}
	for index, pub := range pointUpdateBuffers {
		if pub.UUID == uuid {
			pointUpdateBuffers[index] = pointUpdateBuffer
			return
		}
	}
	pointUpdateBuffers = append(pointUpdateBuffers, pointUpdateBuffer)
}

func (d *GormDatabase) bufferPointUpdateBody(uuid string, body *model.Point) {
	pointUpdateBuffer := &interfaces.PointUpdateBuffer{
		UUID:  uuid,
		Body:  body,
		Point: nil,
		State: interfaces.Initialize,
	}
	for index, pub := range pointUpdateBuffers {
		if pub.UUID == uuid {
			pointUpdateBuffers[index].Body = body
			return
		}
	}
	pointUpdateBuffers = append(pointUpdateBuffers, pointUpdateBuffer)
}

func (d *GormDatabase) FlushPointUpdateBuffers() {
	log.Debug("Flush point update buffers has is been called...")
	if len(pointUpdateBuffers) == 0 {
		log.Debug("Point update buffers not found")
		return
	}

	d.pointBuffersMutex.Lock()
	var tempPointUpdateBuffers []*interfaces.PointUpdateBuffer
	tempPointUpdateBuffers = append(tempPointUpdateBuffers, pointUpdateBuffers...)
	for _, pub := range pointUpdateBuffers {
		if pub.State != interfaces.Initialize {
			pub.State = interfaces.Updating
		}
	}
	chuckPointUpdateBuffers := ChuckPointUpdateBuffer(tempPointUpdateBuffers, ChuckSize)
	d.pointBuffersMutex.Unlock()
	for _, chuckPointUpdateBuffer := range chuckPointUpdateBuffers {
		wg := &sync.WaitGroup{}
		for _, point := range chuckPointUpdateBuffer {
			if point.State != interfaces.Initialize {
				wg.Add(1)
				go func(point *interfaces.PointUpdateBuffer) {
					defer wg.Done()
					_, _ = d.UpdatePoint(point.UUID, point.Body, false)
					d.removeUpdatePointBuffer(point.UUID)
				}(point)
				time.Sleep(200 * time.Millisecond) // for don't let them call at once
			}
		}
		wg.Wait()
	}
	log.Debug("Finished flush point update buffers process")
}

func ChuckPointUpdateBuffer(array []*interfaces.PointUpdateBuffer, chunkSize int) [][]*interfaces.PointUpdateBuffer {
	var chucks [][]*interfaces.PointUpdateBuffer
	for i := 0; i < len(array); i += chunkSize {
		end := i + chunkSize
		if end > len(array) {
			end = len(array)
		}
		chucks = append(chucks, array[i:end])
	}
	return chucks
}

func (d *GormDatabase) updateUpdatePointBufferPointWriter(uuid string, priority *map[string]*float64) {
	d.pointBuffersMutex.Lock()
	defer d.pointBuffersMutex.Unlock()
	for _, pub := range pointUpdateBuffers {
		if pub.UUID == uuid {
			_, _ = priorityarray.ApplyMapToPriorityArray(pub.Body, priority)
			return
		}
	}
}

func (d *GormDatabase) updateUpdatePointBufferPoint(point *model.Point) {
	d.pointBuffersMutex.Lock()
	defer d.pointBuffersMutex.Unlock()
	for _, pub := range pointUpdateBuffers {
		if pub.UUID == point.UUID {
			// keep tags and meta_tags as is, this will be called from point write where these attributes don't exist
			point.MetaTags = pub.Point.MetaTags
			point.Tags = pub.Point.Tags
			pub.Point = point
			return
		}
	}
}

// There is also a chance of next update request when it's updating, so for such case don't remove that record
// For eg: when it's updating the record, update new request been called, so don't remove those recently created record
func (d *GormDatabase) removeUpdatePointBuffer(uuid string) {
	d.pointBuffersMutex.Lock()
	defer d.pointBuffersMutex.Unlock()
	for i, pub := range pointUpdateBuffers {
		if pub.UUID == uuid {
			if pub.State == interfaces.Updating {
				pointUpdateBuffers = d.removeIndex(i)
			}
			return
		}
	}
}

func (d *GormDatabase) removeIndex(index int) []*interfaces.PointUpdateBuffer {
	return append(pointUpdateBuffers[:index], pointUpdateBuffers[index+1:]...)
}

func (d *GormDatabase) DeletePointByName(networkName, deviceName, pointName string, args api.Args) (bool, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.Joins("JOIN devices ON points.device_uuid = devices.uuid").
		Joins("JOIN networks ON devices.network_uuid = networks.uuid").
		Where("networks.name = ?", networkName).Where("devices.name = ?", deviceName).
		Where("points.name = ?", pointName).
		First(&pointModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(pointModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) getUpdatePointBufferBody(uuid string) *model.Point {
	for _, pub := range pointUpdateBuffers {
		if pub.UUID == uuid {
			return pub.Body
		}
	}
	return nil
}
