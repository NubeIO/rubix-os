package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
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

// CreatePoint creates a device.
func (d *GormDatabase) CreatePoint(body *model.Point) (*model.Point, error) {
	var deviceModel *model.Device
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Point)
	deviceUUID := body.DeviceUUID
	body.Name = nameIsNil(body.Name)
	query := d.DB.Where("uuid = ? ", deviceUUID).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	body.CommonEnable.Enable = true
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
	return body, query.Error
}

// UpdatePoint returns the device for the given id or nil.
func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point, writeValue bool) (*model.Point, error) {
	var pointModel *model.Point

	if writeValue {
		//TODO point cov event
	}
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").Find(&pointModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&pointModel.Priority).Updates(&body.Priority)
	query = d.DB.Model(&pointModel).Updates(&body)
	if query.Error != nil {
		return nil, query.Error
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

// UpdatePointByField get by field and update.
func (d *GormDatabase) UpdatePointByField(field string, value string, body *model.Point, writeValue bool) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).Find(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if writeValue {
		if body.IsProducer {
			//producer, err := d.UpdateProducerByField("producer_thing_uuid", pointModel.UUID)
			//if err != nil {
			//	return nil, err
			//}

		}

		query = d.DB.Model(&pointModel).Updates(body)
		if query.Error != nil {
			return nil, query.Error
		}
	}

	//query := d.DB.Where(f, value).Find(&pointModel)
	//if query.Error != nil {
	//	return nil, query.Error
	//}
	//query = d.DB.Model(&pointModel).Updates(body)
	//if query.Error != nil {
	//	return nil, query.Error
	//}
	return pointModel, nil
}

// DeletePoint delete a Device.
func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ? ", uuid).Delete(&pointModel)
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
