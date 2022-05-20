package api

import (
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type NetworkDatabase interface {
	GetNetworkByName(name string, args Args) (*model.Network, error)
	GetNetworksByName(name string, args Args) ([]*model.Network, error)
	GetNetworkByPluginName(name string, args Args) (*model.Network, error)
	GetNetworksByPluginName(name string, args Args) ([]*model.Network, error)
	GetNetworks(args Args) ([]*model.Network, error)
	GetNetwork(uuid string, args Args) (*model.Network, error)
	CreateNetwork(network *model.Network, fromPlugin bool) (*model.Network, error)
	UpdateNetwork(uuid string, body *model.Network, fromPlugin bool) (*model.Network, error)
	DeleteNetwork(uuid string) (bool, error)

	CreateNetworkPlugin(network *model.Network) (*model.Network, error)
	UpdateNetworkPlugin(uuid string, body *model.Network) (*model.Network, error)
	DeleteNetworkPlugin(uuid string) (bool, error)
}
type NetworksAPI struct {
	DB  NetworkDatabase
	Bus eventbus.BusService
}

func (a *NetworksAPI) GetNetworkByName(ctx *gin.Context) {
	name := resolveName(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworkByName(name, args) //TODO fix this need to add in like "serial"
	responseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetworksByName(ctx *gin.Context) {
	name := resolveName(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworksByName(name, args)
	responseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetworkByPluginName(ctx *gin.Context) {
	name := resolveName(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworkByPluginName(name, args) //TODO fix this need to add in like "serial"
	responseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetworksByPluginName(ctx *gin.Context) {
	name := resolveName(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworksByPluginName(name, args)
	responseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetworks(ctx *gin.Context) {
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetworks(args)
	responseHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildNetworkArgs(ctx)
	q, err := a.DB.GetNetwork(uuid, args)
	responseHandler(q, err, ctx)
}

func (a *NetworksAPI) CreateNetwork(ctx *gin.Context) {
	body, _ := getBODYNetwork(ctx)
	q, err := a.DB.CreateNetworkPlugin(body)
	responseHandler(q, err, ctx)
}

func (a *NetworksAPI) UpdateNetwork(ctx *gin.Context) {
	body, _ := getBODYNetwork(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateNetworkPlugin(uuid, body)
	responseHandler(q, err, ctx)
}

func (a *NetworksAPI) DeleteNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteNetworkPlugin(uuid)
	responseHandler(q, err, ctx)
}
