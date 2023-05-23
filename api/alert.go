package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type AlertStatus struct {
	Status string `json:"status"`
}

func getAlertStatus(ctx *gin.Context) (status string, err error) {
	statusStruct := AlertStatus{}
	err = ctx.ShouldBindJSON(&statusStruct)
	return statusStruct.Status, err
}

type AlertDatabase interface {
	GetAlert(uuid string) (*model.Alert, error)
	GetAlerts() ([]*model.Alert, error)
	GetAlertsByHost(hostUUID string) ([]*model.Alert, error)
	GetAlertByField(field string, value string) (*model.Alert, error)
	CreateAlert(body *model.Alert) (*model.Alert, error)
	UpdateAlertStatus(uuid string, status string) (alert *model.Alert, err error)
	DeleteAlert(uuid string) (bool, error)
	DropAlerts() (bool, error)
}
type AlertAPI struct {
	DB AlertDatabase
}

func (a *AlertAPI) AlertsSchema(ctx *gin.Context) {
}

func (a *AlertAPI) GetAlert(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetAlert(uuid)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *AlertAPI) GetAlertsByHost(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetAlertsByHost(uuid)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *AlertAPI) GetAlerts(ctx *gin.Context) {
	q, err := a.DB.GetAlerts()
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *AlertAPI) CreateAlert(ctx *gin.Context) {
	body, _ := getBodyAlert(ctx)
	q, err := a.DB.CreateAlert(body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *AlertAPI) UpdateAlertStatus(ctx *gin.Context) {
	status, _ := getAlertStatus(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateAlertStatus(uuid, status)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *AlertAPI) DeleteAlert(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteAlert(uuid)
	if err != nil {
		ResponseHandler(nil, err, ctx)
	} else {
		ResponseHandler(q, err, ctx)
	}
}

func (a *AlertAPI) DropAlerts(c *gin.Context) {
	q, err := a.DB.DropAlerts()
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(q, err, c)
}