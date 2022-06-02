package nerrors

import (
	"fmt"
	"github.com/NubeIO/flow-framework/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFound creates a gin middleware for handling page not found.
func NotFound() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 404, "api not found")
		ctx.JSON(http.StatusNotFound, interfaces.Message{Message: message})
	}
}
