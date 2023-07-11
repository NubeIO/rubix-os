package database

import (
	"errors"
	"fmt"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/config"
	"github.com/NubeIO/rubix-os/interfaces"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/float"
	"github.com/NubeIO/rubix-os/utils/integer"
	"github.com/NubeIO/rubix-os/utils/nmath"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"github.com/NubeIO/rubix-os/utils/priorityarray"
)

func (d *GormDatabase) GetPoints(args argspkg.Args) ([]*model.Point, error) {
	var pointsModel []*model.Point
	query := d.buildPointQuery(args)
	if err := query.Find(&pointsModel).Error; err != nil {
		return nil, err
	}
	return pointsModel, nil
}

func (d *GormDatabase) GetPointsBulkUUIs() ([]string, error) {
	var uuids []string
	if err := d.DB.Model(&model.Point{}).Select("uuid").Find(&uuids).Error; err != nil {
		return nil, err
	}
	return uuids, nil
}

func (d *GormDatabase) GetPointsBulk(bulkPoints []*model.Point) ([]*model.Point, error) {
	var pointsModel []*model.Point
	points, err := d.GetPoints(argspkg.Args{WithPriority: true})
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

func (d *GormDatabase) GetPoint(uuid string, args argspkg.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&pointModel).Error; err != nil {
		return nil, err
	}
	return pointModel, nil
}

func (d *GormDatabase) GetPointByName(networkName, deviceName, pointName string, args argspkg.Args) (*model.Point, error) {
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

func (d *GormDatabase) GetPointByUUID(networkUUID, deviceUUID, pointUUID string, args argspkg.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.Joins("JOIN devices ON points.device_uuid = devices.uuid").
		Joins("JOIN networks ON devices.network_uuid = networks.uuid").
		Where("networks.uuid = ?", networkUUID).Where("devices.uuid = ?", deviceUUID).
		Where("points.uuid = ?", pointUUID).
		First(&pointModel).Error; err != nil {
		return nil, err
	}
	return pointModel, nil
}

func (d *GormDatabase) GetOnePointByArgs(args argspkg.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.First(&pointModel).Error; err != nil {
		return nil, err
	}
	return pointModel, nil
}

func (d *GormDatabase) CreatePoint(body *model.Point) (*model.Point, error) {
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	var device *model.Device
	query := d.DB.Where("uuid = ? ", body.DeviceUUID).First(&device)
	if query.Error != nil {
		return nil, fmt.Errorf("no such parent device with uuid %s", body.DeviceUUID)
	}
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Point)
	body.Name = name
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
	if body.HistoryEnable == nil {
		body.HistoryEnable = boolean.NewTrue()
	}
	body.HistoryCOVThreshold = defaultHistoryCOVThreshold(body.HistoryCOVThreshold)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	go d.PublishPointsList("")
	return body, nil
}

func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point) (*model.Point, error) {
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	body.Name = name
	pointModel, err := d.GetPoint(uuid, argspkg.Args{WithPriority: true})
	if err != nil {
		return nil, err
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
	if err = d.updateTags(&pointModel, body.Tags); err != nil {
		return nil, err
	}
	publishPointList := body.Name != pointModel.Name
	if err = d.DB.Model(&pointModel).Select("*").Updates(&body).Error; err != nil {
		return nil, err
	}

	if body.Priority != nil {
		pointModel.Priority = body.Priority
	}
	priorityMap := priorityarray.ConvertToMap(*pointModel.Priority)
	pointWriter := &model.PointWriter{
		Priority: &priorityMap,
	}
	pnt, _, _, _, err := d.updatePointValue(pointModel, pointWriter)
	if publishPointList {
		go d.PublishPointsList("")
	}
	return pnt, err
}

func (d *GormDatabase) PointWrite(uuid string, body *model.PointWriter) (
	returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, false, false, false, query.Error
	}
	if body == nil || body.Priority == nil {
		return nil, false, false, false, errors.New("no priority value is been sent")
	} else {
		pointModel.ValueUpdatedFlag = boolean.NewTrue()
	}
	point, isPresentValueChange, isWriteValueChange, isPriorityChanged, err :=
		d.updatePointValue(pointModel, body)
	return point, isPresentValueChange, isWriteValueChange, isPriorityChanged, err
}

