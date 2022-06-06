package api

import (
	"encoding/json"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type FlowNetworkCloneDatabase interface {
	GetFlowNetworkClones(args Args) ([]*model.FlowNetworkClone, error)
	GetFlowNetworkClone(uuid string, args Args) (*model.FlowNetworkClone, error)
	DeleteFlowNetworkClone(uuid string) (bool, error)
	GetOneFlowNetworkCloneByArgs(args Args) (*model.FlowNetworkClone, error)
	DeleteOneFlowNetworkCloneByArgs(args Args) (bool, error)
	RefreshFlowNetworkClonesConnections() (*bool, error)
	SyncFlowNetworkClones(args Args) ([]*interfaces.SyncModel, error)
	SyncFlowNetworkCloneStreamClones(uuid string, args Args) ([]*interfaces.SyncModel, error)
}

type FlowNetworkClonesAPI struct {
	DB FlowNetworkCloneDatabase
}

func (a *FlowNetworkClonesAPI) GetFlowNetworkClones(ctx *gin.Context) {
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetFlowNetworkClones(args)
	if err == nil && args.IsMetadata {
		var flowNetworkClonesMetaData []*interfaces.FlowNetworkCloneMetadata
		out, _ := json.Marshal(q)
		_ = json.Unmarshal(out, &flowNetworkClonesMetaData)
		ResponseHandler(flowNetworkClonesMetaData, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) GetFlowNetworkClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetFlowNetworkClone(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) DeleteFlowNetworkClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteFlowNetworkClone(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) GetOneFlowNetworkCloneByArgs(ctx *gin.Context) {
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetOneFlowNetworkCloneByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) DeleteOneFlowNetworkCloneByArgs(ctx *gin.Context) {
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.DeleteOneFlowNetworkCloneByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) RefreshFlowNetworkClonesConnections(ctx *gin.Context) {
	q, err := a.DB.RefreshFlowNetworkClonesConnections()
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) SyncFlowNetworkClones(ctx *gin.Context) {
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.SyncFlowNetworkClones(args)
	ResponseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) SyncFlowNetworkCloneStreamClones(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.SyncFlowNetworkCloneStreamClones(uuid, args)
	ResponseHandler(q, err, ctx)
}
