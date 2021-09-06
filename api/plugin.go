package api

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/NubeDev/flow-framework/auth"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin"
	"github.com/NubeDev/flow-framework/plugin/compat"
	"github.com/NubeDev/location"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// The PluginDatabase interface for encapsulating database access.
type PluginDatabase interface {
	GetPluginConfByUser(userid uint) ([]*model.PluginConf, error)
	UpdatePluginConf(p *model.PluginConf) error
	GetPluginConfByID(uuid string) (*model.PluginConf, error)
	GetPlugin(uuid string) (*model.PluginConf, error)
	GetPluginByPath(name string) (*model.PluginConf, error)
}

// The PluginAPI provides handlers for managing plugins.
type PluginAPI struct {
	Notifier Notifier
	Manager  *plugin.Manager
	DB       PluginDatabase
}


func (c *PluginAPI) GetPlugin(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := c.DB.GetPlugin(uuid)
	reposeHandler(q, err, ctx)
}

func (c *PluginAPI) GetPluginByPath(ctx *gin.Context) {
	path := resolvePath(ctx)
	q, err := c.DB.GetPluginByPath(path)
	reposeHandler(q, err, ctx)
}


// GetPlugins returns all plugins a user has.
func (c *PluginAPI) GetPlugins(ctx *gin.Context) {
	userID := auth.GetUserID(ctx)
	plugins, err := c.DB.GetPluginConfByUser(userID)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	result := make([]model.PluginConfExternal, 0)
	for _, conf := range plugins {
		if inst, err := c.Manager.Instance(conf.UUID); err == nil {
			info := c.Manager.PluginInfo(conf.ModulePath)
			result = append(result, model.PluginConfExternal{
				UUID:   	conf.UUID,
				Name:         info.String(),
				Token:        conf.Token,
				ModulePath:   conf.ModulePath,
				Author:       info.Author,
				Website:      info.Website,
				License:      info.License,
				Enabled:      conf.Enabled,
				Capabilities: inst.Supports().Strings(),
			})
		}
	}
	reposeHandler(result, err, ctx)
}


// EnablePluginByName enables a plugin.
func (c *PluginAPI) EnablePluginByName(ctx *gin.Context) {
	//uuid := resolveID(ctx)
	body, err := getBODYPlugin(ctx);if err != nil {
		reposeHandler("error on body", err, ctx)
	}
	conf, err := c.DB.GetPluginByPath(body.ModulePath)
	uuid := conf.UUID
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil || !isPluginOwner(ctx, conf) {
		reposeHandler("unknown plugin", err, ctx)
		return
	}
	_, err = c.Manager.Instance(uuid)
	if err != nil {
		reposeHandler("plugin not found", err, ctx)
		return
	}
	if err := c.Manager.SetPluginEnabled(uuid, true); err == plugin.ErrAlreadyEnabledOrDisabled {
		reposeHandler("err:", err, ctx)
	} else if err != nil {
		reposeHandler("err:", err, ctx)
	}
	reposeHandler("enabled", err, ctx)
}


// GetDisplay get display info for Displayer plugin.
func (c *PluginAPI) GetDisplay(ctx *gin.Context) {
	uuid := resolveID(ctx)
	conf, err := c.DB.GetPluginConfByID(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil || !isPluginOwner(ctx, conf) {
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
	uuid := resolveID(ctx)
	conf, err := c.DB.GetPluginConfByID(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil || !isPluginOwner(ctx, conf) {
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
	uuid := resolveID(ctx)
	conf, err := c.DB.GetPluginConfByID(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil || !isPluginOwner(ctx, conf) {
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
	newconfBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	if err := yaml.Unmarshal(newconfBytes, newConf); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	if err := instance.ValidateAndSetConfig(newConf); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	conf.Config = newconfBytes
	successOrAbort(ctx, 500, c.DB.UpdatePluginConf(conf))
	
}

func isPluginOwner(ctx *gin.Context, conf *model.PluginConf) bool {
	return conf.UserID == auth.GetUserID(ctx)
}

func supportOrAbort(ctx *gin.Context, instance compat.PluginInstance, module compat.Capability) (aborted bool) {
	if compat.HasSupport(instance, module) {
		return false
	}
	ctx.AbortWithError(400, fmt.Errorf("plugin does not support %s", module))
	return true
}
