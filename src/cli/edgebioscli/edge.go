package edgebioscli

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
	"path/filepath"
)

const rubixEdgeName = "rubix-edge"

func (inst *BiosClient) RubixEdgeUpload(body *interfaces.FileUpload) (*interfaces.Message, error) {
	uploadLocation := fmt.Sprintf("/data/rubix-service/apps/download/%s/%s", rubixEdgeName, body.Version)
	url := fmt.Sprintf("/api/dirs/create?path=%s", uploadLocation)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		Post(url))

	url = fmt.Sprintf("/api/files/upload?destination=%s", uploadLocation)
	reader, err := os.Open(body.File)
	if err != nil {
		return nil, err
	}
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&ebmodel.UploadResponse{}).
		SetFileReader("file", filepath.Base(body.File), reader).
		Post(url))
	if err != nil {
		return nil, err
	}
	upload := resp.Result().(*ebmodel.UploadResponse)

	url = fmt.Sprintf("/api/zip/unzip?source=%s&destination=%s", upload.Destination, uploadLocation)
	resp, err = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	unzippedFiles := resp.Result().(*[]fileutils.FileDetails)

	url = fmt.Sprintf("/api/files/delete?file=%s", upload.Destination)
	resp, err = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		Delete(url))
	if err != nil {
		return nil, err
	}

	for _, f := range *unzippedFiles {
		from := path.Join(uploadLocation, f.Name)
		to := path.Join(uploadLocation, "app")
		url = fmt.Sprintf("/api/files/move?from=%s&to=%s", from, to)
		resp, err = nresty.FormatRestyResponse(inst.Rest.R().
			SetResult(&interfaces.Message{}).
			Post(url))
		if err != nil {
			return nil, err
		}
	}
	return &interfaces.Message{Message: "successfully uploaded the rubix-edge in edge device"}, nil
}

func (inst *BiosClient) RubixEdgeInstall(version string) (*interfaces.Message, error) {
	// delete installed files
	installationDirectory := fmt.Sprintf("/data/rubix-service/apps/install/%s", rubixEdgeName)
	url := fmt.Sprintf("/api/files/delete-all?path=%s", installationDirectory)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		Delete(url))
	log.Println("deleted installed files, if any")

	downloadedFile := fmt.Sprintf("/data/rubix-service/apps/download/%s/%s/app", rubixEdgeName, version)
	installationFile := fmt.Sprintf("/data/rubix-service/apps/install/%s/%s/app", rubixEdgeName, version)

	// create installation directory
	installationDirectoryWithVersion := filepath.Dir(installationFile)
	url = fmt.Sprintf("/api/dirs/create?path=%s", installationDirectoryWithVersion)
	_, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	log.Info("created installation directory")

	// move downloaded file to installation directory
	url = fmt.Sprintf("/api/files/move?from=%s&to=%s", downloadedFile, installationFile)
	_, err = nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&interfaces.Message{}).
		Post(url))
	if err != nil {
		return nil, err
	}
	log.Info("moved downloaded file to installation directory")

	tmpDir, absoluteServiceFileName, err := systemctl.GenerateServiceFile(&systemctl.ServiceFile{
		Name:                        rubixEdgeName,
		Version:                     version,
		ExecStart:                   "app -p 1661 -r /data -a rubix-edge -d data -c config --prod server",
		AttachWorkingDirOnExecStart: true,
	}, global.Installer)
	if err != nil {
		return nil, err
	}
	log.Info("created service file locally")

	message, err := inst.installServiceFile(rubixEdgeName, absoluteServiceFileName)
	if err != nil {
		return message, err
	}
	err = fileutils.RmRF(tmpDir)
	if err != nil {
		log.Errorf("delete tmp generated service file %s", absoluteServiceFileName)
	}
	log.Infof("deleted tmp generated local service file %s", absoluteServiceFileName)
	return &interfaces.Message{Message: "successfully installed the rubix-edge in edge device"}, nil
}

func (inst *BiosClient) installServiceFile(appName, absoluteServiceFileName string) (*interfaces.Message, error) {
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
