package api

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/gin-gonic/gin"
)

type AutoMappingDatabase interface {
	CreateAutoMapping(body *interfaces.AutoMapping) *interfaces.AutoMappingResponse
}

type AutoMappingAPI struct {
	DB AutoMappingDatabase
}

func (a *AutoMappingAPI) CreateAutoMapping(ctx *gin.Context) {
	body, _ := getBodyAutoMapping(ctx)
	resp := a.DB.CreateAutoMapping(body)
	ResponseHandler(resp, nil, ctx)
}
