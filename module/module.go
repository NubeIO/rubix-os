package module

import (
	"fmt"
	"github.com/NubeIO/flow-framework/module/shared"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

var modules = map[string] /*action name@action version*/ *plugin.Client{}

func ReLoadModulesWithDir(dir string) error {
	var failedModules []string
	UninstallModules()
	if len(failedModules) > 0 {
		return fmt.Errorf("modules [%v] uninstall failed, please retry loading all module after processing", strings.Join(failedModules, ", "))
	}
	return LoadModuleWithLocalDir(dir)
}

func UninstallModules() {
	for _, client := range modules {
		client.Kill()
	}

	var current []string
	for s := range modules {
		current = append(current, s)
	}
	log.Warningf("uninstall all modules, current working modules: %v", strings.Join(current, ";"))
	modules = map[string]*plugin.Client{}
}

func UninstallModule(actionName string, actionVersion float32) {
	key := fmt.Sprintf("%v@%v", actionName, actionVersion)
	log.Infof("uninstall the module for action %v", key)
	modules[key].Kill()
	delete(modules, key)
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

var NameOfModule = "action"

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
	counter := raw.(shared.Module)

	_ = counter.Put("test", 34, &addHelper{})
	a, err := counter.Get("test")
	fmt.Println("c>>>>>>>>>>", strconv.Itoa(int(a)), err)

	a, err = counter.Get("test")
	fmt.Println("c>>>>>>>>>>", strconv.Itoa(int(a)), err)

	a, err = counter.Get("test")
	fmt.Println("c>>>>>>>>>>", strconv.Itoa(int(a)), err)
	return nil
}

type addHelper struct{}

func (*addHelper) Sum(a, b int64) (int64, error) {
	return a + b, nil
}
