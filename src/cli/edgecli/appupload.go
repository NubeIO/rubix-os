package edgecli

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/NubeIO/rubix-os/src/cli/edgebioscli/ebmodel"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func (inst *Client) AppUpload(body *interfaces.AppUpload) (*interfaces.Message, error) {
	url := fmt.Sprintf("/api/files/delete-all?path=%s", global.Installer.GetAppDownloadPath(body.Name))
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Delete(url))

	uploadLocation := global.Installer.GetAppDownloadPathWithVersion(body.Name, body.Version)
	url = fmt.Sprintf("/api/dirs/create?path=%s", uploadLocation)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Post(url))

	appStoreFile, err := findAppOnAppStoreFile(body.Name, body.Arch, body.Version)
	if err != nil {
		return nil, err
	}
	url = fmt.Sprintf("/api/files/upload?destination=%s", uploadLocation)
	reader, err := os.Open(*appStoreFile)
	if err != nil {
		return nil, err
	}
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&ebmodel.UploadResponse{}).
		SetFileReader("file", filepath.Base(*appStoreFile), reader).
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
	_, err = nresty.FormatRestyResponse(inst.Rest.R().Delete(url))
	if err != nil {
		return nil, err
	}

	if body.MoveExtractedFileToNameApp {
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
			return &interfaces.Message{Message: "uploaded successfully"}, nil
		}
	}
	if body.MoveOneLevelInsideFileToOutside {
		tmpFolder := global.Installer.GetEmptyNewTmpFolder()
		if unzippedFiles != nil && len(*unzippedFiles) > 0 {
			extractedFile := path.Join(uploadLocation, (*unzippedFiles)[0].Name)
			url = fmt.Sprintf("/api/files/list?path=%s", extractedFile)
			resp, err = nresty.FormatRestyResponse(inst.Rest.R().
				SetResult(&[]fileutils.FileDetails{}).
				Get(url))
			if err != nil {
				return nil, err
			}
			files := resp.Result().(*[]fileutils.FileDetails)
			for _, file := range *files {
				if file.IsDir {
					from := path.Join(extractedFile, file.Name)
					url = fmt.Sprintf("/api/files/move?from=%s&to=%s", from, tmpFolder)
					resp, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url))
					if err != nil {
						return nil, err
					}

					to := global.Installer.GetAppDownloadPathWithVersion(body.Name, body.Version)
					url = fmt.Sprintf("/api/files/delete-all?path=%s", to)
					_, err = nresty.FormatRestyResponse(inst.Rest.R().Delete(url))
					if err != nil {
						return nil, err
					}

					url = fmt.Sprintf("/api/files/move?from=%s&to=%s", tmpFolder, to)
					_, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url))
					if err != nil {
						return nil, err
					}
					return &interfaces.Message{Message: "uploaded successfully"}, nil
				}
			}
		}
	}
	return nil, nil
}

func findAppOnAppStoreFile(appName, arch, version string) (*string, error) {
	storePath := global.Installer.GetAppsStoreAppPathWithArchVersion(appName, arch, version)
	files, err := ioutil.ReadDir(storePath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New(fmt.Sprintf("%s store file doesn't exist (arch: %s, version: %s)", appName, arch, version))
	}
	appStoreFile := path.Join(storePath, files[0].Name())
	return &appStoreFile, nil
}
