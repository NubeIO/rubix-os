package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type Integration struct {
	*model.Integration
}

// GetIntegrations get all of them
func (d *GormDatabase) GetIntegrations() ([]*model.Integration, error) {
	var integrations []*model.Integration
	query := d.DB.Find(&integrations)
	if query.Error != nil {
		return nil, query.Error
	}
	return integrations, nil
}

// CreateIntegration make it
func (d *GormDatabase) CreateIntegration(body *model.Integration) (*model.Integration, error) {
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Integration)
	body.Name = nameIsNil(body.Name)
	body.PluginName = pluginIsNil(body.PluginName)
	body.IntegrationType = typeIsNil(body.IntegrationType, "mqtt")
	p, err := d.GetPluginByPath(body.PluginName)
	if err != nil { //the integration can be added by the pluginName
		return nil, errors.New("invalid plugin name or id")
	}
	body.PluginConfId = p.UUID
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetIntegration get it
func (d *GormDatabase) GetIntegration(uuid string) (*model.Integration, error) {
	var integration *model.Integration
	query := d.DB.Where("uuid = ? ", uuid).First(&integration)
	if query.Error != nil {
		return nil, query.Error
	}
	return integration, nil
}

// GetIntegrationByName get it by name
func (d *GormDatabase) GetIntegrationByName(name string) (*model.Integration, error) {
	var integration *model.Integration
	query := d.DB.Where("name = ? ", name).First(&integration)
	if query.Error != nil {
		return nil, query.Error
	}
	return integration, nil
}

// GetEnabledIntegrationByPluginConfId get it
func (d *GormDatabase) GetEnabledIntegrationByPluginConfId(pcId string) ([]*model.Integration, error) {
	var integration []*model.Integration
	query := d.DB.Where("plugin_conf_id = ? ", pcId).
		Where("enable = ?", true).Find(&integration)
	if query.Error != nil {
		return nil, query.Error
	}
	return integration, nil
}

// DeleteIntegration deletes it
func (d *GormDatabase) DeleteIntegration(uuid string) (bool, error) {
	var integration *model.Integration
	query := d.DB.Where("uuid = ? ", uuid).Delete(&integration)
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

// UpdateIntegration  update it
func (d *GormDatabase) UpdateIntegration(uuid string, body *model.Integration) (*model.Integration, error) {
	var integration *model.Integration
	query := d.DB.Where("uuid = ?", uuid).First(&integration)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&integration).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return integration, nil
}

// DropIntegrationsList delete all.
func (d *GormDatabase) DropIntegrationsList() (bool, error) {
	var integration *model.Integration
	query := d.DB.Where("1 = 1").Delete(&integration)
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
