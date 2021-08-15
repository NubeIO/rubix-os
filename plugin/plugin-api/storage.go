package plugin

import "github.com/NubeDev/plug-framework/model"

// StorageHandler consists callbacks used to perform read/writes to the persistent storage for plugins.
type StorageHandler interface {
	Save(b []byte) error
	Load() ([]byte, error)
	GetNet() ([]*model.Network, error)
}

// Storager is the interface plugins should implement to use persistent storage.
type Storager interface {
	Plugin
	SetStorageHandler(handler StorageHandler)
}
