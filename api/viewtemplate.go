package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/gin-gonic/gin"
)

type ViewTemplateDatabase interface {
	GetViewTemplates(args argspkg.Args) ([]*model.ViewTemplate, error)
	GetViewTemplate(uuid string, args argspkg.Args) (*model.ViewTemplate, error)
	CreateViewTemplate(body *model.ViewTemplate) (*model.ViewTemplate, error)
	UpdateViewTemplate(uuid string, body *model.ViewTemplate) (*model.ViewTemplate, error)
	DeleteViewTemplate(uuid string) (bool, error)
}

type ViewTemplateAPI struct {
	DB ViewTemplateDatabase
}

func (a *ViewTemplateAPI) GetViewTemplates(ctx *gin.Context) {
	args := buildViewTemplateArgs(ctx)
	q, err := a.DB.GetViewTemplates(args)
	ResponseHandler(q, err, ctx)
}

func (a *ViewTemplateAPI) GetViewTemplate(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildViewTemplateArgs(ctx)
	q, err := a.DB.GetViewTemplate(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *ViewTemplateAPI) CreateViewTemplate(ctx *gin.Context) {
	body, _ := getBodyViewTemplate(ctx)
	q, err := a.DB.CreateViewTemplate(body)
	ResponseHandler(q, err, ctx)
}

func (a *ViewTemplateAPI) UpdateViewTemplate(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBodyViewTemplate(ctx)
	q, err := a.DB.UpdateViewTemplate(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *ViewTemplateAPI) DeleteViewTemplate(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteViewTemplate(uuid)
	ResponseHandler(q, err, ctx)
}
