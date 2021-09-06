package api

import (
	"github.com/NubeDev/flow-framework/handler"
	"github.com/NubeDev/flow-framework/system"
	"github.com/gin-gonic/gin"
	"net/http"
)

//api/system/ping
//api/system/time
//api/system/memory
//api/system/disc
//POST
//api/system/host/restart
// body {action:"restart"}

//3rd party apps like lorawan

//SYSTEM SERVICE
// GET
///api/system/service
/// return all the services
//POST
///api/system/service
// body {service:"lorawan" , action:"restart"}

const (
	HealhStatusUp   = "UP"
	HealhStatusDown = "DOWN"
)

type HealthsAPI struct {
	Handler *handler.Handler
}

type Healths struct {
	Status string `json:"status"`
}

func (healths *HealthsAPI) Hostname(ctx *gin.Context) {
	healths.Handler.Get()
	h := Healths{}
	h.Status, _ = system.Hostname()
	ctx.JSON(http.StatusOK, h)
}
