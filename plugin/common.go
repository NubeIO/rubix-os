package plugin

import (
	"fmt"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

const (
	NetworksURL    = "/networks"
	DevicesURL     = "/devices"
	PointsURL      = "/points"
	PointsWriteURL = "/points/write/:uuid"
)

func GetBODYNetwork(ctx *gin.Context) (dto *model.Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func GetBODYDevice(ctx *gin.Context) (dto *model.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func GetBODYPoint(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func GetBodyPointWriter(ctx *gin.Context) (dto *model.PointWriter, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func ResolveID(ctx *gin.Context) string {
	return ctx.Param("uuid")
}

func PointWrite(pnt *model.Point) (out float64) {
	out = float.NonNil(pnt.WriteValue)
	return
}

func SetStatusCode(code, defaultCode int) int {
	if code == 0 {
		return defaultCode
	} else {
		return code
	}
}

func ResponseHandler(body interface{}, err error, statusCode int, ctx *gin.Context) {
	if err == nil {
		ctx.JSON(SetStatusCode(statusCode, http.StatusOK), body)
	} else {
		switch err {
		case gorm.ErrRecordNotFound:
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 404, err.Error())
			ctx.JSON(http.StatusNotFound, interfaces.Message{Message: message})
		case gorm.ErrInvalidTransaction,
			gorm.ErrNotImplemented,
			gorm.ErrMissingWhereClause,
			gorm.ErrUnsupportedRelation,
			gorm.ErrPrimaryKeyRequired,
			gorm.ErrModelValueRequired,
			gorm.ErrInvalidData,
			gorm.ErrUnsupportedDriver,
			gorm.ErrRegistered,
			gorm.ErrInvalidField,
			gorm.ErrEmptySlice,
			gorm.ErrDryRunModeUnsupported,
			gorm.ErrInvalidDB,
			gorm.ErrInvalidValue,
			gorm.ErrInvalidValueOfLength:
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 500, err.Error())
			ctx.JSON(http.StatusInternalServerError, interfaces.Message{Message: message})
		default:
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 400, err.Error())
			ctx.JSON(http.StatusBadRequest, interfaces.Message{Message: message})
		}
	}
}
