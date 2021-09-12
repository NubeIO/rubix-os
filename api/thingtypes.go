package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)


// The ThingAPI
type ThingAPI struct {

}


func (t *ThingAPI) ThingClass(ctx *gin.Context) {
	reposeHandler(model.ThingClass, nil, ctx)
}

var ThingTypes interface{}

func (t *ThingAPI) ThingTypes(ctx *gin.Context) {
	reposeHandler(ThingTypes, nil, ctx)
}

