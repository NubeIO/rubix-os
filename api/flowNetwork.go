package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The FlowNetworkDatabase interface for encapsulating database access.
type FlowNetworkDatabase interface {
	GetFlowNetworks(args Args) ([]*model.FlowNetwork, error)
	GetFlowNetwork(uuid string, args Args) (*model.FlowNetwork, error)
	GetOneFlowNetworkByArgs(args Args) (*model.FlowNetwork, error)
	CreateFlowNetwork(network *model.FlowNetwork) (*model.FlowNetwork, error)
	UpdateFlowNetwork(uuid string, body *model.FlowNetwork) (*model.FlowNetwork, error)
	DeleteFlowNetwork(uuid string) (bool, error)
	DropFlowNetworks() (bool, error)
	RefreshFlowNetworksConnections() (*bool, error)
}
type FlowNetworksAPI struct {
	DB FlowNetworkDatabase
}

func (a *FlowNetworksAPI) GetFlowNetworks(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetFlowNetworks(args)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) GetFlowNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetFlowNetwork(uuid, args)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) GetOneFlowNetworkByArgs(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetOneFlowNetworkByArgs(args)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) UpdateFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYFlowNetwork(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateFlowNetwork(uuid, body)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) CreateFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYFlowNetwork(ctx)
	q, err := a.DB.CreateFlowNetwork(body)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) DeleteFlowNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteFlowNetwork(uuid)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) DropFlowNetworks(ctx *gin.Context) {
	q, err := a.DB.DropFlowNetworks()
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) RefreshFlowNetworksConnections(ctx *gin.Context) {
	q, err := a.DB.RefreshFlowNetworksConnections()
	reposeHandler(q, err, ctx)
}
