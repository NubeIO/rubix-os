package api

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/gin-gonic/gin"
)

type AutoMappingDatabase interface {
	CreateAutoMapping(body *interfaces.AutoMapping) error
}

type AutoMappingAPI struct {
	DB AutoMappingDatabase
}

func (a *AutoMappingAPI) CreateAutoMapping(ctx *gin.Context) {
	body, _ := getBodyAutoMapping(ctx)
	err := a.DB.CreateAutoMapping(body)
	if err != nil {
		ResponseHandler(false, err, ctx)
		return
	}
	ResponseHandler(true, err, ctx)
}
