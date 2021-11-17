package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type LocalStorageFlowNetworkDatabase interface {
	GetLocalStorageFlowNetwork() (*model.LocalStorageFlowNetwork, error)
	UpdateLocalStorageFlowNetwork(network *model.LocalStorageFlowNetwork) (*model.LocalStorageFlowNetwork, error)
	RefreshLocalStorageFlowToken() (*bool, error)
}

type LocalStorageFlowNetworkAPI struct {
	DB LocalStorageFlowNetworkDatabase
}

func (a *LocalStorageFlowNetworkAPI) GetLocalStorageFlowNetwork(ctx *gin.Context) {
	q, err := a.DB.GetLocalStorageFlowNetwork()
	reposeHandler(q, err, ctx)
}

func (a *LocalStorageFlowNetworkAPI) UpdateLocalStorageFlowNetwork(ctx *gin.Context) {
	body, _ := getBodyLocalStorageFlowNetwork(ctx)
	q, err := a.DB.UpdateLocalStorageFlowNetwork(body)
	reposeHandler(q, err, ctx)
}

func (a *LocalStorageFlowNetworkAPI) RefreshLocalStorageFlowToken(ctx *gin.Context) {
	q, err := a.DB.RefreshLocalStorageFlowToken()
	reposeHandler(q, err, ctx)
}
