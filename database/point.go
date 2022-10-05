package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nmath"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) GetPoints(args api.Args) ([]*model.Point, error) {
	var pointsModel []*model.Point
	query := d.buildPointQuery(args)
	if err := query.Find(&pointsModel).Error; err != nil {
		return nil, err
	}
	return pointsModel, nil
}

func (d *GormDatabase) GetPointsBulk(bulkPoints []*model.Point) ([]*model.Point, error) {
	var pointsModel []*model.Point
	points, err := d.GetPoints(api.Args{WithPriority: true})
	if err != nil {
		return nil, err
	}
	for _, pnt := range points {
		for _, search := range bulkPoints {
			if pnt.UUID == search.UUID {
				pointsModel = append(pointsModel, pnt)
			}
		}
	}
	return pointsModel, nil
}

func (d *GormDatabase) GetPoint(uuid string, args api.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&pointModel).Error; err != nil {
		return nil, err
	}
	return pointModel, nil
}

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

func (d *GormDatabase) GetOnePointByArgs(args api.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.First(&pointModel).Error; err != nil {
		return nil, err
	}
	return pointModel, nil
}

func (d *GormDatabase) CreatePoint(body *model.Point, fromPlugin bool) (*model.Point, error) {
	network, err := d.GetNetworkByDeviceUUID(body.DeviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Point)
	if body.Decimal == nil {
		body.Decimal = nils.NewUint32(2)
	}
	obj, err := checkObjectType(body.ObjectType)
	if err != nil {
		return nil, err
	}
	body.ObjectType = string(obj)
	if body.PointPriorityArrayMode == "" {
		body.PointPriorityArrayMode = model.PriorityArrayToPresentValue // sets default priority array mode.
	}
	body.ThingClass = model.ThingClass.Point
	body.InSync = boolean.NewFalse()
	if body.Priority == nil {
		body.Priority = &model.Priority{}
	}
	if body.ScaleEnable == nil {
		body.ScaleEnable = boolean.NewFalse()
	}
	if body.ScaleInMin == nil {
		body.ScaleInMin = float.New(0)
	}
	if body.ScaleInMax == nil {
		body.ScaleInMax = float.New(0)
	}
	if body.ScaleOutMin == nil {
		body.ScaleOutMin = float.New(0)
	}
	if body.ScaleOutMax == nil {
		body.ScaleOutMax = float.New(0)
	}
	if body.Offset == nil {
		body.Offset = float.New(0)
	}
	if body.PollRate == "" {
		body.PollRate = model.RATE_NORMAL
	}
	if body.PollPriority == "" {
		body.PollPriority = model.PRIORITY_NORMAL
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}

	// check for mapping
	if network.AutoMappingNetworksSelection != "" {
		if network.AutoMappingNetworksSelection != "disable" {
			pointMapping := &model.PointMapping{}
			pointMapping.Point = body
			pointMapping.AutoMappingFlowNetworkName = network.AutoMappingFlowNetworkName
			pointMapping.AutoMappingFlowNetworkUUID = network.AutoMappingFlowNetworkUUID
			pointMapping.AutoMappingNetworksSelection = []string{network.AutoMappingNetworksSelection}
			pointMapping.AutoMappingEnableHistories = network.AutoMappingEnableHistories
			pointMapping, err = d.CreatePointMapping(pointMapping)
			if err != nil {
				log.Errorln("points.db.CreatePoint() failed to make auto point mapping")
				return nil, err
			} else {
				log.Println("points.db.CreatePoint() added point new mapping")
			}
		}
	}
	return body, nil
}

