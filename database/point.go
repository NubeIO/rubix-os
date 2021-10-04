package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"reflect"
	"time"
)

func (d *GormDatabase) GetPoints(args api.Args) ([]*model.Point, error) {
	var pointsModel []*model.Point
	query := d.buildPointQuery(args)
	if err := query.Find(&pointsModel).Error; err != nil {
		return nil, err
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

// PointDeviceByAddressID will query by device_uuid = ? AND object_type = ? AND address_id = ?
func (d *GormDatabase) PointDeviceByAddressID(pointUUID string, body *model.Point) (*model.Point, bool) {
	var pointModel *model.Point
	deviceUUID := body.DeviceUUID
	objType := body.ObjectType
	addressID := body.AddressId
	f := fmt.Sprintf("device_uuid = ? AND object_type = ? AND address_id = ?")
	query := d.DB.Where(f, deviceUUID, objType, addressID, pointUUID).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, false
	}
	return pointModel, true
}

// CreatePoint creates an object.
func (d *GormDatabase) CreatePoint(body *model.Point, streamUUID string) (*model.Point, error) {
	var deviceModel *model.Device
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Point)
	deviceUUID := body.DeviceUUID
	body.Name = nameIsNil(body.Name)
	obj, err := checkObjectType(body.ObjectType)
	if err != nil {
		return nil, err
	}
	body.ObjectType = obj
	query := d.DB.Where("uuid = ? ", deviceUUID).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Description == "" {
		body.Description = "na"
	}
	body.ThingClass = model.ThingClass.Point
	body.CommonEnable.Enable = utils.NewTrue()
	body.CommonFault.InFault = true
	body.CommonFault.MessageLevel = model.MessageLevel.NoneCritical
	body.CommonFault.MessageCode = model.CommonFaultCode.PluginNotEnabled
	body.CommonFault.Message = model.CommonFaultMessage.PluginNotEnabled
	body.CommonFault.LastFail = time.Now().UTC()
	body.CommonFault.LastOk = time.Now().UTC()
	body.InSync = utils.NewFalse()
	if body.Priority == nil {
		body.Priority = &model.Priority{}
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, query.Error
	}
	if streamUUID != "" {
		producerModel := new(model.Producer)
		producerModel.StreamUUID = streamUUID
		producerModel.ProducerThingUUID = body.UUID
		producerModel.ProducerThingClass = model.ThingClass.Point
		producerModel.ProducerThingType = model.ThingClass.Point
		_, err := d.CreateProducer(producerModel)
		if err != nil {
			return nil, errors.New("ERROR on create new producer to an existing stream")
		}
	}
	plug, err := d.GetPluginIDFromDevice(deviceUUID)
	if err != nil {
		return nil, errors.New("ERROR failed to get plugin uuid")
	}
	t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, plug.PluginConfId, body.UUID)
	d.Bus.RegisterTopic(t)
	err = d.Bus.Emit(eventbus.CTX(), t, body)
	if err != nil {
		return nil, errors.New("ERROR on device eventbus")
	}
	return body, query.Error
}

// UpdatePoint returns the device for the given id or nil.
func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").Find(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&pointModel, body.Tags); err != nil {
			return nil, err
		}
	}
	body.InSync = utils.NewFalse()
	query = d.DB.Model(&pointModel).Updates(&body)
	pnt, err := d.UpdatePointValue(uuid, pointModel, false)
	if err != nil {
		return nil, err
	}
	if !fromPlugin { //stop looping
		plug, err := d.GetPluginIDFromDevice(pointModel.DeviceUUID)
		if err != nil {
			return nil, errors.New("ERROR failed to get plugin uuid")
		}
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, plug.PluginConfId, pointModel.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, pointModel)
		if err != nil {
			return nil, errors.New("ERROR on device eventbus")
		}
	}
	return pnt, nil
}

