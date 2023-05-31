package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type ViewTemplateWidgetDatabase interface {
	UpdateViewTemplateWidget(uuid string, body *model.ViewTemplateWidget) (*model.ViewTemplateWidget, error)
	DeleteViewTemplateWidget(uuid string) (bool, error)
}

type ViewTemplateWidgetAPI struct {
	DB ViewTemplateWidgetDatabase
}

func (a *ViewTemplateWidgetAPI) UpdateViewTemplateWidget(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBodyViewTemplateWidget(ctx)
	q, err := a.DB.UpdateViewTemplateWidget(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *ViewTemplateWidgetAPI) DeleteViewTemplateWidget(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteViewTemplateWidget(uuid)
	ResponseHandler(q, err, ctx)
}
