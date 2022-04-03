package error

import (
	"net/http"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

// NotFound creates a gin middleware for handling page not found.
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, &model.Error{
			Error:            http.StatusText(http.StatusNotFound),
			ErrorCode:        http.StatusNotFound,
			ErrorDescription: "page not found",
		})
	}
}
