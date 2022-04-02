package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

// The HealthDatabase interface for encapsulating database access.
type HealthDatabase interface {
	Ping() error
}

// The HealthAPI provides handlers for the health information.
type HealthAPI struct {
	DB HealthDatabase
}

// Health returns health information.
func (a *HealthAPI) Health(ctx *gin.Context) {
	if err := a.DB.Ping(); err != nil {
		ctx.JSON(500, model.Health{
			Health:   model.StatusOrange,
			Database: model.StatusRed,
		})
		return
	}
	ctx.JSON(200, model.Health{
		Health:   model.StatusGreen,
		Database: model.StatusGreen,
	})
}
