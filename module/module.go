package module

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/database"
	"github.com/NubeIO/rubix-os/module/shared"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

var clients = map[string]*plugin.Client{}
var modules = map[string]shared.Module{}

func ReLoadModulesWithDir(dir string) (map[string]shared.Module, error) {
	var failedModules []string
	UninstallModules()
	if len(failedModules) > 0 {
		return nil, fmt.Errorf("modules [%v] uninstall failed, please retry loading all module after processing", strings.Join(failedModules, ", "))
	}
	return modules, LoadModuleWithLocalDir(dir)
}

func UninstallModules() {
	for _, client := range clients {
		client.Kill()
	}

	var current []string
	for s := range clients {
		current = append(current, s)
	}
	log.Warningf("uninstall all modules, current working modules: %v", strings.Join(current, ";"))
	clients = map[string]*plugin.Client{}
	modules = map[string]shared.Module{}
}

func LoadModuleWithLocalDir(dir string) error {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range fs {
		err = LoadModuleWithLocal(path.Join(dir, f.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

var NameOfModule = "nube-module"

func LoadModuleWithLocal(path string) error {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			NameOfModule: &shared.NubeModule{},
		},
		Cmd:              exec.Command(path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	raw, err := rpcClient.Dispense(NameOfModule)
	module := raw.(shared.Module)

	moduleName := getModuleName(path)
	_ = module.Init(&dbHelper{}, moduleName)
	pluginConf, err := createPluginConf(module, moduleName)
	if err != nil {
		log.Error(err)
	}
	if pluginConf.Enabled {
		log.Infof("enabling module %s", moduleName)
		if err = module.Enable(); err != nil {
			log.Warningf("failed to enable module %s, err: %s", moduleName, err)
		}
	}
	clients[moduleName] = client
	modules[moduleName] = module
	return nil
}

func createPluginConf(module shared.Module, moduleName string) (*model.PluginConf, error) {
	info, err := module.GetInfo()
	if err != nil {
		return nil, err
	}
	pluginConf, _ := database.GlobalGormDatabase.GetPluginByPath(moduleName)

	if pluginConf == nil {
		pluginConf = &model.PluginConf{
			Name:       info.Name,
			ModulePath: moduleName,
			HasNetwork: info.HasNetwork,
		}
		if err := database.GlobalGormDatabase.CreatePlugin(pluginConf); err != nil {
			return nil, err
		}
	}
	return pluginConf, nil
}

// Module naming convention
// ------------------------
// module-core-xxx (for open module, e.g. module-core-lora)
// module-oem-xxx (for private module, e.g. module-oem-cps)
// module-contrib-xxx (for open module, developed by third party)
// module-contrib-oem-xxx (for private module, developed by third party)
//
// moduleName, modulePath and pluginName are same
func getModuleName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
