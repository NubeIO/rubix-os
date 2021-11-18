package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The ScheduleDatabase interface for encapsulating database access.
type ScheduleDatabase interface {
	GetSchedule(uuid string) (*model.Schedule, error)
	GetScheduleByField(field, value string) (*model.Schedule, error)
	GetSchedules() ([]*model.Schedule, error)
	CreateSchedule(body *model.Schedule) (*model.Schedule, error)
	UpdateSchedule(uuid string, body *model.Schedule) (*model.Schedule, error)
	DeleteSchedule(uuid string) (bool, error)
	DropSchedules() (bool, error)
}
type ScheduleAPI struct {
	DB ScheduleDatabase
}

func (a *ScheduleAPI) GetSchedule(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetSchedule(uuid)
	responseHandler(q, err, ctx)
}

func (a *ScheduleAPI) GetScheduleByField(ctx *gin.Context) {
	field, value := withFieldsArgs(ctx)
	q, err := a.DB.GetScheduleByField(field, value)
	responseHandler(q, err, ctx)
}

func (a *ScheduleAPI) GetSchedules(ctx *gin.Context) {
	q, err := a.DB.GetSchedules()
	responseHandler(q, err, ctx)
}

func (a *ScheduleAPI) UpdateSchedule(ctx *gin.Context) {
	body, _ := getBODYSchedule(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateSchedule(uuid, body)
	responseHandler(q, err, ctx)
}

func (a *ScheduleAPI) CreateSchedule(ctx *gin.Context) {
	body, _ := getBODYSchedule(ctx)
	q, err := a.DB.CreateSchedule(body)
	responseHandler(q, err, ctx) //TODO
}

func (a *ScheduleAPI) DeleteSchedule(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteSchedule(uuid)
	responseHandler(q, err, ctx)
}

func (a *ScheduleAPI) DropSchedules(ctx *gin.Context) {
	q, err := a.DB.DropSchedules()
	responseHandler(q, err, ctx)
}
