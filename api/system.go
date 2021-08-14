package api

import (
	"github.com/NubeDev/plug-framework/system"
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

type Health struct {
	Status string `json:"status"`
}

func Hostname(ctx *gin.Context) {
	h := Health{}
	h.Status, _ = system.Hostname()
	ctx.JSON(http.StatusOK, h)
}
