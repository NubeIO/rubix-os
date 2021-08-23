package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

/*
Gateway
*/


// The GatewayDatabase interface for encapsulating database access.
type GatewayDatabase interface {
	GetGateway(uuid string) (*model.Gateway, error)
	GetGateways(withChildren bool) ([]model.Gateway, error)
	CreateGateway(body *model.Gateway) error
	UpdateGateway(uuid string, body *model.Gateway) (*model.Gateway, error)
	DeleteGateway(uuid string) (bool, error)
}

type GatewayAPI struct {
	DB GatewayDatabase
}


func (j *GatewayAPI) GetGateway(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetGateway(uuid)
	reposeHandler(q, err, ctx)
}


func (j *GatewayAPI) GetGateways(ctx *gin.Context) {
	withChildren, _ := withChildrenArgs(ctx)
	q, err := j.DB.GetGateways(withChildren)
	reposeHandler(q, err, ctx)

}

func (j *GatewayAPI) CreateGateway(ctx *gin.Context) {
	body, _ := getBODYGateway(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	err = j.DB.CreateGateway(body)
	if err != nil {
		reposeHandlerError(err, ctx)
	} else {
		reposeHandler(body, err, ctx)
	}

}


func (j *GatewayAPI) UpdateGateway(ctx *gin.Context) {
	body, _ := getBODYGateway(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateGateway(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *GatewayAPI) DeleteGateway(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteGateway(uuid)
	reposeHandler(q, err, ctx)
}
