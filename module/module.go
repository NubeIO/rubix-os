package module

import (
	"fmt"
	"github.com/NubeIO/flow-framework/database"
	"github.com/NubeIO/flow-framework/module/shared"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

var clients = map[string]*plugin.Client{}
var modules = map[string]*shared.Module{}

func ReLoadModulesWithDir(dir string, mux *gin.RouterGroup) error {
	var failedModules []string
	UninstallModules()
	if len(failedModules) > 0 {
		return fmt.Errorf("modules [%v] uninstall failed, please retry loading all module after processing", strings.Join(failedModules, ", "))
	}
	return LoadModuleWithLocalDir(dir, mux)
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
	modules = map[string]*shared.Module{}
}

func LoadModuleWithLocalDir(dir string, mux *gin.RouterGroup) error {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range fs {
		err = LoadModuleWithLocal(path.Join(dir, f.Name()), mux)
		if err != nil {
			return err
		}
	}
	return nil
}

var NameOfModule = "nube-module"

func LoadModuleWithLocal(path string, mux *gin.RouterGroup) error {
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
	_, err = createPluginConf(module, moduleName)
	if err != nil {
		log.Error(err)
	}
	urlPrefix, err := module.GetUrlPrefix()
	if err != nil {
		log.Error(err)
	} else if urlPrefix == nil {
		log.Errorf("url prefix is empty for module %s", path)
	} else {
		clients[*urlPrefix] = client
		modules[*urlPrefix] = &module
		mux.Any(fmt.Sprintf("/%s/*proxyPath", *urlPrefix), ProxyModule)
	}
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

// moduleName, modulePath and pluginName are same
func getModuleName(path string) string {
	parts := strings.Split(path, "/")
	module := parts[len(parts)-1]
	return fmt.Sprintf("%s-module", module)
}
