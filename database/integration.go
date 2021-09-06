package database

import (
	"errors"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


type Integration struct {
	*model.Integration
}


// GetIntegrationsList get all of them
func (d *GormDatabase) GetIntegrationsList() ([]*model.Integration, error) {
	var producersModel []*model.Integration
	query := d.DB.Find(&producersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return producersModel, nil
}

// CreateIntegration make it
func (d *GormDatabase) CreateIntegration(body *model.Integration) (*model.Integration, error) {
	body.UUID = utils.MakeTopicUUID("")
	body.Name = nameIsNil(body.Name)
	body.PluginName = pluginIsNil(body.PluginName)
	body.IntegrationType = typeIsNil(body.IntegrationType, "mqtt")
	p, err := d.GetPluginByPath(body.PluginName);if err != nil { //the integration can be added by the pluginName
		return nil, errors.New("invalid plugin name or id")
	}
	body.PluginConfId = p.UUID
	query := d.DB.Create(body);if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

// GetIntegration get it
func (d *GormDatabase) GetIntegration(uuid string) (*model.Integration, error) {
	var wcm *model.Integration
	query := d.DB.Where("uuid = ? ", uuid).First(&wcm); if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

// DeleteIntegration deletes it
func (d *GormDatabase) DeleteIntegration(uuid string) (bool, error) {
	var wcm *model.Integration
	query := d.DB.Where("uuid = ? ", uuid).Delete(&wcm);if query.Error != nil {
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
	var wcm *model.Integration
	query := d.DB.Where("uuid = ?", uuid).Find(&wcm);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&wcm).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

// DropIntegrationsList delete all.
func (d *GormDatabase) DropIntegrationsList() (bool, error) {
	var wcm *model.Integration
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

