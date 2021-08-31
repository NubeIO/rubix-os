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
	Sort         string
	Order        string
	Offset       string
	Limit        string
	Search       string
	WithChildren string
	WithPoints   string
	AskRefresh   string
	AskResponse  string
	Write 		 string
	ThingType 	 string
	FlowNetworkUUID string
}

var ArgsType = struct {
	Sort         string
	Order        string
	Offset       string
	Limit        string
	Search       string
	WithChildren string
	WithPoints   string
	AskRefresh 	 string
	AskResponse  string
	Write 		 string
	ThingType 	 string
	FlowNetworkUUID string
}{
	Sort:         "sort",
	Order:        "order",
	Offset:       "offset",
	Limit:        "limit",
	Search:       "search",
	WithChildren: "with_children",
	WithPoints:   "with_points",
	AskRefresh:   "ask_refresh",
	AskResponse:  "ask_response",
	Write:  	  "write", //consumer to write a value
	ThingType:    	"thing_type", //the type of thing like a point
	FlowNetworkUUID:"flow_network_uuid", //the type of thing like a point


}

var ArgsDefault = struct {
	Sort         string
	Order        string
	Offset       string
	Limit        string
	Search       string
	WithChildren string
	WithPoints   string
	AskRefresh 	 string
	AskResponse  string
	Write        string
	ThingType 	 string
	FlowNetworkUUID string
}{
	Sort:         "ID",
	Order:        "DESC",
	Offset:       "0",
	Limit:        "25",
	Search:       "",
	WithChildren: "false",
	WithPoints:   "false",
	AskRefresh:   "false",
	AskResponse:  "false",
	Write:        "false",
	ThingType:    "point",
	FlowNetworkUUID:    "",
}

func withID(ctx *gin.Context, name string, f func(id uint)) {
	if id, err := strconv.ParseUint(ctx.Param(name), 10, bits.UintSize); err == nil {
		f(uint(id))
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}

func getBODYRubixPlat(ctx *gin.Context) (dto *model.RubixPlat, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYFlowNetwork(ctx *gin.Context) (dto *model.FlowNetwork, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYNetwork(ctx *gin.Context) (dto *model.Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYHistory(ctx *gin.Context) (dto *model.ProducerHistory, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYBulkHistory(ctx *gin.Context) (dto []*model.ProducerHistory, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYDevice(ctx *gin.Context) (dto *model.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYProducer(ctx *gin.Context) (dto *model.Producer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYConsumer(ctx *gin.Context) (dto *model.Consumer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriter(ctx *gin.Context) (dto *model.Writer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterClone(ctx *gin.Context) (dto *model.WriterClone, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYGateway(ctx *gin.Context) (dto *model.Stream, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYCommandGroup(ctx *gin.Context) (dto *model.CommandGroup, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYJobs(ctx *gin.Context) (dto *model.Job, err error) {
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
func resolveName(ctx *gin.Context) string {
	id := ctx.Param("name")
	return id
}

func resolvePath(ctx *gin.Context) string {
	id := ctx.Param("path")
	return id
}

func WithChildren(value string) (bool, error) {
	if value == "" {
		return false, nil
	} else {
		c, err := bools.Boolean(value)
		return c, err
	}
}

func toBool(value string) (bool, error) {
	if value == "" {
		return false, nil
	} else {
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
		b, _ := json.MarshalIndent(model, "", "  ")
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
	code     int
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
		"response": Response,
	}}
}

func Failed(code int, Response interface{}) Response {
	return &BaseResponse{code: code, Response: JSON{
		"status": "failed",
		"error":  Response,
	}}
}
