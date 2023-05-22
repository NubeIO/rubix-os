package installer

import (
	"path"
)

func (inst *Installer) GetAppsStorePath() string {
	return path.Join(inst.StoreDir, "apps") // /data/store/apps
}

func (inst *Installer) GetAppsStoreAppPath(appName string) string {
	return path.Join(inst.GetAppsStorePath(), appName) // /data/store/apps/<app_name>
}

func (inst *Installer) GetAppsStoreAppPathWithArchVersion(appName, arch, version string) string {
	return path.Join(inst.GetAppsStoreAppPath(appName), arch, version) // /data/store/apps/<app_name>/<arch>/<version>
}

func (inst *Installer) GetPluginsStorePath() string {
	return path.Join(inst.StoreDir, "plugins") // /data/store/plugins
}

func (inst *Installer) GetPluginsStoreWithFile(fileName string) string {
	p := path.Join(inst.GetPluginsStorePath(), fileName) // /data/store/plugins/<plugin_file>
	return p
}

func (inst *Installer) GetPluginInstallationPath(appName string) string {
	appDataPath := inst.GetAppDataPath(appName)
	return path.Join(appDataPath, "data/plugins") // /data/flow-framework/plugins
}
