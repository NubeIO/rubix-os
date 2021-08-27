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
}

// The PluginAPI provides handlers for managing plugins.
type PluginAPI struct {
	Notifier Notifier
	Manager  *plugin.Manager
	DB       PluginDatabase
}

// GetPlugins returns all plugins a user has.
// swagger:operation GET /plugin plugin getPlugins
//
// Return all plugins.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//     schema:
//       type: array
//       items:
//         $ref: "#/definitions/PluginConf"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
//   500:
//     description: Internal Server Error
//     schema:
//         $ref: "#/definitions/Error"
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
	ctx.JSON(200, result)
}

// EnablePlugin enables a plugin.
// swagger:operation POST /plugin/{uuid}/enable plugin enablePlugin
//
// Enable a plugin.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// parameters:
// - name: uuid
//   in: path
//   description: the plugin uuid
//   required: true
//   type: integer
//   format: int64
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
//   500:
//     description: Internal Server Error
//     schema:
//         $ref: "#/definitions/Error"
func (c *PluginAPI) EnablePlugin(ctx *gin.Context) {
	uuid := resolveID(ctx)
	conf, err := c.DB.GetPluginConfByID(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil || !isPluginOwner(ctx, conf) {
		ctx.AbortWithError(404, errors.New("unknown plugin"))
		return
	}
	_, err = c.Manager.Instance(uuid)
	if err != nil {
		ctx.AbortWithError(404, errors.New("plugin instance not found"))
		return
	}
	if err := c.Manager.SetPluginEnabled(uuid, true); err == plugin.ErrAlreadyEnabledOrDisabled {
		ctx.AbortWithError(400, err)
	} else if err != nil {
		ctx.AbortWithError(500, err)
	}
	
	
	//withUUID(ctx, "uuid", func(uuid string) {
	//	conf, err := c.DB.GetPluginConfByID(uuid)
	//	if success := successOrAbort(ctx, 500, err); !success {
	//		return
	//	}
	//	if conf == nil || !isPluginOwner(ctx, conf) {
	//		ctx.AbortWithError(404, errors.New("unknown plugin"))
	//		return
	//	}
	//	_, err = c.Manager.Instance(uuid)
	//	if err != nil {
	//		ctx.AbortWithError(404, errors.New("plugin instance not found"))
	//		return
	//	}
	//	if err := c.Manager.SetPluginEnabled(uuid, true); err == plugin.ErrAlreadyEnabledOrDisabled {
	//		ctx.AbortWithError(400, err)
	//	} else if err != nil {
	//		ctx.AbortWithError(500, err)
	//	}
	//})
}

// DisablePlugin disables a plugin.
// swagger:operation POST /plugin/{uuid}/disable plugin disablePlugin
//
// Disable a plugin.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// parameters:
// - name: uuid
//   in: path
//   description: the plugin uuid
//   required: true
//   type: integer
//   format: int64
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
//   500:
//     description: Internal Server Error
//     schema:
//         $ref: "#/definitions/Error"



func (c *PluginAPI) DisablePlugin(ctx *gin.Context) {
	uuid := resolveID(ctx)
	conf, err := c.DB.GetPluginConfByID(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil || !isPluginOwner(ctx, conf) {
		ctx.AbortWithError(404, errors.New("unknown plugin"))
		return
	}
	_, err = c.Manager.Instance(uuid)
	if err != nil {
		ctx.AbortWithError(404, errors.New("plugin instance not found"))
		return
	}
	if err := c.Manager.SetPluginEnabled(uuid, false); err == plugin.ErrAlreadyEnabledOrDisabled {
		ctx.AbortWithError(400, err)
	} else if err != nil {
		ctx.AbortWithError(500, err)
	}

}

// GetDisplay get display info for Displayer plugin.
// swagger:operation GET /plugin/{uuid}/display plugin getPluginDisplay
//
// Get display info for a Displayer plugin.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// parameters:
// - name: uuid
//   in: path
//   description: the plugin uuid
//   required: true
//   type: integer
//   format: int64
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//     schema:
//       type: string
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
//   500:
//     description: Internal Server Error
//     schema:
//         $ref: "#/definitions/Error"
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
// swagger:operation GET /plugin/{uuid}/config plugin getPluginConfig
//
// Get YAML configuration for Configurer plugin.
//
// ---
// consumes: [application/json]
// produces: [application/x-yaml]
// parameters:
// - name: uuid
//   in: path
//   description: the plugin uuid
//   required: true
//   type: integer
//   format: int64
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//     schema:
//         type: object
//         description: plugin configuration
//   400:
//     description: Bad Request
//     schema:
//         $ref: "#/definitions/Error"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
//   500:
//     description: Internal Server Error
//     schema:
//         $ref: "#/definitions/Error"
func (c *PluginAPI) GetConfig(ctx *gin.Context) {
	uuid := resolveID(ctx)
	fmt.Println(uuid)
	fmt.Println(uuid)
	fmt.Println(uuid)
	fmt.Println(uuid)
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
// swagger:operation POST /plugin/{uuid}/config plugin updatePluginConfig
//
// Update YAML configuration for Configurer plugin.
//
// ---
// consumes: [application/x-yaml]
// produces: [application/json]
// parameters:
// - name: uuid
//   in: path
//   description: the plugin uuid
//   required: true
//   type: integer
//   format: int64
// security: [clientTokenHeader: [], clientTokenQuery: [], basicAuth: []]
// responses:
//   200:
//     description: Ok
//   400:
//     description: Bad Request
//     schema:
//         $ref: "#/definitions/Error"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
//   500:
//     description: Internal Server Error
//     schema:
//         $ref: "#/definitions/Error"
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
