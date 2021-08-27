package plugin

import "github.com/NubeDev/flow-framework/model"

type dbStorageHandler struct {
	pluginID string
	db       Database
}

func (c dbStorageHandler) Save(b []byte) error {
	conf, err := c.db.GetPluginConfByID(c.pluginID)
	if err != nil {
		return err
	}
	conf.Storage = b
	return c.db.UpdatePluginConf(conf)
}

func (c dbStorageHandler) Load() ([]byte, error) {
	pluginConf, err := c.db.GetPluginConfByID(c.pluginID)
	if err != nil {
		return nil, err
	}
	return pluginConf.Storage, nil
}

func (c dbStorageHandler) GetNet() ([]*model.Network, error) {
	net, err := c.db.GetNetworks(false, false)
	if err != nil {
		return nil, err
	}
	return net, err
}
