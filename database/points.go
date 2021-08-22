package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


var pointsModel []*model.Point
var pointModel *model.Point
var pointStoreChildTable = "PointStore"

// GetPoints returns all devices.
func (d *GormDatabase) GetPoints(withChildren bool) ([]*model.Point, error) {
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload(pointStoreChildTable).Find(&pointsModel);if query.Error != nil {
			return nil, query.Error
		}
		return pointsModel, nil
	} else {
		query := d.DB.Find(&pointsModel);if query.Error != nil {
			return nil, query.Error
		}
		return pointsModel, nil
	}

}

// GetPoint returns the device for the given id or nil.
func (d *GormDatabase) GetPoint(uuid string, withChildren bool) (*model.Point, error) {
	if withChildren { // drop child to reduce json size
		query := d.DB.Where("uuid = ? ", uuid).Preload(pointStoreChildTable).First(&pointModel);if query.Error != nil {
			return nil, query.Error
		}
		return pointModel, nil
	} else {
		query := d.DB.Where("uuid = ? ", uuid).First(&pointModel); if query.Error != nil {
			return nil, query.Error
		}
		return pointModel, nil
	}
}

// CreatePoint creates a device.
func (d *GormDatabase) CreatePoint( body *model.Point) error {
	body.UUID, _ = utils.MakeUUID()
	deviceUUID := body.DeviceUUID
	query := d.DB.Where("uuid = ? ", deviceUUID).First(&deviceModel);if query.Error != nil {
		return query.Error
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return query.Error
	}
	return query.Error
}


// UpdatePoint returns the device for the given id or nil.
func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point) (*model.Point, error) {
	query := d.DB.Where("uuid = ?", uuid).Find(&pointModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&pointModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

// DeletePoint delete a Device.
func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ? ", uuid).Delete(&pointModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}
