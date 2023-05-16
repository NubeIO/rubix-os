package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/interfaces"
	"io/ioutil"
	"strings"
)

// GetPluginDetails takes in the name (influx-amd64.so) and returns the info
func (inst *Installer) GetPluginDetails(pluginName string) *interfaces.Plugin {
	parts := strings.Split(pluginName, "-")
	plugin := interfaces.Plugin{}
	for _, part := range parts {
		plugin.Name = parts[0]
		if strings.Contains(part, "amd64") {
			plugin.Arch = "amd64"
		}
		if strings.Contains(part, "armv7") {
			plugin.Arch = "armv7"
		}
		if strings.Contains(part, ".so") {
			plugin.Extension = ".so"
		}
	}
	return &plugin
}

// ValidateBinaryPlugin check if all the details of a binary name is correct (influx-amd64.so)
func (inst *Installer) ValidateBinaryPlugin(pluginName string) error {
	plugin := inst.GetPluginDetails(pluginName)
	if plugin.Name == "" {
		return errors.New(fmt.Sprintf("plugin name is incorrect: %s", pluginName))
	}
	if plugin.Arch == "" {
		return errors.New(fmt.Sprintf("plugin arch is incorrect: %s", pluginName))
	}
	if plugin.Extension == "" {
		return errors.New(fmt.Sprintf("plugin extension is incorrect: %s", pluginName))
	}
	return nil
}

func (inst *Installer) GetPluginsStorePluginFile(plugin interfaces.Plugin) (pluginsPath string, err error) {
	plugins, err := inst.GetPluginsStorePlugins()
	if err != nil {
		return "", err
	}
	var matchName bool
	var matchArch bool
	for _, plg := range plugins {
		if plg.Name == plugin.Name {
			matchName = true
			if plg.Arch == plugin.Arch {
				matchArch = true
				if plg.Version == plugin.Version {
					return inst.GetPluginsStoreWithFile(plg.ZipName), nil
				}
			}
		}
	}
	if !matchName {
		return "", errors.New(fmt.Sprintf("failed to find plugin name: %s", plugin.Name))
	}
	if !matchArch {
		return "", errors.New(fmt.Sprintf("failed to find plugin arch: %s", plugin.Arch))
	}
	return "", errors.New(fmt.Sprintf("failed to find plugin: %s, version: %s", plugin.Name, plugin.Version))
}

func (inst *Installer) GetPluginsStorePlugins() ([]BuildDetails, error) {
	pluginStore := inst.GetPluginsStorePath()
	files, err := ioutil.ReadDir(pluginStore)
	if err != nil {
		return nil, err
	}
	plugins := make([]BuildDetails, 0)
	for _, file := range files {
		plugins = append(plugins, *inst.GetZipBuildDetails(file.Name()))
	}
	return plugins, err
}
