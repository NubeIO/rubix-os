package nerrors

import (
	"fmt"
	"github.com/NubeIO/rubix-os/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFoundHandler creates a gin middleware for handling page not found.
func NotFoundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 404, "api not found")
		ctx.JSON(http.StatusNotFound, interfaces.Message{Message: message})
	}
}
