package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

/*
Stream
*/


// The GatewayDatabase interface for encapsulating database access.
type GatewayDatabase interface {
	GetStreamGateway(uuid string) (*model.Stream, error)
	GetStreamGateways(withChildren bool) ([]*model.Stream, error)
	CreateStreamGateway(body *model.Stream) error
	UpdateStreamGateway(uuid string, body *model.Stream) (*model.Stream, error)
	DeleteStreamGateway(uuid string) (bool, error)
}

type GatewayAPI struct {
	DB GatewayDatabase
}


func (j *GatewayAPI) GetStreamGateway(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetStreamGateway(uuid)
	reposeHandler(q, err, ctx)
}


func (j *GatewayAPI) GetStreamGateways(ctx *gin.Context) {
	withChildren, _ := withChildrenArgs(ctx)
	q, err := j.DB.GetStreamGateways(withChildren)
	reposeHandler(q, err, ctx)

}

func (j *GatewayAPI) CreateStreamGateway(ctx *gin.Context) {
	body, _ := getBODYGateway(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	err = j.DB.CreateStreamGateway(body)
	if err != nil {
		reposeHandlerError(err, ctx)
	} else {
		reposeHandler(body, err, ctx)
	}

}


func (j *GatewayAPI) UpdateGateway(ctx *gin.Context) {
	body, _ := getBODYGateway(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateStreamGateway(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *GatewayAPI) DeleteStreamGateway(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteStreamGateway(uuid)
	reposeHandler(q, err, ctx)
}
