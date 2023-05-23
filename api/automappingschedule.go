package api

import (
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type AutoMappingScheduleDatabase interface {
	CreateAutoMappingSchedule(body *interfaces.AutoMapping) *interfaces.AutoMappingScheduleResponse
}

type AutoMappingScheduleAPI struct {
	DB AutoMappingScheduleDatabase
}

func (a *AutoMappingScheduleAPI) CreateAutoMappingSchedule(ctx *gin.Context) {
	body, _ := getBodyAutoMapping(ctx)
	resp := a.DB.CreateAutoMappingSchedule(body)
	ResponseHandler(resp, nil, ctx)
}
