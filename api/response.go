package api

import (
	"fmt"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/nerrors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ResponseHandler(body interface{}, err error, ctx *gin.Context) {
	if err == nil {
		ctx.JSON(http.StatusOK, body)
	} else {
		switch err.(type) {
		case *nerrors.ErrConflict:
			statusCode := http.StatusConflict
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, statusCode, err.Error())
			ctx.JSON(statusCode, interfaces.Message{Message: message})
		case *nerrors.NotFound:
			statusCode := http.StatusNotFound
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, statusCode, err.Error())
			ctx.JSON(statusCode, interfaces.Message{Message: message})
		case *nerrors.ErrUnauthorized:
			statusCode := http.StatusUnauthorized
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, statusCode, err.Error())
			ctx.JSON(statusCode, interfaces.Message{Message: message})
		default:
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
				statusCode := http.StatusInternalServerError
				message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, statusCode, err.Error())
				ctx.JSON(statusCode, interfaces.Message{Message: message})
			default:
				statusCode := http.StatusBadRequest
				message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, statusCode, err.Error())
				ctx.JSON(statusCode, interfaces.Message{Message: message})
			}
		}
	}
}
