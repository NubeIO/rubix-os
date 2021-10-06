package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetPluginByPath(name string) (*model.PluginConf, error) {
	q, err := getDb().GetPluginByPath(name)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPlugin(uuid string) (*model.PluginConf, error) {
	q, err := getDb().GetPlugin(uuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}
