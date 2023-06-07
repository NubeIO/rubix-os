package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/rubix-os/module/shared"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"

	"github.com/NubeDev/location"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/plugin"
	"github.com/NubeIO/rubix-os/plugin/compat"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// The PluginDatabase interface for encapsulating database access.
type PluginDatabase interface {
	GetPlugins() ([]*model.PluginConf, error)
	GetPluginByPath(path string) (*model.PluginConf, error)
	GetPlugin(uuid string) (*model.PluginConf, error)
	UpdatePluginConf(p *model.PluginConf) error
}

// The PluginAPI provides handlers for managing plugins.
type PluginAPI struct {
	Manager *plugin.Manager
	Modules map[string]shared.Module
	DB      PluginDatabase
}

func (c *PluginAPI) GetPlugin(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	q, err := c.DB.GetPlugin(uuid)
	ResponseHandler(q, err, ctx)
}

func (c *PluginAPI) GetPluginByPath(ctx *gin.Context) {
	path := resolvePath(ctx)
	q, err := c.DB.GetPluginByPath(path)
	ResponseHandler(q, err, ctx)
}

// GetPlugins returns all plugins a user has.
func (c *PluginAPI) GetPlugins(ctx *gin.Context) {
	plugins, err := c.DB.GetPlugins()
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	result := make([]model.PluginConfExternal, 0)
	for _, conf := range plugins {
		if strings.HasPrefix(conf.ModulePath, "module") {
			if module, found := c.Modules[conf.ModulePath]; found {
				info, err := module.GetInfo()
				if err != nil {
					log.Errorf("can't get info details from module %s", conf.ModulePath)
					continue
				}
				result = append(result, model.PluginConfExternal{
					UUID:         conf.UUID,
					Name:         info.Name,
					ModulePath:   conf.ModulePath,
					Author:       info.Author,
					Website:      info.Website,
					License:      info.License,
					Enabled:      conf.Enabled,
					HasNetwork:   info.HasNetwork,
					Capabilities: []string{},
				})
			}
		} else {
			if inst, err := c.Manager.Instance(conf.UUID); err == nil {
				info := c.Manager.PluginInfo(conf.ModulePath)
				result = append(result, model.PluginConfExternal{
					UUID:         conf.UUID,
					Name:         info.String(),
					ModulePath:   conf.ModulePath,
					Author:       info.Author,
					Website:      info.Website,
					License:      info.License,
					Enabled:      conf.Enabled,
					HasNetwork:   info.HasNetwork,
					Capabilities: inst.Supports().Strings(),
				})
			}
		}

	}
	ResponseHandler(result, err, ctx)
}

// buildUUID a way to query a plugin by its name or uuid
func (c *PluginAPI) buildUUID(ctx *gin.Context) string {
	nameOrUUID := resolveID(ctx) // system?by_plugin_name=true is passed in then enable plugin by its name
	uuid := ""
	args := buildPluginArgs(ctx)
	if args.ByPluginName {
		path, err := c.DB.GetPluginByPath(nameOrUUID)
		if err != nil {
			ResponseHandler("err: no plugin with that name was found", nil, ctx)
		} else {
			uuid = path.UUID
		}
	} else {
		uuid = resolveID(ctx)
	}
	if uuid == "" {
		ResponseHandler("err: no valid uuid found", nil, ctx)
	}
	return uuid
}

// EnablePluginByUUID enables a plugin.
func (c *PluginAPI) EnablePluginByUUID(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	body, err := getBODYPlugin(ctx)
	if err != nil {
		ResponseHandler("error on body", err, ctx)
	}
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		ResponseHandler("unknown plugin", err, ctx)
		return
	}
	if strings.HasPrefix(conf.ModulePath, "module") {
		conf, err = c.DB.GetPlugin(conf.UUID)
		if err != nil {
			ResponseHandler(nil, err, ctx)
			return
		}
		if conf.Enabled == body.Enabled {
			ResponseHandler(nil, errors.New("config is already on your state"), ctx)
		}
		module, found := c.Modules[conf.ModulePath]
		if !found {
			errMsg := fmt.Sprintf("not found module %s", conf.ModulePath)
			ResponseHandler(nil, errors.New(errMsg), ctx)
			return
		}
		if body.Enabled {
			err = module.Enable()
			if err != nil {
				ResponseHandler(nil, err, ctx)
				return
			}
		} else {
			err = module.Disable()
			if err != nil {
				ResponseHandler(nil, err, ctx)
				return
			}
		}
		conf.Enabled = body.Enabled
		err = c.DB.UpdatePluginConf(conf)
		if err != nil {
			ResponseHandler(nil, err, ctx)
			return
		}
	} else {
		_, err = c.Manager.Instance(uuid)
		if err != nil {
			ResponseHandler("plugin not found", err, ctx)
			return
		}
		if err = c.Manager.SetPluginEnabled(uuid, body.Enabled); err == plugin.ErrAlreadyEnabledOrDisabled {
			ResponseHandler(nil, err, ctx)
			return
		} else if err != nil {
			ResponseHandler(nil, err, ctx)
			return
		}
	}
	if body.Enabled {
		ResponseHandler(map[string]string{"state": "enabled"}, err, ctx)
	} else {
		ResponseHandler(map[string]string{"state": "disabled"}, err, ctx)
	}
}

