package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type ViewWidgetDatabase interface {
	CreateViewWidget(body *model.ViewWidget) (*model.ViewWidget, error)
	UpdateViewWidget(uuid string, body *model.ViewWidget) (*model.ViewWidget, error)
	DeleteViewWidget(uuid string) (bool, error)
}

type ViewWidgetAPI struct {
	DB ViewWidgetDatabase
}

func (a *ViewWidgetAPI) CreateViewWidget(ctx *gin.Context) {
	body, _ := getBodyViewWidget(ctx)
	q, err := a.DB.CreateViewWidget(body)
	ResponseHandler(q, err, ctx)
}

func (a *ViewWidgetAPI) UpdateViewWidget(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBodyViewWidget(ctx)
	q, err := a.DB.UpdateViewWidget(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *ViewWidgetAPI) DeleteViewWidget(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteViewWidget(uuid)
	ResponseHandler(q, err, ctx)
}
