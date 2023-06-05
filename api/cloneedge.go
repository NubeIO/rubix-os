package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type CloneEdgeDatabase interface {
	ResolveHost(uuid string, name string) (*model.Host, error)
	CloneEdge(host *model.Host) error
}
type CloneEdgeApi struct {
	DB CloneEdgeDatabase
}

func (a *CloneEdgeApi) CloneEdge(ctx *gin.Context) {
	matchHostUUIDName(ctx)
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	err = a.DB.CloneEdge(host)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(model.Message{Message: "cloned edge successfully"}, nil, ctx)
}
