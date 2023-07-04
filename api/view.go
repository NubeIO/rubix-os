package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
	"github.com/gin-gonic/gin"
)

type ViewDatabase interface {
	GetViews(args args.Args) ([]*model.View, error)
	GetView(uuid string, args args.Args) (*model.View, error)
	CreateView(body *model.View) (*model.View, error)
	UpdateView(uuid string, body *model.View) (*model.View, error)
	DeleteView(uuid string) (bool, error)

	GenerateViewTemplate(uuid string, templateUUID string) (bool, error)
	AssignViewTemplate(uuid string, viewTemplateUUID string, hostUUID string) (bool, error)
}

type ViewAPI struct {
	DB ViewDatabase
}

func (a *ViewAPI) GetViews(ctx *gin.Context) {
	args := buildViewArgs(ctx)
	q, err := a.DB.GetViews(args)
	ResponseHandler(q, err, ctx)
}

func (a *ViewAPI) GetView(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildViewArgs(ctx)
	q, err := a.DB.GetView(uuid, args)
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

func (a *ViewAPI) GenerateViewTemplate(ctx *gin.Context) {
	body, _ := getBodyGenerateViewTemplate(ctx)
	q, err := a.DB.GenerateViewTemplate(body.ViewUUID, body.Name)
	ResponseHandler(q, err, ctx)
}

func (a *ViewAPI) AssignViewTemplate(ctx *gin.Context) {
	body, _ := getBodyAssignViewTemplate(ctx)
	q, err := a.DB.AssignViewTemplate(body.ViewUUID, body.ViewTemplateUUID, body.HostUUID)
	ResponseHandler(q, err, ctx)
}
