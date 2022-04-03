package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type FlowNetworkCloneDatabase interface {
	GetFlowNetworkClones(args Args) ([]*model.FlowNetworkClone, error)
	GetFlowNetworkClone(uuid string, args Args) (*model.FlowNetworkClone, error)
	DeleteFlowNetworkClone(uuid string) error
	GetOneFlowNetworkCloneByArgs(args Args) (*model.FlowNetworkClone, error)
	RefreshFlowNetworkClonesConnections() (*bool, error)
}

type FlowNetworkClonesAPI struct {
	DB FlowNetworkCloneDatabase
}

func (a *FlowNetworkClonesAPI) GetFlowNetworkClones(ctx *gin.Context) {
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetFlowNetworkClones(args)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) GetFlowNetworkClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetFlowNetworkClone(uuid, args)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) DeleteFlowNetworkClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	err := a.DB.DeleteFlowNetworkClone(uuid)
	responseHandler(nil, err, ctx)
}

func (a *FlowNetworkClonesAPI) GetOneFlowNetworkCloneByArgs(ctx *gin.Context) {
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetOneFlowNetworkCloneByArgs(args)
	responseHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) RefreshFlowNetworkClonesConnections(ctx *gin.Context) {
	q, err := a.DB.RefreshFlowNetworkClonesConnections()
	responseHandler(q, err, ctx)
}
