package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"time"
)




var subscriberChildTable = "PointSubscriberLedger"
var SubscriptionChildTable = "PointSubscriptionLedger"

// GetPoints returns all devices.
func (d *GormDatabase) GetPoints(withChildren bool) ([]*model.Point, error) {
	var pointsModel []*model.Point
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload(subscriberChildTable).Preload(SubscriptionChildTable).Find(&pointsModel);if query.Error != nil {
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
	if withChildren { // drop child to reduce json size
		query := d.DB.Where("uuid = ? ", uuid).Preload(subscriberChildTable).First(&pointModel);if query.Error != nil {
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
	return body, query.Error
}


// UpdatePoint returns the device for the given id or nil.
func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Preload(subscriberChildTable).Preload(SubscriptionChildTable).Where("uuid = ?", uuid).Find(&pointModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&pointModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	gatewayUUID := pointModel.PointSubscriberLedger
	for _, e := range gatewayUUID {
		busUpdate(e.UUID, pointModel)
	}

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


func busUpdate(UUID string, body *model.Point){
	fmt.Println(123123123123)
	payload := new(eventbus.BusPayload)
	payload.GatewayUUID = UUID
	payload.ThingName = body.Name
	payload.MessageString = "what up"
	payload.MessageTS = time.Now().Format(time.RFC850)
	topic := "gateway"
	err := eventbus.BUS.Emit(eventbus.BusBackground, topic, payload)
	fmt.Println("topics", eventbus.BUS.Topics())
	if err != nil {
		fmt.Println("error", err)
	}
}