package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type FlowNetworkCloneDatabase interface {
	GetFlowNetworkClones(args Args) ([]*model.FlowNetworkClone, error)
	GetFlowNetworkClone(uuid string, args Args) (*model.FlowNetworkClone, error)
	GetOneFlowNetworkCloneByArgs(args Args) (*model.FlowNetworkClone, error)
	RefreshFlowNetworkClonesConnections() (*bool, error)
}

type FlowNetworkClonesAPI struct {
	DB FlowNetworkCloneDatabase
}

func (a *FlowNetworkClonesAPI) GetFlowNetworkClones(ctx *gin.Context) {
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetFlowNetworkClones(args)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) GetFlowNetworkClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetFlowNetworkClone(uuid, args)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) GetOneFlowNetworkCloneByArgs(ctx *gin.Context) {
	args := buildFlowNetworkCloneArgs(ctx)
	q, err := a.DB.GetOneFlowNetworkCloneByArgs(args)
	reposeHandler(q, err, ctx)
}

func (a *FlowNetworkClonesAPI) RefreshFlowNetworkClonesConnections(ctx *gin.Context) {
	q, err := a.DB.RefreshFlowNetworkClonesConnections()
	reposeHandler(q, err, ctx)
}
