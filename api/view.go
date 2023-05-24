package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type ViewDatabase interface {
	GetView(uuid string) (*model.View, error)
	CreateView(body *model.View) (*model.View, error)
	UpdateView(uuid string, body *model.View) (*model.View, error)
	DeleteView(uuid string) (bool, error)
}

type ViewAPI struct {
	DB ViewDatabase
}

func (a *ViewAPI) GetView(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetView(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *ViewAPI) CreateView(ctx *gin.Context) {
	body, _ := getBodyView(ctx)
	q, err := a.DB.CreateView(body)
	ResponseHandler(q, err, ctx)
}

func (a *ViewAPI) UpdateView(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBodyView(ctx)
	q, err := a.DB.UpdateView(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *ViewAPI) DeleteView(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteView(uuid)
	ResponseHandler(q, err, ctx)
}
