package edgecli

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/NubeIO/rubix-os/services/systemctl"
	"github.com/NubeIO/rubix-os/src/cli/constants"
	"github.com/NubeIO/rubix-os/src/cli/edgebioscli/ebmodel"
	"github.com/NubeIO/rubix-os/utils/namings"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func (inst *Client) AppInstall(app *systemctl.ServiceFile) (*interfaces.Message, error) {
	err := inst.backupAppDataDir(app.Name)
	if err != nil {
		return nil, err
	}

	err = inst.moveAppFromDownloadToInstallDir(app)
	if err != nil {
		return nil, err
	}

	installPath := global.Installer.GetAppInstallPath(app.Name)
	url := fmt.Sprintf("/api/files/delete-all?path=%s", installPath)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Delete(url))

	tmpDir, absoluteServiceFileName, err := systemctl.GenerateServiceFile(app, global.Installer)
	_, err = inst.installServiceFile(app.Name, absoluteServiceFileName)
	if err != nil {
		return nil, err
	}
	err = fileutils.RmRF(tmpDir)
	if err != nil {
		log.Errorf("delete tmp generated service file %s", absoluteServiceFileName)
	}
	log.Infof("deleted tmp generated local service file %s", absoluteServiceFileName)
	return &interfaces.Message{Message: "successfully installed the app"}, nil
}

func (inst *Client) moveAppFromDownloadToInstallDir(app *systemctl.ServiceFile) error {
	from := global.Installer.GetAppDownloadPathWithVersion(app.Name, app.Version)
	to := global.Installer.GetAppInstallPathWithVersion(app.Name, app.Version)
	url := fmt.Sprintf("/api/files/delete-all?path=%s", to)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Delete(url))
	url = fmt.Sprintf("/api/dirs/create?path=%s", path.Dir(to))
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Post(url))
	url = fmt.Sprintf("/api/files/move?from=%s&to=%s", from, to)
	_, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url))
	if err != nil {
		return err
	}
	return nil
}

func (inst *Client) installServiceFile(appName, absoluteServiceFileName string) (*interfaces.Message, error) {
	serviceFileName := namings.GetServiceNameFromAppName(appName)
	serviceFile := path.Join(constants.ServiceDir, serviceFileName)
	symlinkServiceFile := path.Join(constants.ServiceDirSoftLink, serviceFileName)
	url := fmt.Sprintf("/api/files/upload?destination=%s", constants.ServiceDir)
	reader, err := os.Open(absoluteServiceFileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error open service file: %s err: %s", absoluteServiceFileName, err.Error()))
	}
	if _, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetFileReader("file", serviceFileName, reader).
		SetResult(&ebmodel.UploadResponse{}).
		Post(url)); err != nil {
		return nil, err
	}
	log.Info("service file is uploaded successfully")

	url = fmt.Sprintf("/api/syscall/unlink?path=%s", symlinkServiceFile)
	if _, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("soft un-linked %s", symlinkServiceFile)

	url = fmt.Sprintf("/api/syscall/link?path=%s&link=%s", serviceFile, symlinkServiceFile)
	if _, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("soft linked %s to %s", serviceFile, symlinkServiceFile)

	url = "/api/systemctl/daemon-reload"
	if _, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("daemon reloaded")

	url = fmt.Sprintf("/api/systemctl/enable?unit=%s", serviceFileName)
	if _, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("enabled service %s", serviceFileName)

	url = fmt.Sprintf("/api/systemctl/restart?unit=%s", serviceFileName)
	if _, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url)); err != nil {
		log.Error(err)
	}
	log.Infof("started service %s", serviceFileName)
	return nil, nil
}

func (inst *Client) backupAppDataDir(appName string) error {
	appVersion, connectionErr, requestErr := inst.getAppVersion(appName)
	if requestErr != nil {
		log.Warnf("we haven't found %s installed version (%s), so skipping backup process", appName, requestErr)
		return nil
	} else if connectionErr != nil {
		return connectionErr
	}

	from := global.Installer.GetAppDataDataPath(appName)
	to := global.Installer.GetAppBackupPath(appName, appVersion)
	if appName == constants.RubixOs { // otherwise, plugins & images folders also gets copied which will be large
		from = path.Join(from, "data.db")
		to = path.Join(to, "data.db")
	}
	url := fmt.Sprintf("/api/files/copy?from=%s&to=%s", from, to)
	log.Infof("backing up from %s to %s", from, to)
	_, connectionErr, requestErr = nresty.FormatRestyV2Response(inst.Rest.R().Post(url))
	if requestErr != nil {
		log.Warnf("backup process from %s to %s got failed (%s)", from, to, requestErr)
	} else if connectionErr != nil {
		return connectionErr
	} else {
		log.Infof("backup process has been completed from %s to %s", from, to)
	}
	return nil
}
