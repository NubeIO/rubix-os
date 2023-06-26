package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

type MqttConnection struct {
	*model.MqttConnection
}

// GetMqttConnectionsList get all of them
func (d *GormDatabase) GetMqttConnectionsList() ([]*model.MqttConnection, error) {
	var mqttConnectionsModel []*model.MqttConnection
	query := d.DB.Find(&mqttConnectionsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return mqttConnectionsModel, nil
}

// CreateMqttConnection make it
func (d *GormDatabase) CreateMqttConnection(body *model.MqttConnection) (*model.MqttConnection, error) {
	body.UUID = nuuid.MakeTopicUUID("")
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
	return d.deleteResponseBuilder(query)
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
