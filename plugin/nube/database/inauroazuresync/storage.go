package main

import (
	"encoding/json"
	"time"
)

type PluginConfStorage struct { // defines the structure of the plugin storage (stored as []byte)
	LastSyncByGateway map[string]time.Time `json:"lastSyncByGateway"` // stores the last sync time (to azure iot hub) by gateway/host uuid.
}

func (inst *Instance) getPluginConfStorage() (*PluginConfStorage, error) {
	conf, err := inst.db.GetPluginByPath(name)
	if err != nil {
		return nil, err
	}
	var storedStruct PluginConfStorage
	err = json.Unmarshal(conf.Storage, &storedStruct)
	if err != nil {
		return nil, err
	}
	return &storedStruct, nil
}

func (inst *Instance) setPluginConfStorage(storeStruct *PluginConfStorage) error {
	storageBytes, err := json.Marshal(storeStruct)
	if err != nil {
		return err
	}

	err = inst.db.UpdatePluginConfStorage(name, storageBytes)
	if err != nil {
		return err
	}
	return nil
}
