package api

import (
	"github.com/NubeDev/flow-framework/model"
	unit "github.com/NubeDev/flow-framework/src/units"
	"github.com/gin-gonic/gin"
)

// The ThingAPI
type ThingAPI struct {
}

func (t *ThingAPI) ThingClass(ctx *gin.Context) {
	reposeHandler(model.ThingClass, nil, ctx)
}

func (t *ThingAPI) WriterActions(ctx *gin.Context) {
	reposeHandler(model.WriterActions, nil, ctx)
}

func (t *ThingAPI) ThingUnits(ctx *gin.Context) {
	reposeHandler(unit.UnitsMap, nil, ctx)
}
