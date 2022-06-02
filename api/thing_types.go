package api

import (
	"github.com/NubeIO/flow-framework/src/units"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

// The ThingAPI
type ThingAPI struct {
}

func (t *ThingAPI) ThingClass(ctx *gin.Context) {
	ResponseHandler(model.ThingClass, nil, ctx)
}

func (t *ThingAPI) WriterActions(ctx *gin.Context) {
	ResponseHandler(model.WriterActions, nil, ctx)
}

func (t *ThingAPI) ThingUnits(ctx *gin.Context) {
	ResponseHandler(units.UnitsMap, nil, ctx)
}
