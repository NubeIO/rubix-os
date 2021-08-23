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
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())
}


func (a *RubixPlatAPI) UpdateRubixPlat(ctx *gin.Context) {
	body, _ := getBODYRubixPlat(ctx)
	q, err := a.DB.UpdateRubixPlat(body)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}