func (d *GormDatabase) updatePointValue(pointModel *model.Point, pointWriter *model.PointWriter) (
	returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	priority := pointWriter.Priority
	if pointModel.PointPriorityArrayMode == "" {
		pointModel.PointPriorityArrayMode = model.PriorityArrayToPresentValue // sets default priority array mode
	}

	pointModel, priority, presentValue, writeValue, isPriorityChanged := d.updatePriority(pointModel, priority)
	presentValueTransformFault := false
	if pointWriter.PresentValue == nil {
		ov := float.Copy(presentValue)
		pointModel.OriginalValue = ov
		wv := float.Copy(writeValue)
		pointModel.WriteValueOriginal = wv

		transform := PointValueTransformOnRead(presentValue, pointModel.ScaleEnable, pointModel.MultiplicationFactor,
			pointModel.ScaleInMin, pointModel.ScaleInMax, pointModel.ScaleOutMin, pointModel.ScaleOutMax, pointModel.Offset)
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

		if !integer.IsUnit32Nil(pointModel.Decimal) && presentValue != nil {
			value := nmath.RoundTo(*presentValue, *pointModel.Decimal)
			presentValue = &value
		}
	} else {
		presentValue = pointWriter.PresentValue
	}

	isPresentValueChange = !float.ComparePtrValues(pointModel.PresentValue, presentValue) // Use for createCOVHistory
	// isWriteValueChange
	// in some cases we don't write presentValue directly to the model.
	// Instead of writing that, we firstly write it on writeValue & and we read that writeValue to update presentValue.
	// Examples are: modbus, edge28 plugins
	// So for such cases, to trigger that value we do this comparison
	isWriteValueChange = !float.ComparePtrValues(pointModel.WriteValue, writeValue)
	cov, _ := nmath.Cov(float.NonNil(presentValue), float.NonNil(pointModel.LastHistoryValue),
		*defaultHistoryCOVThreshold(pointModel.HistoryCOVThreshold))
	// If the present value transformations have resulted in an error, DB needs to be updated with the errors,
	// but PresentValue should not change
	if !presentValueTransformFault {
		pointModel.PresentValue = presentValue
	}
	pointModel.WriteValue = writeValue

	if cov && boolean.IsTrue(config.Get().PointHistory.Enable) && checkHistoryCovType(string(pointModel.HistoryType)) &&
		d.CheckHistoryEnable(pointModel.UUID) {
		pointHistory := &model.PointHistory{
			PointUUID: pointModel.UUID,
			Value:     pointModel.PresentValue,
			Timestamp: time.Now().UTC(),
		}
		_, err = d.CreatePointHistory(pointHistory)
		if err != nil {
			log.Errorf("point: issue on write history for point: %v\n", err)
			return nil, false, false, false, err
		}
		pointModel.LastHistoryValue = pointModel.PresentValue
	}

	_ = d.DB.Model(&pointModel).Select("*").Updates(&pointModel)

	if isPresentValueChange {
		err = d.PublishPointCov(pointModel.UUID)
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

// UpdatePointSuccess will only update the CommonFault properties of the point, all other properties will not be updated.
// Does not update `LastFail`.
func (d *GormDatabase) UpdatePointSuccess(uuid string, body *model.Point) error {
	return d.DB.Model(&body).
		Where("uuid = ?", uuid).
		Select("InFault", "MessageLevel", "MessageCode", "Message", "LastOk", "InSync").
		Updates(&body).
		Error
}

func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Point{})
	go d.PublishPointsList("")
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) PointWriteByName(networkName, deviceName, pointName string, body *model.PointWriter) (*model.Point, error) {
	point, err := d.GetPointByName(networkName, deviceName, pointName, argspkg.Args{})
	if err != nil {
		return nil, err
	}
	point, err = d.WritePointPlugin(point.UUID, body)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (d *GormDatabase) DeleteOnePointByArgs(args argspkg.Args) (bool, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.First(&pointModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&pointModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeletePointByName(networkName, deviceName, pointName string, args argspkg.Args) (bool, error) {
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

func (d *GormDatabase) GetPointWithParent(uuid string) (*interfaces.PointWithParent, error) {
	var pointWithParent *interfaces.PointWithParent
	if err := d.DB.Table("points").
		Select("points.uuid, points.name, devices.uuid AS device_uuid, devices.name AS device_name, "+
			"networks.uuid AS network_names, networks.name AS network_name").
		Joins("JOIN devices ON points.device_uuid = devices.uuid").
		Joins("JOIN networks ON devices.network_uuid = networks.uuid").
		Where("points.uuid = ?", uuid).
		First(&pointWithParent).Error; err != nil {
		return nil, err
	}
	return pointWithParent, nil
}

func (d *GormDatabase) GetPointsForCreateInterval() ([]*interfaces.PointHistoryInterval, error) {
	var pointIntervalHistory []*interfaces.PointHistoryInterval
	if err := d.DB.Table("points").
		Select("points.uuid, points.history_interval, point_histories.timestamp AS timestamp, "+
			"points.present_value").
		Joins("JOIN devices ON points.device_uuid = devices.uuid").
		Joins("JOIN networks ON devices.network_uuid = networks.uuid").
		Joins("LEFT JOIN (SELECT point_uuid, MAX(timestamp) AS timestamp "+
			"FROM point_histories "+
			"GROUP BY point_uuid) point_histories "+
			"ON points.uuid = point_histories.point_uuid").
		Where("points.history_enable IS NOT FALSE").
		Where("devices.history_enable IS NOT FALSE").
		Where("networks.history_enable IS TRUE").
		Where("points.history_interval > ?", 0).
		Where(fmt.Sprintf("points.history_type != '%s'", model.HistoryTypeCov)). // points.history_type != "COV" is invalid
		Scan(&pointIntervalHistory).Error; err != nil {
		return nil, err
	}
	return pointIntervalHistory, nil
}

func (d *GormDatabase) CheckHistoryEnable(uuid string) bool {
	var count int64
	if err := d.DB.Table("points").
		Select("points.uuid").
		Joins("JOIN devices ON points.device_uuid = devices.uuid").
		Joins("JOIN networks ON devices.network_uuid = networks.uuid").
		Where("points.uuid = ?", uuid).
		Where("points.history_enable IS NOT FALSE").
		Where("devices.history_enable IS NOT FALSE").
		Where("networks.history_enable IS TRUE").
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
