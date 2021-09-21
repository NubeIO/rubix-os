package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

// GetPoints returns all devices.
func (d *GormDatabase) GetPoints(withChildren bool) ([]*model.Point, error) {
	var pointsModel []*model.Point
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload("Priority").Find(&pointsModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return pointsModel, nil
	} else {
		query := d.DB.Find(&pointsModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return pointsModel, nil
	}
}

// GetPoint returns the device for the given id or nil.
func (d *GormDatabase) GetPoint(uuid string, withChildren bool) (*model.Point, error) {
	var pointModel *model.Point
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload("Priority").Where("uuid = ? ", uuid).First(&pointModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return pointModel, nil
	} else {
		query := d.DB.Where("uuid = ? ", uuid).First(&pointModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return pointModel, nil
	}
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

// CreatePoint creates a device.
func (d *GormDatabase) CreatePoint(body *model.Point, streamUUID string) (*model.Point, error) {
	var deviceModel *model.Device
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Point)
	deviceUUID := body.DeviceUUID
	body.Name = nameIsNil(body.Name)
	_, err := checkObjectType(body.ObjectType)
	if err != nil {
		return nil, err
	}
	query := d.DB.Where("uuid = ? ", deviceUUID).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	//check if there is an existing device with this address code
	_, existing := d.PointDeviceByAddressID("", body)
	if existing {
		return nil, errors.New("an existing point of that ObjectType & id exists")
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
func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point, writeValue, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	//TODO add in a check to make sure user doesn't set the addressID and the ObjectType the same as another point
	//check if there is an existing device with this address code
	_, err := checkObjectType(body.ObjectType)
	if err != nil {
		return nil, err
	}
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").Find(&pointModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	if writeValue {
		//TODO add this in to save a few DB read/write for the priority
		//query = d.DB.Model(&pointModel.Priority).Updates(&body.Priority)
		//query = d.DB.Model(&pointModel).Updates(&body)
	}
	query = d.DB.Model(&pointModel.Priority).Updates(&body.Priority)
	query = d.DB.Model(&pointModel).Updates(&body)

	if *pointModel.IsProducer && *body.IsProducer {
		if compare(pointModel, body) {
			_, err := d.ProducerWrite("point", pointModel)
			if err != nil {
				log.Errorf("ERROR ProducerPointCOV at func UpdatePointByFieldAndType")
				return nil, err
			}
		}
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

// UpdatePointByFieldAndType get by field and update.
func (d *GormDatabase) UpdatePointByFieldAndType(field string, value string, body *model.Point, writeValue bool) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("%s = ? AND thing_type = ?", field)
	query := d.DB.Where(f, value, body.ThingType).Preload("Priority").Find(&pointModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	if *pointModel.IsProducer {
		if compare(pointModel, body) {
			log.Errorf("UpdatePointByFieldAndType")
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
	point, err := d.GetPoint(uuid, false)
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
		fmt.Println(point.DeviceUUID, point)
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