// RestartPlugin enables a plugin.
func (c *PluginAPI) RestartPlugin(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		ResponseHandler("unknown plugin", err, ctx)
		return
	}
	_, err = c.Manager.Instance(uuid)
	if err != nil {
		ResponseHandler("plugin not found", err, ctx)
		return
	}
	if res, err := c.Manager.RestartPlugin(uuid); err == plugin.ErrAlreadyEnabledOrDisabled {
		ResponseHandler(res, err, ctx)
	} else if err != nil {
		ResponseHandler(res, nil, ctx)
	}
	ResponseHandler("plugin restart ok", err, ctx)
}

// RestartPluginByName restart a plugin.
func (c *PluginAPI) RestartPluginByName(ctx *gin.Context) {
	name := ctx.Param("name")
	plugins, err := c.DB.GetPlugins()
	if err != nil {
		ResponseHandler("plugin", err, ctx)
		return
	}
	for _, conf := range plugins {
		if conf.Name == name {
			uuid := conf.UUID
			if success := successOrAbort(ctx, 500, err); !success {
				return
			}
			if conf == nil {
				ResponseHandler("unknown plugin", err, ctx)
				return
			}
			_, err = c.Manager.Instance(uuid)
			if err != nil {
				ResponseHandler("plugin not found", err, ctx)
				return
			}
			if res, err := c.Manager.RestartPlugin(uuid); err == plugin.ErrAlreadyEnabledOrDisabled {
				ResponseHandler(res, err, ctx)
			} else if err != nil {
				ResponseHandler(res, nil, ctx)
			}
			ResponseHandler("plugin restart ok", err, ctx)
			return
		}
	}
	ResponseHandler(fmt.Sprintf("plugin not found with that name:%s", name), nil, ctx)
	return

}

// GetDisplay get display info for Displayer plugin.
func (c *PluginAPI) GetDisplay(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		ctx.AbortWithError(404, errors.New("unknown plugin"))
		return
	}
	instance, err := c.Manager.Instance(uuid)
	if err != nil {
		ctx.AbortWithError(404, errors.New("plugin instance not found"))
		return
	}
	ctx.JSON(200, instance.GetDisplay(location.Get(ctx)))

}

// GetConfig returns Configurer plugin configuration in YAML format.
func (c *PluginAPI) GetConfig(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		ctx.AbortWithError(404, errors.New("unknown plugin"))
		return
	}
	instance, err := c.Manager.Instance(uuid)
	if err != nil {
		ctx.AbortWithError(404, errors.New("plugin instance not found"))
		return
	}

	if aborted := supportOrAbort(ctx, instance, compat.Configurer); aborted {
		return
	}

	ctx.Header("content-type", "application/x-yaml")
	ctx.Writer.Write(conf.Config)

}

// UpdateConfig updates Configurer plugin configuration in YAML format.
func (c *PluginAPI) UpdateConfig(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		ctx.AbortWithError(404, errors.New("unknown plugin"))
		return
	}
	instance, err := c.Manager.Instance(uuid)
	if err != nil {
		ctx.AbortWithError(404, errors.New("plugin instance not found"))
		return
	}

	if aborted := supportOrAbort(ctx, instance, compat.Configurer); aborted {
		return
	}

	newConf := instance.DefaultConfig()
	_ = yaml.Unmarshal(conf.Config, newConf)
	newConfBytes, err := ioutil.ReadAll(ctx.Request.Body)
	var jsonConf map[string]string
	err = json.Unmarshal(newConfBytes, &jsonConf)
	if err != nil {
		ctx.AbortWithError(400, errors.New("invalid data"))
		return
	}
	if err := yaml.Unmarshal([]byte(jsonConf["data"]), newConf); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	if err := instance.ValidateAndSetConfig(newConf); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	config, _ := yaml.Marshal(instance.GetConfig())
	conf.Config = config
	successOrAbort(ctx, 500, c.DB.UpdatePluginConf(conf))
}

func supportOrAbort(ctx *gin.Context, instance compat.PluginInstance, module compat.Capability) (aborted bool) {
	if compat.HasSupport(instance, module) {
		return false
	}
	ctx.AbortWithError(400, fmt.Errorf("plugin does not support %s", module))
	return true
}
