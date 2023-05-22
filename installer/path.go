package installer

import (
	"fmt"
	"github.com/NubeIO/flow-framework/src/cli/constants"
	"github.com/NubeIO/flow-framework/utils/namings"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"os"
	"path"
	"time"
)

func (inst *Installer) GetAppDataPath(appName string) string {
	dataDirName := namings.GetDataDirNameFromAppName(appName)
	return path.Join(inst.RootDir, dataDirName) // /data/rubix-wires
}

func (inst *Installer) GetAppDataDataPath(appName string) string {
	dataDirName := namings.GetDataDirNameFromAppName(appName)
	return path.Join(inst.RootDir, dataDirName, "data") // /data/rubix-wires/data
}

func (inst *Installer) GetAppDataConfigPath(appName string) string {
	dataDirName := namings.GetDataDirNameFromAppName(appName)
	return path.Join(inst.RootDir, dataDirName, "config") // /data/rubix-wires/config
}

func (inst *Installer) GetAppInstallPath(appName string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsInstallDir, repoName) // /data/rubix-service/apps/install/wires-builds
}

func (inst *Installer) GetAppInstallPathWithVersion(appName, version string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsInstallDir, repoName, version) // /data/rubix-service/apps/install/wires-builds/v0.0.1
}

func (inst *Installer) GetAppDownloadPath(appName string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsDownloadDir, repoName) // /data/rubix-service/apps/download/wires-builds
}

func (inst *Installer) GetAppDownloadPathWithVersion(appName, version string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsDownloadDir, repoName, version) // /data/rubix-service/apps/download/wires-builds/v0.0.1
}

func (inst *Installer) GetEmptyNewTmpFolder() string {
	return path.Join(inst.TmpDir, nuuid.ShortUUID("tmp")) // /data/tmp/tmp_45EA34EB
}

func (inst *Installer) MakeTmpDir() error {
	return os.MkdirAll(inst.TmpDir, os.FileMode(inst.FileMode)) // /data/tmp
}

func (inst *Installer) MakeTmpDirUpload() (string, error) {
	tmpDir := inst.GetEmptyNewTmpFolder()
	err := os.MkdirAll(tmpDir, os.FileMode(inst.FileMode)) // /data/tmp/tmp_45EA34EB
	return tmpDir, err
}

func (inst *Installer) GetAppPluginDownloadPath() string {
	repoName := namings.GetRepoNameFromAppName(constants.FlowFramework)
	return path.Join(inst.AppsDownloadDir, repoName, "plugins") // /data/rubix-service/apps/download/flow-framework/plugins
}

func (inst *Installer) GetAppPluginInstallPath() string {
	return path.Join(inst.GetAppDataDataPath(constants.FlowFramework), "plugins") // /data/flow-framework/data/plugins
}

func (inst *Installer) GetAppPluginInstallFilePath(pluginName, arch string) string {
	return path.Join(inst.GetAppPluginInstallPath(), fmt.Sprintf("%s-%s.so", pluginName, arch)) // /data/flow-framework/data/plugins/system-amd64.so
}

func (inst *Installer) GetAppBackupPath(appName, version string) string {
	return path.Join(inst.BackupDir, appName,
		fmt.Sprintf("%s_%s", time.Now().UTC().Format("20060102150405"), version)) // /data/rubix-wires/<time_value>_v0.0.1
}