// PointWrite returns the device for the given id or nil.
func (d *GormDatabase) PointWrite(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").Find(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Priority != nil {
		priority := map[string]interface{}{}
		priorityValue := reflect.ValueOf(*body.Priority)
		typeOfPriority := priorityValue.Type()
		highestPri := utils.NewArray()
		highestValue := utils.NewMap()
		for i := 0; i < priorityValue.NumField(); i++ {
			if priorityValue.Field(i).Type().Kind().String() == "ptr" {
				val := priorityValue.Field(i).Interface().(*float64)
				if val == nil {
					priority[typeOfPriority.Field(i).Name] = nil
				} else {
					highestPri.Add(i)
					highestValue.Set(i, *val)
					priority[typeOfPriority.Field(i).Name] = *val
				}
			}
		}
		notNil := false
		for _, v := range priority { //check if there is a value in priority array
			if v != nil {
				notNil = true
			}
		}
		if notNil {
			min, _ := highestPri.MinMaxInt() //get the highest priority
			val := highestValue.Get(min)     //get the highest priority value
			body.CurrentPriority = &min      //TODO check conversion
			v := val.(float64)
			body.PresentValue = &v //process the units as in temperature conversion
		}
		d.DB.Model(&pointModel.Priority).Updates(&priority)
	}
	if !fromPlugin { //stop looping
		plug, err := d.GetPluginIDFromDevice(pointModel.DeviceUUID)
		if err != nil {
			return nil, errors.New("ERROR failed to get plugin uuid")
		}
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, plug.PluginConfId, pointModel.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, pointModel)
		if err != nil {
			return nil, errors.New("ERROR on device eventbus")
		}
	}
	return pointModel, nil
}

// UpdatePointValue returns the device for the given id or nil.
func (d *GormDatabase) UpdatePointValue(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	//TODO add in a check to make sure user doesn't set the addressID and the ObjectType the same as another point
	//check if there is an existing device with this address code
	_, err := checkObjectType(body.ObjectType)
	if err != nil {
		return nil, err
	}
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").Find(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	var presentValue *float64
	if body.PresentValue == nil {
		if pointModel.PresentValue != nil {
			presentValue = pointModel.PresentValue
		} else {
			presentValue = body.PresentValue
		}
	} else {
		presentValue = body.PresentValue
	}
	value := *body.PresentValue
	body.ValueOriginal = &value
	var limitMin *float64
	if body.LimitMin == nil {
		if pointModel.LimitMin != nil {
			limitMin = pointModel.LimitMin
		} else {
			limitMin = body.LimitMin
		}
	}
	var limitMax *float64
	if body.LimitMax == nil {
		if pointModel.LimitMax != nil {
			limitMax = pointModel.LimitMax
		} else {
			limitMax = body.LimitMax
		}
	}
	var scaleInMin *float64
	if body.ScaleInMin == nil {
		if pointModel.ScaleInMin != nil {
			scaleInMin = pointModel.ScaleInMin
		} else {
			scaleInMin = body.ScaleInMin
		}
	}
	var scaleInMax *float64
	if body.ScaleInMax == nil {
		if pointModel.ScaleInMax != nil {
			scaleInMax = pointModel.ScaleInMax
		} else {
			scaleInMax = body.ScaleInMax
		}
	}
	var scaleOutMin *float64
	if body.ScaleOutMin == nil {
		if pointModel.ScaleOutMin != nil {
			scaleOutMin = pointModel.ScaleOutMin
		} else {
			scaleOutMin = body.ScaleOutMin
		}
	}
	var scaleOutMax *float64
	if body.ScaleOutMax == nil {
		if pointModel.ScaleOutMax != nil {
			scaleOutMax = pointModel.ScaleOutMax
		} else {
			scaleOutMax = body.ScaleOutMax
		}
	}
	if body.Priority != nil {
		priority := map[string]interface{}{}
		priorityValue := reflect.ValueOf(*body.Priority)
		typeOfPriority := priorityValue.Type()
		highestPri := utils.NewArray()
		highestValue := utils.NewMap()
		for i := 0; i < priorityValue.NumField(); i++ {
			if priorityValue.Field(i).Type().Kind().String() == "ptr" {
				val := priorityValue.Field(i).Interface().(*float64)
				if val == nil {
					priority[typeOfPriority.Field(i).Name] = nil
				} else {
					highestPri.Add(i)
					highestValue.Set(i, *val)
					priority[typeOfPriority.Field(i).Name] = *val
				}
			}
		}
		notNil := false
		for _, v := range priority { //check if there is a value in priority array
			if v != nil {
				notNil = true
			}
		}
		if notNil {
			min, _ := highestPri.MinMaxInt() //get the highest priority
			val := highestValue.Get(min)     //get the highest priority value
			body.CurrentPriority = &min      //TODO check conversion
			v := val.(float64)
			body.PresentValue = &v //process the units as in temperature conversion
		}
		d.DB.Model(&pointModel.Priority).Updates(&priority)
	}
	presentValue = pointScale(presentValue, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax)
	presentValue = pointRange(presentValue, limitMin, limitMax)
	eval, err := pointEval(presentValue, body.ValueOriginal, pointModel.EvalMode, pointModel.Eval)
	if err != nil {
		log.Errorf("ERROR on point invalid point unit")
		return nil, err
	} else {
		pointModel.PresentValue = eval
	}
	vv, display, ok, err := pointUnits(pointModel)
	if err != nil {
		log.Errorf("ERROR on point invalid point unit")
		return nil, err
	}
	if ok {
		presentValue = &vv
		body.ValueDisplay = display
	}
	if !utils.Unit32NilCheck(pointModel.Decimal) {
		val := utils.RoundTo(*presentValue, *pointModel.Decimal)
		presentValue = &val
	}
	if utils.BoolIsNil(pointModel.IsProducer) && utils.BoolIsNil(body.IsProducer) {
		if compare(pointModel, body) {
			_, err := d.ProducerWrite("point", pointModel)
			if err != nil {
				log.Errorf("ERROR ProducerPointCOV at func UpdatePointByFieldAndType")
				return nil, err
			}
		}
	}
	body.PresentValue = presentValue
	query = d.DB.Model(&pointModel).Updates(&body)
	if !fromPlugin { //stop looping
		plug, err := d.GetPluginIDFromDevice(pointModel.DeviceUUID)
		if err != nil {
			return nil, errors.New("ERROR failed to get plugin uuid")
		}
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, plug.PluginConfId, pointModel.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, pointModel)
		if err != nil {
			return nil, errors.New("ERROR on device eventbus")
		}
	}
	return pointModel, nil
}

