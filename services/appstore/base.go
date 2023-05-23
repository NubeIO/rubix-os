package appstore

import (
	"errors"
	"github.com/NubeIO/rubix-os/global"
	"os"
)

type StoreDatabase interface {
}
type Store struct {
}

func New(store *Store) (*Store, error) {
	if store == nil {
		return nil, errors.New("app store can not be empty")
	}
	err := store.initMakeAllDirs()
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (inst *Store) initMakeAllDirs() error {
	if err := os.MkdirAll(global.Installer.GetAppsStorePath(), os.FileMode(global.Installer.FileMode)); err != nil {
		return err
	}
	if err := os.MkdirAll(global.Installer.GetPluginsStorePath(), os.FileMode(global.Installer.FileMode)); err != nil {
		return err
	}
	return nil
}
