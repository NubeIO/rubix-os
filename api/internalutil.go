package api

import (
	"errors"
	"github.com/NubeDev/plug-framework/model"
	"math/bits"
	"strconv"

	"github.com/gin-gonic/gin"
)

func withID(ctx *gin.Context, name string, f func(id uint)) {
	if id, err := strconv.ParseUint(ctx.Param(name), 10, bits.UintSize); err == nil {
		f(uint(id))
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}



func getBODY(ctx *gin.Context) (dto *model.Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	id := ctx.Param("uuid")
	return id
}
