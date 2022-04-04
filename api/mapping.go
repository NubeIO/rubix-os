package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type MappingDatabase interface {
	CreatePointMapping(body *model.PointMapping) (*model.PointMapping, error)
}
type MappingAPI struct {
	DB MappingDatabase
}

func (m *MappingAPI) CreatePointMapping(ctx *gin.Context) {
	body, _ := getBODYPointMapping(ctx)
	q, err := m.DB.CreatePointMapping(body)
	responseHandler(q, err, ctx)
}
