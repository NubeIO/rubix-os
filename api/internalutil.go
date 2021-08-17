package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/bools"
	"math/bits"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)


type Args struct {
	Sort   			string
	Order  			string
	Offset 			string
	Limit  			string
	Search 			string
	WithChildren 	string
	WithPoints          string
}


var ArgsType = struct {
	Sort   				string
	Order   			string
	Offset   			string
	Limit   			string
	Search   			string
	WithChildren 		string
	WithPoints          string
}{
	Sort:   			"Sort",
	Order:   			"Order",
	Offset:   			"Offset",
	Limit:   			"Limit",
	Search:   			"Search",
	WithChildren:   	"with_children",
	WithPoints:   		"with_points",
}

var ArgsDefault = struct {
	Sort   			string
	Order   		string
	Offset   		string
	Limit   		string
	Search   		string
	WithChildren   	string
	WithPoints      string
}{
	Sort:   			"ID",
	Order:   			"DESC",
	Offset:   			"0",
	Limit:   			"25",
	Search:   			"",
	WithChildren:   	"false",
	WithPoints:   		"false",
}




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

func getBODYDevice(ctx *gin.Context) (dto *model.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPoint(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	id := ctx.Param("uuid")
	return id
}


func WithChildren(value string) (bool, error)  {
	if value == "" {
		return false, nil
	} else  {
		c, err := bools.Boolean(value)
		return c, err
	}
}



func OK(resp interface{}) Response {
	return Success(http.StatusOK, resp)
}

func OKWithMessage(resp string) Response {
	return Success(http.StatusOK, resp)
}


func BadEntity(excepted string) Response {
	return Failed(http.StatusUnprocessableEntity, excepted)
}


func NotFound(err string) Response {
	return Failed(http.StatusNotFound, err)
}


func Created(id string) Response {
	return Success(http.StatusCreated, JSON{"id": id})
}

func Data(model interface{}) Response {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Slice {
		b,_:=json.MarshalIndent(model, "", "  ")
		fmt.Print(string(b))
		return Success(http.StatusOK, JSON{"count": v.Len(), "items": model})
	}
	return Success(http.StatusOK, model)
}

type JSON map[string]interface{}

type Response interface {
	GetResponse() map[string]interface{}
	GetStatusCode() int
}

type BaseResponse struct {
	Response JSON
	code int
}

func (r *BaseResponse) GetResponse() map[string]interface{} {
	return r.Response
}

func (r *BaseResponse) GetStatusCode() int {
	return r.code
}

func Success(code int, Response interface{}) Response {
	return &BaseResponse{code: code, Response: JSON{
		"status":   "success",
		"Response": Response,
	}}
}

func Failed(code int, Response interface{}) Response {
	return &BaseResponse{code: code, Response: JSON{
		"status": "failed",
		"error":  Response,
	}}
}

