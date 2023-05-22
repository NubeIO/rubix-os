package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type HostTagDatabase interface {
	UpdateHostTags(uuid string, body []*model.HostTag) ([]*model.HostTag, error)
}

type HostTagAPI struct {
	DB HostTagDatabase
}

func (a *HostTagAPI) UpdateHostTags(ctx *gin.Context) {
	body, _ := getBodyHostTags(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateHostTags(uuid, body)
	ResponseHandler(q, err, ctx)
}
