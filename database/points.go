package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

// GetPoints returns all devices.
func (d *GormDatabase) GetPoints(withChildren bool) ([]*model.Point, error) {
	var pointsModel []*model.Point
	if withChildren { // drop child to reduce json size
		query := d.DB.Find(&pointsModel);if query.Error != nil {
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
	var pointModel *model.Point
	fmt.Println(1010101)
	fmt.Println(eventbus.BusContext.Value(uuid))
	fmt.Println(1010101)
	if withChildren { // drop child to reduce json size
		query := d.DB.Where("uuid = ? ", uuid).First(&pointModel);if query.Error != nil {
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
func (d *GormDatabase) CreatePoint( body *model.Point) (*model.Point, error) {
	var deviceModel *model.Device
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Point)
	deviceUUID := body.DeviceUUID
	query := d.DB.Where("uuid = ? ", deviceUUID).First(&deviceModel);if query.Error != nil {
		return nil, query.Error
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return  nil, query.Error
	}
	busUpdate(body.UUID, "create", body)

	return body, query.Error
}


// UpdatePoint returns the device for the given id or nil.
func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Find(&pointModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&pointModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	busUpdate(pointModel.UUID, "updates", pointModel)
	return pointModel, nil
}

// DeletePoint delete a Device.
func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ? ", uuid).Delete(&pointModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		busUpdate(pointModel.UUID, "delete", pointModel)
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


/*
update a point value
need network, device and point uuid
need to check if point device and network are enabled
if all are enabled then check COV
if COV is out of range to update db and publish a message on the eventbus
send data to gateway
publish to MQTT if there is an external subscriber
*/



var GetDatabaseBus eventbus.NotificationService

func DataBus() {
	notificationService := eventbus.NewNotificationService(eventbus.BUS)
	GetDatabaseBus = notificationService

}


func busUpdate(UUID string, action string, body *model.Point){
	notificationService := eventbus.NewNotificationService(eventbus.BUS)
	notificationService.Emit(eventbus.BusContext, eventbus.PointUpdated, body)
	fmt.Println("topics", eventbus.BUS.Topics())
}

