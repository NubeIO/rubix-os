package module

import (
	"fmt"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/module/shared"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
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

	_ = module.Init(&dbHelper{})
	urlPrefix, err := module.GetUrlPrefix()
	if err != nil {
		log.Error(err)
	} else if urlPrefix == "" {
		log.Errorf("url prefix is empty for module %s", path)
	} else {
		clients[urlPrefix] = client
		modules[urlPrefix] = &module
		mux.Any(fmt.Sprintf("/%s/*proxyPath", urlPrefix), ProxyModule)
	}
	return nil
}

func ProxyModule(c *gin.Context) {
	method := c.Request.Method
	proxyPath := c.Param("proxyPath")
	fullPath := c.FullPath()
	fullPathParts := strings.Split(fullPath, "/")
	if len(fullPathParts) < 4 {
		log.Error("Error on module framework")
		c.JSON(http.StatusInternalServerError, interfaces.Message{Message: "error on module framework!"})
	}
	urlPrefix := fullPathParts[3]
	module := modules[urlPrefix]
	var res []byte
	var err error
	var status int
	if method == "GET" {
		status = http.StatusOK
		res, err = (*module).Get(proxyPath)
	} else if method == "POST" {
		status = http.StatusCreated
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		res, err = (*module).Post(proxyPath, requestBody)
	} else if method == "PUT" {
		status = http.StatusOK
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		res, err = (*module).Put(proxyPath, requestBody)
	} else if method == "PATCH" {
		status = http.StatusOK
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, interfaces.Message{Message: err.Error()})
			return
		}
		res, err = (*module).Patch(proxyPath, requestBody)
	} else if method == "DELETE" {
		status = http.StatusNoContent
		res, err = (*module).Delete(proxyPath)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, interfaces.Message{Message: err.Error()})
		return
	}
	c.Data(status, "application/json", res)
}