// GetPointByField returns the point for the given field ie name or nil.
func (d *GormDatabase) GetPointByField(field string, value string, withChildren bool) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("%s = ? ", field)
	if withChildren { // drop child to reduce json size
		query := d.DB.Where(f, value).Preload("Priority").First(&pointModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return pointModel, nil
	} else {
		query := d.DB.Where(f, value).First(&pointModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return pointModel, nil
	}
}

// PointAndQuery will do an SQL AND
func (d *GormDatabase) PointAndQuery(value1 string, value2 string) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("object_type = ? AND address_id = ?")
	query := d.DB.Where(f, value1, value2).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

// UpdatePointByFieldAndUnit get by field and update.
func (d *GormDatabase) UpdatePointByFieldAndUnit(field string, value string, body *model.Point, writeValue bool) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("%s = ? AND unit_type = ?", field)
	query := d.DB.Where(f, value, body.UnitType).Preload("Priority").Find(&pointModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	if utils.BoolIsNil(pointModel.IsProducer) {
		if compare(pointModel, body) {
			_, err := d.ProducerWrite("point", pointModel)
			if err != nil {
				log.Errorf("ERROR ProducerPointCOV at func UpdatePointByFieldAndType")
				return nil, err
			}
		}
	}
	return pointModel, nil
}

// DeletePoint delete a Device.
func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	var pointModel *model.Point
	point, err := d.GetPoint(uuid, api.Args{})
	if err != nil {
		return false, errors.New("point not exist")
	}
	query := d.DB.Where("uuid = ? ", uuid).Delete(&pointModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		plug, err := d.GetPluginIDFromDevice(point.DeviceUUID)
		if err != nil {
			return false, errors.New("ERROR failed to get plugin uuid")
		}
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsDeleted, plug.PluginConfId, point.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, point)
		if err != nil {
			return false, errors.New("ERROR on device eventbus")
		}
		return true, nil
	}

}

// DropPoints delete all points.
func (d *GormDatabase) DropPoints() (bool, error) {
	var pointModel *model.Point
	query := d.DB.Where("1 = 1").Delete(&pointModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}
