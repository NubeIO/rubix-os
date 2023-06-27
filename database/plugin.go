package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
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

func (d *GormDatabase) GetPluginByPathTransaction(db *gorm.DB, path string) (*model.PluginConf, error) {
	var plugin *model.PluginConf
	query := db.Where("module_path = ? ", path).First(&plugin)
	if query.Error != nil {
		return nil, query.Error
	}
	return plugin, nil
}

func (d *GormDatabase) GetPluginByPath(path string) (*model.PluginConf, error) {
	return d.GetPluginByPathTransaction(d.DB, path)
}

func (d *GormDatabase) CreatePlugin(p *model.PluginConf) error {
	p.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Plugin)
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

func (d *GormDatabase) UpdatePluginConfStorage(path string, data []byte) error {
	var plugin *model.PluginConf
	query := d.DB.Where("module_path = ? ", path).First(&plugin)
	if query.Error != nil {
		return query.Error
	}
	plugin.Storage = data
	return d.DB.Model(&plugin).Updates(&plugin).Error
}
