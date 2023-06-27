package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nerrors"
	"github.com/gin-gonic/gin"
)

type TicketCommentDatabase interface {
	GetTicketComment(uuid string) (*model.TicketComment, error)
	CreateTicketComment(body *model.TicketComment) (*model.TicketComment, error)
	UpdateTicketComment(uuid string, body *model.TicketComment) (*model.TicketComment, error)
	DeleteTicketComment(uuid string) (bool, error)
}

type TicketCommentAPI struct {
	DB TicketCommentDatabase
}

func (a *TicketCommentAPI) GetTicketComment(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetTicketComment(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *TicketCommentAPI) CreateTicketComment(ctx *gin.Context) {
	username, err := getAuthorizedOrDefaultUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	body, _ := getBodyTicketComment(ctx)
	body.Owner = username
	q, err := a.DB.CreateTicketComment(body)
	ResponseHandler(q, err, ctx)
}

func (a *TicketCommentAPI) UpdateTicketComment(ctx *gin.Context) {
	username, err := getAuthorizedOrDefaultUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	body, _ := getBodyTicketComment(ctx)
	body.Owner = username
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateTicketComment(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *TicketCommentAPI) DeleteTicketComment(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteTicketComment(uuid)
	ResponseHandler(q, err, ctx)
}
