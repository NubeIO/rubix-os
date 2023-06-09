package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/plugin"
	"github.com/gin-gonic/gin"
	"strconv"
)

type NetworkDatabase interface {
	GetNetworkByName(name string, args argspkg.Args) (*model.Network, error)
	GetNetworkByPluginName(name string, args argspkg.Args) (*model.Network, error)
	GetNetworksByPluginName(name string, args argspkg.Args) ([]*model.Network, error)
	GetNetworks(args argspkg.Args) ([]*model.Network, error)
	GetNetwork(uuid string, args argspkg.Args) (*model.Network, error)
	CreateNetwork(network *model.Network) (*model.Network, error)
	UpdateNetwork(uuid string, body *model.Network) (*model.Network, error)
	DeleteNetwork(uuid string) (bool, error)
	DeleteOneNetworkByArgs(args argspkg.Args) (bool, error)
	DeleteNetworkByName(name string, args argspkg.Args) (bool, error)

	CreateNetworkPlugin(network *model.Network) (*model.Network, error)
	UpdateNetworkPlugin(uuid string, body *model.Network) (*model.Network, error)
	DeleteNetworkPlugin(uuid string) (bool, error)

	CreateNetworkMetaTags(networkUUID string, networkMetaTags []*model.NetworkMetaTag) ([]*model.NetworkMetaTag, error)
}

type NetworksAPI struct {
	DB     NetworkDatabase
	Bus    eventbus.BusService
	Plugin *plugin.Manager
}

func (a *NetworksAPI) GetNetworkByName(ctx *gin.Context) {
	name := resolveName(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworkByName(name, args) // TODO fix this need to add in like "serial"
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetworkByPluginName(ctx *gin.Context) {
	name := resolveName(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworkByPluginName(name, args) // TODO fix this need to add in like "serial"
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetworksByPluginName(ctx *gin.Context) {
	name := resolveName(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworksByPluginName(name, args)
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetworks(ctx *gin.Context) {
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworks(args)
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetwork(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) CreateNetwork(ctx *gin.Context) {
	body, _ := getBODYNetwork(ctx)
	restart, _ := strconv.ParseBool(ctx.Query("restart_plugin"))
	q, err := a.DB.CreateNetworkPlugin(body)
	if err != nil {
		ResponseHandler(q, err, ctx)
		return
	}
	if restart {
		if q.PluginConfId != "" {
			err = a.Plugin.RestartPlugin(q.PluginConfId)
			if err != nil {
				ResponseHandler(nil, err, ctx)
				return
			}
		}
	}
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) UpdateNetwork(ctx *gin.Context) {
	body, _ := getBODYNetwork(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateNetworkPlugin(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) DeleteNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteNetworkPlugin(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) DeleteOneNetworkByArgs(ctx *gin.Context) {
	args := buildNetworkArgs(ctx)
	q, err := a.DB.DeleteOneNetworkByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) DeleteNetworkByName(ctx *gin.Context) {
	name := resolveName(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.DeleteNetworkByName(name, args)
	ResponseHandler(q, err, ctx)
}

func (a *NetworksAPI) CreateNetworkMetaTags(ctx *gin.Context) {
	networkUUID := resolveID(ctx)
	body, _ := getBodyBulkNetworkMetaTags(ctx)
	q, err := a.DB.CreateNetworkMetaTags(networkUUID, body)
	if err != nil {
		ResponseHandler(q, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}
