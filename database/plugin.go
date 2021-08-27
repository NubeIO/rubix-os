package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"gorm.io/gorm"
)

// GetPluginConfByUser gets plugin configurations from a user.
func (d *GormDatabase) GetPluginConfByUser(userid uint) ([]*model.PluginConf, error) {
	var plugins []*model.PluginConf
	err := d.DB.Where("user_id = ?", userid).Find(&plugins).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return plugins, err
}

func (d *GormDatabase) GetPlugins() ([]*model.PluginConf, error) {
	var plugins []*model.PluginConf
	query := d.DB.Find(&plugins)
	if query.Error != nil {
		return nil, query.Error
	}
	return plugins, nil
}


// GetPluginByPath get object by uuid.
func (d *GormDatabase) GetPluginByPath(path string) (*model.PluginConf, error) {
	var plugin *model.PluginConf
	query := d.DB.Where("module_path = ? ", path).First(&plugin); if query.Error != nil {
		return nil, query.Error
	}
	return plugin, nil
}


// GetPlugin get object by uuid.
func (d *GormDatabase) GetPlugin(uuid string) (*model.PluginConf, error) {
	var plugin *model.PluginConf
	query := d.DB.Where("uuid = ? ", uuid).First(&plugin); if query.Error != nil {
		return nil, query.Error
	}
	return plugin, nil
}


// GetPluginConfByUserAndPath gets plugin configuration by user and file name.
func (d *GormDatabase) GetPluginConfByUserAndPath(userid uint, path string) (*model.PluginConf, error) {
	plugin := new(model.PluginConf)
	err := d.DB.Where("user_id = ? AND module_path = ?", userid, path).First(plugin).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if plugin.ModulePath == path {
		return plugin, err
	}
	return nil, err
}

// GetPluginConfByApplicationID gets plugin configuration by its internal appid.
func (d *GormDatabase) GetPluginConfByApplicationID(appid uint) (*model.PluginConf, error) {
	plugin := new(model.PluginConf)
	err := d.DB.Where("application_id = ?", appid).First(plugin).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if plugin.ApplicationID == appid {
		return plugin, err
	}
	return nil, err
}

// CreatePluginConf creates a new plugin configuration.
func (d *GormDatabase) CreatePluginConf(p *model.PluginConf) error {
	p.UUID = utils.MakeTopicUUID(model.CommonNaming.Plugin)
	return d.DB.Create(p).Error
}

// GetPluginConfByToken gets plugin configuration by plugin token.
func (d *GormDatabase) GetPluginConfByToken(token string) (*model.PluginConf, error) {
	plugin := new(model.PluginConf)
	err := d.DB.Where("token = ?", token).First(plugin).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if plugin.Token == token {
		return plugin, err
	}
	return nil, err
}

// GetPluginConfByID gets plugin configuration by plugin ID.
func (d *GormDatabase) GetPluginConfByID(id string) (*model.PluginConf, error) {
	plugin := new(model.PluginConf)
	err := d.DB.Where("uuid = ?", id).First(plugin).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if plugin.UUID == id {
		return plugin, err
	}
	return nil, err
}

// UpdatePluginConf updates plugin configuration.

func (d *GormDatabase) UpdatePluginConf(p *model.PluginConf) error {
	return d.DB.Save(p).Error
}

// DeletePluginConfByID deletes a plugin configuration by its id.
func (d *GormDatabase) DeletePluginConfByID(id string) error {
	return d.DB.Where("uuid = ?", id).Delete(&model.PluginConf{}).Error
}
