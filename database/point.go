package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/src/client"
	"github.com/NubeIO/flow-framework/urls"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"sync"
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

func (d *GormDatabase) GetPointByName(networkName, deviceName, pointName string, args api.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.Joins("JOIN devices ON points.device_uuid = devices.uuid").
		Joins("JOIN networks ON devices.network_uuid = networks.uuid").
		Where("networks.name = ?", networkName).Where("devices.name = ?", deviceName).
		Where("points.name = ?", pointName).
		First(&pointModel).Error; err != nil {
		return nil, err
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

	go d.PublishPointsList("")
	err = d.CreatePointAutoMapping(body)
	if err != nil {
		log.Errorln("points.db.CreatePointAutoMapping() failed to make auto mapping")
		return nil, err
	} else {
		log.Println("points.db.CreatePointAutoMapping() added point new mapping")
	}
	return body, nil
}

func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point, fromPlugin bool, afterRealDeviceUpdate bool) (
	*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Tags").Preload("MetaTags").
		Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	existingName, existingAddrID := d.pointNameExists(body)
	if existingAddrID && boolean.IsTrue(body.IsBitwise) && body.BitwiseIndex != nil && *body.BitwiseIndex >= 0 {
		existingAddrID = false
	}
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
	publishPointList := body.Name != pointModel.Name
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
	pnt, _, _, _, err := d.updatePointValue(pointModel, &priorityMap, fromPlugin, afterRealDeviceUpdate, nil, false)
	if publishPointList {
		go d.PublishPointsList("")
	}
	d.UpdateProducerByProducerThingUUID(pointModel.UUID, pointModel.Name, pointModel.HistoryEnable,
		pointModel.HistoryType, pointModel.HistoryInterval)
	err = d.UpdatePointAutoMapping(pointModel)
	if err != nil {
		log.Errorln("points.db.UpdatePointAutoMapping() failed to make auto mapping")
		return nil, err
	} else {
		log.Println("points.db.UpdatePointAutoMapping() added point new mapping")
	}
	return pnt, err
}

func (d *GormDatabase) PointWrite(uuid string, body *model.PointWriter, fromPlugin bool, afterRealDeviceUpdate bool,
	currentWriterUUID *string, forceWrite bool) (returnPoint *model.Point, isPresentValueChange, isWriteValueChange,
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
		d.updatePointValue(pointModel, body.Priority, fromPlugin, afterRealDeviceUpdate, currentWriterUUID, forceWrite)
	return point, isPresentValueChange, isWriteValueChange, isPriorityChanged, err
}

func (d *GormDatabase) updatePointValue(pointModel *model.Point, priority *map[string]*float64, fromPlugin bool,
	afterRealDeviceUpdate bool, currentWriterUUID *string, forceWrite bool) (returnPoint *model.Point,
	isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
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
	isChange := isPresentValueChange || isWriteValueChange || isPriorityChanged || forceWrite

	// If the present value transformations have resulted in an error, DB needs to be updated with the errors,
	// but PresentValue should not change
	if !presentValueTransformFault {
		pointModel.PresentValue = presentValue
	}
	pointModel.WriteValue = writeValue
	// last update was ok
	pointModel.MessageLevel = model.MessageLevel.Info
	pointModel.MessageCode = model.CommonFaultCode.Ok
	pointModel.Message = fmt.Sprintf("lastMessage: %s", utilstime.TimeStamp())
	pointModel.LastOk = time.Now()
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

		if isPresentValueChange {
			err = d.PublishPointCov(pointModel.UUID)
		}
	}
	return pointModel, isPresentValueChange, isWriteValueChange, isPriorityChanged, nil
}

// UpdatePointErrors will only update the CommonFault properties of the point, all other properties will not be updated.
// Does not update `LastOk`.
func (d *GormDatabase) UpdatePointErrors(uuid string, body *model.Point) error {
	return d.DB.Model(&body).
		Where("uuid = ?", uuid).
		Select("InFault", "MessageLevel", "MessageCode", "Message", "LastFail", "InSync", "Connection").
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
		var wg sync.WaitGroup
		for _, producer := range producers {
			wg.Add(1)
			producer := producer
			go func() {
				defer wg.Done()
				_, _ = d.DeleteProducer(producer.UUID)
			}()
		}
		wg.Wait()
	}
	writers, _ := d.GetWriters(api.Args{WriterThingUUID: &point.UUID})
	if writers != nil {
		var wg sync.WaitGroup
		for _, writer := range writers {
			wg.Add(1)
			writer := writer
			go func() {
				defer wg.Done()
				_, _ = d.DeleteWriter(writer.UUID)
			}()
		}
		wg.Wait()
	}
	query := d.DB.Delete(&point)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	go d.PublishPointsList("")
	var aType = api.ArgsType
	deviceModel, err := d.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil {
		return false, err
	}
	if boolean.IsTrue(deviceModel.AutoMappingEnable) {
		fn, err := d.selectFlowNetwork(deviceModel.AutoMappingFlowNetworkName, deviceModel.AutoMappingFlowNetworkUUID)
		if err != nil {
			return false, err
		}
		cli := client.NewFlowClientCliFromFN(fn)
		url := urls.SingularUrlByArg(urls.PointUrl, aType.AutoMappingUUID, point.UUID)
		_ = cli.DeleteQuery(url)
	}
	return r != 0, nil
}

func (d *GormDatabase) PointWriteByName(networkName, deviceName, pointName string, body *model.PointWriter,
	fromPlugin bool) (*model.Point, error) {
	point, err := d.GetPointByName(networkName, deviceName, pointName, api.Args{})
	if err != nil {
		return nil, err
	}
	write, _, _, _, err := d.PointWrite(point.UUID, body, fromPlugin, false, nil, false)
	if err != nil {
		return nil, err
	}
	return write, nil
}

func (d *GormDatabase) DeleteOnePointByArgs(args api.Args) (bool, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.First(&pointModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&pointModel)
	return d.deleteResponseBuilder(query)
}
