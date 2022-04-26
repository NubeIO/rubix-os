package plugin

import (
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

const (
	NetworksURL    = "/networks"
	DevicesURL     = "/devices"
	PointsURL      = "/points"
	PointsWriteURL = "/write"
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

func PointWrite(pnt *model.Point) (out float64) {
	out = utils.Float64IsNil(pnt.WriteValue)
	return
}

type Message struct {
	Message string `json:"message"`
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
			ctx.JSON(SetStatusCode(statusCode, http.StatusNotFound), Message{Message: err.Error()})
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
			ctx.JSON(SetStatusCode(statusCode, http.StatusInternalServerError), Message{Message: err.Error()})
		default:
			ctx.JSON(SetStatusCode(statusCode, http.StatusBadRequest), Message{Message: err.Error()})
		}
	}
	//return nil
}
