package api

import (
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The NetworkDatabase interface for encapsulating database access.
type NetworkDatabase interface {
	GetNetwork(uuid string, withChildren bool, withPoints bool) (*model.Network, error)
	GetNetworks(withChildren bool, withPoints bool) ([]*model.Network, error)
	CreateNetwork(network *model.Network) (*model.Network, error)
	UpdateNetwork(uuid string, body *model.Network) (*model.Network, error)
	DeleteNetwork(uuid string) (bool, error)
	DropNetworks() (bool, error)
}
type NetworksAPI struct {
	DB NetworkDatabase
	Bus eventbus.BusService
}


func (a *NetworksAPI) GetNetworks(ctx *gin.Context) {
	withChildren, withPoints := networkArgs(ctx)
	q, err := a.DB.GetNetworks(withChildren, withPoints)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) GetNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	withChildren, withPoints := networkArgs(ctx)
	q, err := a.DB.GetNetwork(uuid, withChildren, withPoints)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) UpdateNetwork(ctx *gin.Context) {
	body, _ := getBODYNetwork(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateNetwork(uuid, body)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) CreateNetwork(ctx *gin.Context) {
	body, _ := getBODYNetwork(ctx)
	q, err := a.DB.CreateNetwork(body)
	reposeHandler(q, err, ctx)
}

func (a *NetworksAPI) DeleteNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteNetwork(uuid)
	reposeHandler(q, err, ctx)

}

func (a *NetworksAPI) DropNetworks(ctx *gin.Context) {
	q, err := a.DB.DropNetworks()
	reposeHandler(q, err, ctx)

}

