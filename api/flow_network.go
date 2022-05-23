package api

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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
	RefreshFlowNetworksConnections() (*bool, error)
	SyncFlowNetworks() []*interfaces.SyncModel
	SyncFlowNetworkStreams(uuid string) []*interfaces.SyncModel
}
type FlowNetworksAPI struct {
	DB FlowNetworkDatabase
}

func (a *FlowNetworksAPI) GetFlowNetworks(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetFlowNetworks(args)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) GetFlowNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetFlowNetwork(uuid, args)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) GetOneFlowNetworkByArgs(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetOneFlowNetworkByArgs(args)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) UpdateFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYFlowNetwork(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateFlowNetwork(uuid, body)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) CreateFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYFlowNetwork(ctx)
	q, err := a.DB.CreateFlowNetwork(body)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) DeleteFlowNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteFlowNetwork(uuid)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) RefreshFlowNetworksConnections(ctx *gin.Context) {
	q, err := a.DB.RefreshFlowNetworksConnections()
	responseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) SyncFlowNetworks(ctx *gin.Context) {
	q := a.DB.SyncFlowNetworks()
	responseHandler(q, nil, ctx)
}

func (a *FlowNetworksAPI) SyncFlowNetworkStreams(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q := a.DB.SyncFlowNetworkStreams(uuid)
	responseHandler(q, nil, ctx)
}
