package database

import (
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type MqttConnection struct {
	*model.MqttConnection
}

// GetMqttConnectionsList get all of them
func (d *GormDatabase) GetMqttConnectionsList() ([]*model.MqttConnection, error) {
	var producersModel []*model.MqttConnection
	query := d.DB.Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateMqttConnection make it
func (d *GormDatabase) CreateMqttConnection(body *model.MqttConnection) (*model.MqttConnection, error) {
	body.UUID = utils.MakeTopicUUID("")
	body.Name = nameIsNil(body.Name)
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetMqttConnection get it
func (d *GormDatabase) GetMqttConnection(uuid string) (*model.MqttConnection, error) {
	var wcm *model.MqttConnection
	query := d.DB.Where("uuid = ? ", uuid).First(&wcm)
	if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

// DeleteMqttConnection deletes it
func (d *GormDatabase) DeleteMqttConnection(uuid string) (bool, error) {
	var wcm *model.MqttConnection
	query := d.DB.Where("uuid = ? ", uuid).Delete(&wcm)
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

// UpdateMqttConnection  update it
func (d *GormDatabase) UpdateMqttConnection(uuid string, body *model.MqttConnection) (*model.MqttConnection, error) {
	var wcm *model.MqttConnection
	query := d.DB.Where("uuid = ?", uuid).First(&wcm)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&wcm).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

// DropMqttConnectionsList delete all.
func (d *GormDatabase) DropMqttConnectionsList() (bool, error) {
	var wcm *model.MqttConnection
	query := d.DB.Where("1 = 1").Delete(&wcm)
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
