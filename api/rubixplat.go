package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The RubixPlatDatabase interface for encapsulating database access.
type RubixPlatDatabase interface {
	GetRubixPlat() (*model.RubixPlat, error)
	UpdateRubixPlat(body *model.RubixPlat) (*model.RubixPlat, error)

}
type RubixPlatAPI struct {
	DB RubixPlatDatabase
}


func (a *RubixPlatAPI) GetRubixPlat(ctx *gin.Context) {
	q, err := a.DB.GetRubixPlat()
	reposeHandler(q, err, ctx)
}


func (a *RubixPlatAPI) UpdateRubixPlat(ctx *gin.Context) {
	body, _ := getBODYRubixPlat(ctx)
	q, err := a.DB.UpdateRubixPlat(body)
	reposeHandler(q, err, ctx)
}
