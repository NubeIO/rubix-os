package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type ScheduleDatabase interface {
	GetSchedules() ([]*model.Schedule, error)
	GetSchedule(uuid string) (*model.Schedule, error)
	GetOneScheduleByArgs(Args) (*model.Schedule, error)
	CreateSchedule(body *model.Schedule) (*model.Schedule, error)
	UpdateSchedule(uuid string, body *model.Schedule) (*model.Schedule, error)
	ScheduleWrite(uuid string, body *model.ScheduleData) error
	DeleteSchedule(uuid string) (bool, error)
}

type ScheduleAPI struct {
	DB ScheduleDatabase
}

func (a *ScheduleAPI) GetSchedule(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetSchedule(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *ScheduleAPI) GetOneScheduleByArgs(ctx *gin.Context) {
	args := buildScheduleArgs(ctx)
	q, err := a.DB.GetOneScheduleByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (a *ScheduleAPI) GetSchedules(ctx *gin.Context) {
	q, err := a.DB.GetSchedules()
	ResponseHandler(q, err, ctx)
}

func (a *ScheduleAPI) UpdateSchedule(ctx *gin.Context) {
	body, _ := getBODYSchedule(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateSchedule(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *ScheduleAPI) ScheduleWrite(ctx *gin.Context) {
	body, _ := getBODYScheduleData(ctx)
	uuid := resolveID(ctx)
	err := a.DB.ScheduleWrite(uuid, body)
	ResponseHandler(nil, err, ctx)
}

func (a *ScheduleAPI) CreateSchedule(ctx *gin.Context) {
	body, _ := getBODYSchedule(ctx)
	q, err := a.DB.CreateSchedule(body)
	ResponseHandler(q, err, ctx)
}

func (a *ScheduleAPI) DeleteSchedule(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteSchedule(uuid)
	ResponseHandler(q, err, ctx)
}
