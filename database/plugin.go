package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"gorm.io/gorm"
)

func (d *GormDatabase) GetPlugins() ([]*model.PluginConf, error) {
	var plugins []*model.PluginConf
	query := d.DB.Find(&plugins)
	if query.Error != nil {
		return nil, query.Error
	}
	return plugins, nil
}

func (d *GormDatabase) GetPluginByPath(path string) (*model.PluginConf, error) {
	var plugin *model.PluginConf
	query := d.DB.Where("module_path = ? ", path).First(&plugin)
	if query.Error != nil {
		return nil, query.Error
	}
	return plugin, nil
}

func (d *GormDatabase) CreatePlugin(p *model.PluginConf) error {
	p.UUID = utils.MakeTopicUUID(model.CommonNaming.Plugin)
	return d.DB.Create(p).Error
}

func (d *GormDatabase) GetPlugin(id string) (*model.PluginConf, error) {
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

func (d *GormDatabase) UpdatePluginConf(p *model.PluginConf) error {
	return d.DB.Save(p).Error
}
