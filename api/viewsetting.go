package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type ViewSettingDatabase interface {
	GetViewSetting() (*model.ViewSetting, error)
	CreateViewSetting(body *model.ViewSetting) (*model.ViewSetting, error)
	DeleteViewSetting() (bool, error)
}

type ViewSettingAPI struct {
	DB ViewSettingDatabase
}

func (a *ViewSettingAPI) GetViewSetting(ctx *gin.Context) {
	q, err := a.DB.GetViewSetting()
	ResponseHandler(q, err, ctx)
}

func (a *ViewSettingAPI) CreateViewSetting(ctx *gin.Context) {
	body, _ := getBodyViewSetting(ctx)
	q, err := a.DB.CreateViewSetting(body)
	ResponseHandler(q, err, ctx)
}

func (a *ViewSettingAPI) DeleteViewSetting(ctx *gin.Context) {
	q, err := a.DB.DeleteViewSetting()
	ResponseHandler(q, err, ctx)
}
