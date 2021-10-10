package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The ScheduleDatabase interface for encapsulating database access.
type ScheduleDatabase interface {
	GetSchedule(uuid string) (*model.Schedule, error)
	GetSchedules() ([]*model.Schedule, error)
	CreateSchedule(body *model.Schedule) (*model.Schedule, error)
	UpdateSchedule(uuid string, body *model.Schedule) (*model.Schedule, error)
	DeleteSchedule(uuid string) (bool, error)
	DropSchedules() (bool, error)
}
type ScheduleAPI struct {
	DB ScheduleDatabase
}

func (a *ScheduleAPI) GetSchedules(ctx *gin.Context) {
	q, err := a.DB.GetSchedules()
	reposeHandler(q, err, ctx)
}

func (a *ScheduleAPI) GetSchedule(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetSchedule(uuid)
	reposeHandler(q, err, ctx)
}

func (a *ScheduleAPI) UpdateSchedule(ctx *gin.Context) {
	body, _ := getBODYSchedule(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateSchedule(uuid, body)
	reposeHandler(q, err, ctx)
}

func (a *ScheduleAPI) CreateSchedule(ctx *gin.Context) {
	body, _ := getBODYSchedules(ctx)
	reposeHandler(body, nil, ctx) //TODO
}

func (a *ScheduleAPI) DeleteSchedule(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteSchedule(uuid)
	reposeHandler(q, err, ctx)
}

func (a *ScheduleAPI) DropSchedules(ctx *gin.Context) {
	q, err := a.DB.DropSchedules()
	reposeHandler(q, err, ctx)
}
