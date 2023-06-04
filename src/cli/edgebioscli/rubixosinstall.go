package edgebioscli

import (
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/NubeIO/rubix-os/src/cli/constants"
	"path"
)

func (inst *BiosClient) MoveAppAndPluginsFromDownloadToInstallDir(version string) error { //
	from := global.Installer.GetAppDownloadPathWithVersion(constants.RubixOs, version)
	to := global.Installer.GetAppInstallPathWithVersion(constants.RubixOs, version)
	url := fmt.Sprintf("/api/files/delete-all?path=%s", to)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Delete(url))
	url = fmt.Sprintf("/api/dirs/create?path=%s", path.Dir(to))
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Post(url))
	url = fmt.Sprintf("/api/files/move?from=%s&to=%s", from, to)
	_, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url))
	if err != nil {
		return err
	}

	from = global.Installer.GetAppPluginDownloadPath()
	to = global.Installer.GetAppPluginInstallPath()
	url = fmt.Sprintf("/api/files/delete-all?path=%s", to)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Delete(url))
	url = fmt.Sprintf("/api/dirs/create?path=%s", path.Dir(to))
	_, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url))
	if err != nil {
		return err
	}
	url = fmt.Sprintf("/api/files/move?from=%s&to=%s", from, to)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Post(url)) // ignore error: sometimes from folder will be empty

	return nil
}

func (inst *BiosClient) MovePluginsFromDownloadToInstallDir() (*interfaces.Message, error) {
	from := global.Installer.GetAppPluginDownloadPath()
	to := global.Installer.GetAppPluginInstallPath()
	url := fmt.Sprintf("/api/dirs/create?path=%s", from)
	_, err := nresty.FormatRestyResponse(inst.Rest.R().Post(url))
	if err != nil {
		return nil, err
	}
	url = fmt.Sprintf("/api/dirs/create?path=%s", to)
	_, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url))
	if err != nil {
		return nil, err
	}
	url = fmt.Sprintf("/api/files/list?path=%s", from)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]fileutils.FileDetails{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	files := (resp.Result()).(*[]fileutils.FileDetails)
	if files != nil {
		for _, file := range *files {
			fromFile := path.Join(from, file.Name)
			toFile := path.Join(to, file.Name)
			url = fmt.Sprintf("/api/files/move?from=%s&to=%s", fromFile, toFile)
			_, err = nresty.FormatRestyResponse(inst.Rest.R().Post(url))
			if err != nil {
				return nil, err
			}
		}
	}
	return &interfaces.Message{Message: "transferred plugins from download to install location"}, nil
}
