package api

import (
	"errors"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/bools"
	"math/bits"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Args struct {
	Sort            string
	Order           string
	Offset          string
	Limit           string
	Search          string
	WithChildren    string
	WithPoints      string
	AskRefresh      string
	AskResponse     string
	Write           string
	ThingType       string
	FlowNetworkUUID string
	WriteHistory    string
	WriteConsumer   string
	Field           string
	Value           string
	UpdateProducer  string
	CompactPayload  string
	CompactWithName string
}

var ArgsType = struct {
	Sort            string
	Order           string
	Offset          string
	Limit           string
	Search          string
	WithChildren    string
	WithPoints      string
	AskRefresh      string
	AskResponse     string
	Write           string
	ThingType       string
	FlowNetworkUUID string
	WriteHistory    string
	WriteConsumer   string
	Field           string
	Value           string
	UpdateProducer  string
	CompactPayload  string //for a point would be presentValue
	CompactWithName string //for a point would be presentValue and pointName
}{
	Sort:            "sort",
	Order:           "order",
	Offset:          "offset",
	Limit:           "limit",
	Search:          "search",
	WithChildren:    "with_children",
	WithPoints:      "with_points",
	AskRefresh:      "ask_refresh",
	AskResponse:     "ask_response",
	Write:           "write",             //consumer to write a value
	ThingType:       "thing_type",        //the type of thing like a point
	FlowNetworkUUID: "flow_network_uuid", //the type of thing like a point
	WriteHistory:    "write_history",
	WriteConsumer:   "write_consumer",
	Field:           "field",
	Value:           "value",
	UpdateProducer:  "update_producer",
	CompactPayload:  "compact_payload",
	CompactWithName: "compact_with_name",
}

var ArgsDefault = struct {
	Sort            string
	Order           string
	Offset          string
	Limit           string
	Search          string
	WithChildren    string
	WithPoints      string
	AskRefresh      string
	AskResponse     string
	Write           string
	ThingType       string
	FlowNetworkUUID string
	Field           string
	Value           string
	UpdateProducer  string
	CompactPayload  string
	CompactWithName string
}{
	Sort:            "ID",
	Order:           "DESC", //ASC or DESC
	Offset:          "0",
	Limit:           "25",
	Search:          "",
	WithChildren:    "false",
	WithPoints:      "false",
	AskRefresh:      "false",
	AskResponse:     "false",
	Write:           "false",
	ThingType:       "point",
	FlowNetworkUUID: "",
	Field:           "name",
	Value:           "",
	UpdateProducer:  "false",
	CompactPayload:  "false",
	CompactWithName: "false",
}

func withID(ctx *gin.Context, name string, f func(id uint)) {
	if id, err := strconv.ParseUint(ctx.Param(name), 10, bits.UintSize); err == nil {
		f(uint(id))
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}

func getBODYSchedule(ctx *gin.Context) (dto *model.Schedule, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
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
func getBODYNode(ctx *gin.Context) (dto *model.Node, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPlugin(ctx *gin.Context) (dto *model.PluginConf, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYMqttConnection(ctx *gin.Context) (dto *model.MqttConnection, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYIntegration(ctx *gin.Context) (dto *model.Integration, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterBody(ctx *gin.Context) (dto *model.WriterBody, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterBulk(ctx *gin.Context) (dto []*model.WriterBulk, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterClone(ctx *gin.Context) (dto *model.WriterClone, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYStream(ctx *gin.Context) (dto *model.Stream, err error) {
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

func resolvePluginUUID(ctx *gin.Context) string {
	id := ctx.Param("plugin_uuid")
	return id
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

func toBool(value string) (bool, error) {
	if value == "" {
		return false, nil
	} else {
		c, err := bools.Boolean(value)
		return c, err
	}
}
