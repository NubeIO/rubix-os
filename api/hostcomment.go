package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type HostCommentDatabase interface {
	CreateHostComment(body *model.HostComment) (*model.HostComment, error)
	UpdateHostComment(uuid string, body *model.HostComment) (*model.HostComment, error)
	DeleteHostComment(uuid string) (*interfaces.Message, error)
}

type HostCommentAPI struct {
	DB HostCommentDatabase
}

func (a *HostCommentAPI) CreateHostComment(ctx *gin.Context) {
	body, _ := getBodyHostComment(ctx)
	q, err := a.DB.CreateHostComment(body)
	ResponseHandler(q, err, ctx)
}

func (a *HostCommentAPI) UpdateHostComment(ctx *gin.Context) {
	body, _ := getBodyHostComment(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateHostComment(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *HostCommentAPI) DeleteHostComment(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteHostComment(uuid)
	ResponseHandler(q, err, ctx)
}
