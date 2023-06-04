package edgecli

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
	"github.com/NubeIO/rubix-os/src/cli/constants"
	"os"
	"path"
	"path/filepath"
)

func (inst *Client) PluginUpload(body *interfaces.Plugin) (*interfaces.Message, error) {
	uploadLocation := global.Installer.GetAppPluginDownloadPath()

	url := fmt.Sprintf("/api/dirs/create?path=%s", uploadLocation)
	_, _ = nresty.FormatRestyResponse(inst.Rest.R().Post(url))

	pluginFile, err := global.Installer.GetPluginsStorePluginFile(interfaces.Plugin{
		Name:      body.Name,
		Arch:      body.Arch,
		Version:   body.Version,
		Extension: body.Extension,
	})
	if err != nil {
		return nil, err
	}
	tmpDir, err := global.Installer.MakeTmpDirUpload()
	if err != nil {
		return nil, err
	}
	fileDetails, err := fileutils.Unzip(pluginFile, tmpDir, os.FileMode(global.Installer.FileMode))
	if err != nil {
		return nil, err
	}
	if len(fileDetails) != 1 {
		return nil, errors.New(fmt.Sprintf("plugins extraction count mismatch %d", len(fileDetails)))
	}
	extractedPluginFile := path.Join(tmpDir, fileDetails[0].Name)
	reader, err := os.Open(extractedPluginFile)
	if err != nil {
		return nil, err
	}
	url = fmt.Sprintf("/api/files/upload?destination=%s", uploadLocation)
	_, err = nresty.FormatRestyResponse(inst.Rest.R().
		SetFileReader("file", filepath.Base(extractedPluginFile), reader).
		Post(url))
	if err != nil {
		return nil, err
	}
	if err = fileutils.RmRF(tmpDir); err != nil {
		return nil, err
	}
	return &interfaces.Message{Message: "successfully uploaded the plugin"}, nil
}

func (inst *Client) ListPlugins() ([]interfaces.Plugin, error, error) {
	p := global.Installer.GetPluginInstallationPath(constants.RubixOs)
	files, connectionErr, requestErr := inst.ListFilesV2(p)
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	var plugins []interfaces.Plugin
	for _, file := range files {
		plugins = append(plugins, *global.Installer.GetPluginDetails(file.Name))
	}
	return plugins, nil, nil
}
