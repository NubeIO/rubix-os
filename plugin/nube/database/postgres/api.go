package main

import (
	"github.com/NubeIO/rubix-os/api"
	"github.com/gin-gonic/gin"
)

const (
	relativePathHistories = "/histories"
	filterRegex           = "[^()!|<=>&']+|([&]+[&])+|([|]+[|])+|([!]+[=])+|([>]+[=])+|([<]+[=])+|[()<>=]"
	operatorFormat        = "{opt}"
	valueFormat           = "{val}"
)

var (
	logicalOperators    = []string{"&&", "||"}
	comparisonOperators = []string{"=", ">", "<", "<=", ">=", "!="}
	orderOperators      = []string{"(", ")"}
	filterQueryMap      = map[string]string{
		"location_uuid": "points.location_uuid" + operatorFormat + valueFormat,
		"location_name": "points.location_name" + operatorFormat + valueFormat,
		"group_uuid":    "points.group_id" + operatorFormat + valueFormat,
		"group_name":    "points.group_name" + operatorFormat + valueFormat,
		"host_uuid":     "points.host_uuid" + operatorFormat + valueFormat,
		"host_name":     "points.host_name" + operatorFormat + valueFormat,
		"global_uuid":   "points.global_uuid" + operatorFormat + valueFormat,
		"network_uuid":  "points.network_uuid" + operatorFormat + valueFormat,
		"network_name":  "points.network_name" + operatorFormat + valueFormat,
		"device_uuid":   "points.device_uuid" + operatorFormat + valueFormat,
		"device_name":   "points.device_name" + operatorFormat + valueFormat,
		"point_uuid":    "points.uuid" + operatorFormat + valueFormat,
		"point_name":    "points.name" + operatorFormat + valueFormat,
		"timestamp":     "histories.timestamp" + operatorFormat + valueFormat,
		"value":         "histories.value" + operatorFormat + valueFormat,
		"tag": "(points.network_uuid in (SELECT network_uuid FROM network_tags WHERE tag" +
			operatorFormat + valueFormat + ") OR points.device_uuid in (SELECT device_uuid FROM device_tags WHERE tag" +
			operatorFormat + valueFormat + ") OR points.uuid in (SELECT point_uuid FROM point_tags WHERE tag" +
			operatorFormat + valueFormat + "))",
		"meta_tag_key": "(points.network_uuid in (SELECT network_uuid FROM network_meta_tags WHERE key" +
			operatorFormat + valueFormat + ") OR points.device_uuid in (SELECT device_uuid FROM device_meta_tags WHERE key" +
			operatorFormat + valueFormat + ") OR points.uuid in (SELECT point_uuid FROM point_meta_tags WHERE key" +
			operatorFormat + valueFormat + "))",
		"meta_tag_value": "(points.network_uuid in (SELECT network_uuid FROM network_meta_tags WHERE value" +
			operatorFormat + valueFormat + ") OR points.device_uuid in (SELECT device_uuid FROM device_meta_tags WHERE value" +
			operatorFormat + valueFormat + ") OR points.uuid in (SELECT point_uuid FROM point_meta_tags WHERE value" +
			operatorFormat + valueFormat + "))",
	}
)

type Args struct {
	Filter     *string
	Limit      *string
	Offset     *string
	OrderBy    *string
	Order      *string
	GroupLimit *string
}

var ArgsType = struct {
	Filter     string
	Limit      string
	Offset     string
	OrderBy    string
	Order      string
	GroupLimit string
}{
	Filter:     "filter",
	Limit:      "limit",
	Offset:     "offset",
	OrderBy:    "order_by",
	Order:      "order",
	GroupLimit: "group_limit",
}

func buildHistoryArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	if value, ok := ctx.GetQuery(aType.Filter); ok {
		args.Filter = &value
	}
	if value, ok := ctx.GetQuery(aType.Limit); ok {
		args.Limit = &value
	}
	if value, ok := ctx.GetQuery(aType.Offset); ok {
		args.Offset = &value
	}
	if value, ok := ctx.GetQuery(aType.OrderBy); ok {
		args.OrderBy = &value
	}
	if value, ok := ctx.GetQuery(aType.Order); ok {
		args.Order = &value
	}
	if value, ok := ctx.GetQuery(aType.GroupLimit); ok {
		args.GroupLimit = &value
	}
	return args
}

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.GET(relativePathHistories, func(ctx *gin.Context) {
		args := buildHistoryArgs(ctx)
		q, err := inst.getHistories(args)
		api.ResponseHandler(q, err, ctx)
	})
}
