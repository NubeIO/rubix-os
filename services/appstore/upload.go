package appstore

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	"os"
	"path"
)

type UploadResponse struct {
	Name         string `json:"name,omitempty"`
	Version      string `json:"version,omitempty"`
	UploadedOk   bool   `json:"uploaded_ok,omitempty"`
	TmpFile      string `json:"tmp_file,omitempty"`
	UploadedFile string `json:"uploaded_file,omitempty"`
}

func (inst *Store) UploadAddOnAppStore(app *interfaces.Upload) (*UploadResponse, error) {
	if app.Name == "" {
		return nil, errors.New("app_name can not be empty")
	}
	if app.Version == "" {
		return nil, errors.New("app_version can not be empty")
	}
	if app.Arch == "" {
		return nil, errors.New("arch_type can not be empty, try armv7 amd64")
	}
	err := os.MkdirAll(global.Installer.GetAppsStoreAppPathWithArchVersion(app.Name, app.Arch, app.Version), os.FileMode(global.Installer.FileMode))
	if err != nil {
		return nil, err
	}
	var file = app.File
	uploadResp := &UploadResponse{
		Name:         app.Name,
		Version:      app.Version,
		UploadedOk:   false,
		TmpFile:      "",
		UploadedFile: "",
	}
	resp, err := global.Installer.Upload(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("upload app: %s", err.Error()))
	}
	uploadResp.TmpFile = resp.TmpFile
	source := resp.UploadedFile
	destination := path.Join(global.Installer.GetAppsStoreAppPathWithArchVersion(app.Name, app.Arch, app.Version), resp.FileName)
	check := fileutils.FileExists(source)
	if !check {
		return nil, errors.New(fmt.Sprintf("upload file tmp dir not found: %s", source))
	}
	uploadResp.UploadedFile = destination
	err = os.Rename(source, destination)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("move build error: %s", err.Error()))
	}
	uploadResp.UploadedOk = true
	return uploadResp, nil
}
