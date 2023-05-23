package api

import (
	"encoding/json"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
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
	SyncFlowNetworks(args Args) []*interfaces.SyncModel
	SyncFlowNetworkStreams(uuid string, args Args) ([]*interfaces.SyncModel, error)
}
type FlowNetworksAPI struct {
	DB FlowNetworkDatabase
}

func (a *FlowNetworksAPI) GetFlowNetworks(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetFlowNetworks(args)
	if err == nil && args.IsMetadata {
		var flowNetworksMetaData []*interfaces.FlowNetworkMetadata
		out, _ := json.Marshal(q)
		_ = json.Unmarshal(out, &flowNetworksMetaData)
		ResponseHandler(flowNetworksMetaData, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) GetFlowNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetFlowNetwork(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) GetOneFlowNetworkByArgs(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.GetOneFlowNetworkByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) UpdateFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYFlowNetwork(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateFlowNetwork(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) CreateFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYFlowNetwork(ctx)
	q, err := a.DB.CreateFlowNetwork(body)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) DeleteFlowNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteFlowNetwork(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) RefreshFlowNetworksConnections(ctx *gin.Context) {
	q, err := a.DB.RefreshFlowNetworksConnections()
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworksAPI) SyncFlowNetworks(ctx *gin.Context) {
	args := buildFlowNetworkArgs(ctx)
	q := a.DB.SyncFlowNetworks(args)
	ResponseHandler(q, nil, ctx)
}

func (a *FlowNetworksAPI) SyncFlowNetworkStreams(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildFlowNetworkArgs(ctx)
	q, err := a.DB.SyncFlowNetworkStreams(uuid, args)
	ResponseHandler(q, err, ctx)
}
