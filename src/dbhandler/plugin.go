package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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

func (h *Handler) UpdatePluginConfStorage(path string, data []byte) error {
	return getDb().UpdatePluginConfStorage(path, data)
}