func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point, fromPlugin bool, afterRealDeviceUpdate bool) (
	*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	existingName, existingAddrID := d.pointNameExists(body)
	if existingName {
		eMsg := fmt.Sprintf("a point with existing name: %s exists", body.Name)
		return nil, errors.New(eMsg)
	}

	if !integer.IsNil(body.AddressID) {
		if existingAddrID {
			eMsg := fmt.Sprintf("a point with existing AddressID: %d exists", integer.NonNil(body.AddressID))
			return nil, errors.New(eMsg)
		}
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&pointModel, body.Tags); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&pointModel).Select("*").Updates(&body).Error; err != nil {
		return nil, err
	}

	// TODO: we need to decide if a read only point needs to have a priority array or if it should just be nil.
	if body.Priority == nil {
		pointModel.Priority = &model.Priority{}
	} else {
		pointModel.Priority = body.Priority
	}

	priorityMap := priorityarray.ConvertToMap(*pointModel.Priority)
	pnt, _, _, _, err := d.updatePointValue(pointModel, &priorityMap, fromPlugin, afterRealDeviceUpdate, nil)
	return pnt, err
}

func (d *GormDatabase) PointWrite(uuid string, body *model.PointWriter, fromPlugin bool, afterRealDeviceUpdate bool,
	currentWriterUUID *string) (returnPoint *model.Point, isPresentValueChange, isWriteValueChange,
	isPriorityChanged bool, err error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, false, false, false, query.Error
	}
	if body == nil || body.Priority == nil {
		return nil, false, false, false,
			errors.New("no priority value is been sent")
	} else {
		pointModel.ValueUpdatedFlag = boolean.NewTrue()
	}
	point, isPresentValueChange, isWriteValueChange, isPriorityChanged, err :=
		d.updatePointValue(pointModel, body.Priority, fromPlugin, afterRealDeviceUpdate, currentWriterUUID)
	return point, isPresentValueChange, isWriteValueChange, isPriorityChanged, err
}

