package api

import (
	"github.com/NubeDev/plug-framework/system"
	"github.com/gin-gonic/gin"
	"net/http"
)

//const (
//	HealhStatusUp   = "UP"
//	HealhStatusDown = "DOWN"
//)
//
//type Health struct {
//	Status string `json:"status"`
//}

func Control(ctx *gin.Context) {
	h := Health{}
	h.Status, _ = system.Hostname()
	ctx.JSON(http.StatusOK, h)
}