func (d *GormDatabase) updatePointValue(pointModel *model.Point, priority *map[string]*float64, fromPlugin bool, afterRealDeviceUpdate bool, currentWriterUUID *string) (returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	if pointModel.PointPriorityArrayMode == "" {
		pointModel.PointPriorityArrayMode = model.PriorityArrayToPresentValue // sets default priority array mode
	}

	pointModel, priority, presentValue, writeValue, isPriorityChanged := d.updatePriority(pointModel, priority)
	ov := float.Copy(presentValue)
	pointModel.OriginalValue = ov
	wv := float.Copy(writeValue)
	pointModel.WriteValueOriginal = wv

	presentValueTransformFault := false
	transform := PointValueTransformOnRead(presentValue, pointModel.ScaleEnable, pointModel.MultiplicationFactor,
		pointModel.ScaleInMin, pointModel.ScaleInMax, pointModel.ScaleOutMin, pointModel.ScaleOutMax, pointModel.Offset)
	if afterRealDeviceUpdate {
		pointModel.CommonFault.InFault = false
		pointModel.CommonFault.MessageLevel = model.MessageLevel.Info
		pointModel.CommonFault.MessageCode = model.CommonFaultCode.PointWriteOk
		pointModel.CommonFault.Message = fmt.Sprintf("last-updated: %s", utilstime.TimeStamp())
		pointModel.CommonFault.LastOk = time.Now().UTC()
	}
	presentValue = transform
	val, err := pointUnits(presentValue, pointModel.Unit, pointModel.UnitTo)
	if err != nil {
		pointModel.CommonFault.InFault = true
		pointModel.CommonFault.MessageLevel = model.MessageLevel.Warning
		pointModel.CommonFault.MessageCode = model.CommonFaultCode.PointError
		pointModel.CommonFault.Message = fmt.Sprint("point.db updatePointValue() invalid point units. error:", err)
		pointModel.CommonFault.LastFail = time.Now().UTC()
		presentValueTransformFault = true
	} else {
		presentValue = val
	}
	writeValue = PointValueTransformOnWrite(writeValue, pointModel.ScaleEnable, pointModel.MultiplicationFactor, pointModel.ScaleInMin, pointModel.ScaleInMax, pointModel.ScaleOutMin, pointModel.ScaleOutMax, pointModel.Offset)

	// example for wires and modbus:
	// if a new value is written from wires then set this to false so the modbus knows on the next poll to write a new
	// value to the modbus point
	if fromPlugin && afterRealDeviceUpdate {
		pointModel.InSync = boolean.NewTrue() // TODO: do we still use InSync?
		// pointModel.WritePollRequired = boolean.NewFalse()  // WritePollRequired should be set by the plugins (they know best)
	} else {
		pointModel.InSync = boolean.NewFalse() // TODO: do we still use InSync?
		// pointModel.WritePollRequired = boolean.NewTrue()  // WritePollRequired should be set by the plugins (they know best)
	}

	if !integer.IsUnit32Nil(pointModel.Decimal) && presentValue != nil {
		value := nmath.RoundTo(*presentValue, *pointModel.Decimal)
		presentValue = &value
	}

	isPresentValueChange = !float.ComparePtrValues(pointModel.PresentValue, presentValue) // Use for createCOVHistory
	// isWriteValueChange
	// in some cases we don't write presentValue directly to the model.
	// Instead of writing that, we firstly write it on writeValue & and we read that writeValue to update presentValue.
	// Examples are: modbus, edge28 plugins
	// So for such cases, to trigger that value we do this comparison
	isWriteValueChange = !float.ComparePtrValues(pointModel.WriteValue, writeValue)
	isChange := isPresentValueChange || isWriteValueChange || isPriorityChanged

	// If the present value transformations have resulted in an error, DB needs to be updated with the errors,
	// but PresentValue should not change
	if !presentValueTransformFault {
		pointModel.PresentValue = presentValue
	}
	pointModel.WriteValue = writeValue

	_ = d.DB.Model(&pointModel).Select("*").Updates(&pointModel)
	if isChange {
		err = d.ProducersPointWrite(pointModel.UUID, priority, pointModel.PresentValue, isPresentValueChange,
			currentWriterUUID)
		if err != nil {
			return nil, false, false, false, err
		}
		d.DB.Model(&model.Writer{}).
			Where("writer_thing_uuid = ?", pointModel.UUID).
			Update("present_value", pointModel.PresentValue)

		// TODO: this section is added temporarily to support edge->influx plugin while we sort the official histories
		if isPresentValueChange {
			cli := client.NewLocalClient()
			_, err = cli.WritePointPlugin(pointModel.UUID, &model.PointWriter{}, "edgeinflux")
		}
	}
	return pointModel, isPresentValueChange, isWriteValueChange, isPriorityChanged, nil
}

// UpdatePointErrors will only update the CommonFault properties of the point, all other properties will not be updated.
// Does not update `LastOk`.
func (d *GormDatabase) UpdatePointErrors(uuid string, body *model.Point) error {
	return d.DB.Model(&body).
		Where("uuid = ?", uuid).
		Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail", "InSync").
		Updates(&body).
		Error
}

func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	point, err := d.GetPoint(uuid, api.Args{})
	if err != nil {
		return false, err
	}
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &point.UUID})
	if producers != nil {
		for _, producer := range producers {
			_, _ = d.DeleteProducer(producer.UUID)
		}
	}
	writers, _ := d.GetWriters(api.Args{WriterThingUUID: &point.UUID})
	if writers != nil {
		for _, writer := range writers {
			_, _ = d.DeleteWriter(writer.UUID)
		}
	}
	query := d.DB.Delete(&point)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	return r != 0, nil
}

func (d *GormDatabase) PointWriteByName(networkName, deviceName, pointName string, body *model.PointWriter,
	fromPlugin bool) (*model.Point, error) {
	point, err := d.GetPointByName(networkName, deviceName, pointName)
	if err != nil {
		return nil, err
	}
	write, _, _, _, err := d.PointWrite(point.UUID, body, fromPlugin, false, nil)
	if err != nil {
		return nil, err
	}
	return write, nil
}
