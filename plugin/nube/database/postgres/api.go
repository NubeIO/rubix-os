package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/gin-gonic/gin"
)

const (
	relativePathHistories = "/histories"
	filterRegex           = "[^()!|<=>&']+|([&]+[&])+|([|]+[|])+|([!]+[=])+|([>]+[=])+|([<]+[=])+|[()<>=]"
	operatorFormat        = "{opt}"
	valueFormat           = "{val}"
)

var (
	logicalOperators        = []string{"&&", "||"}
	comparisonOperators     = []string{"=", ">", "<", "<=", ">=", "!="}
	orderOperators          = []string{"(", ")"}
	flowNetworkCloneColumns = []string{"global_uuid", "client_id", "client_name", "site_id", "site_name", "device_id",
		"device_name"}
	filterQueryMap = map[string]string{
		"timestamp":          "histories.timestamp" + operatorFormat + valueFormat,
		"value":              "histories.value" + operatorFormat + valueFormat,
		"rubix_network_uuid": "networks.uuid" + operatorFormat + valueFormat,
		"rubix_network_name": "networks.name" + operatorFormat + valueFormat,
		"rubix_device_uuid":  "devices.uuid" + operatorFormat + valueFormat,
		"rubix_device_name":  "devices.name" + operatorFormat + valueFormat,
		"rubix_point_uuid":   "points.uuid" + operatorFormat + valueFormat,
		"rubix_point_name":   "points.name" + operatorFormat + valueFormat,
		"global_uuid":        "flow_network_clones.global_uuid" + operatorFormat + valueFormat,
		"client_id":          "flow_network_clones.client_id" + operatorFormat + valueFormat,
		"client_name":        "flow_network_clones.client_name" + operatorFormat + valueFormat,
		"site_id":            "flow_network_clones.site_id" + operatorFormat + valueFormat,
		"site_name":          "flow_network_clones.site_name" + operatorFormat + valueFormat,
		"device_id":          "flow_network_clones.device_id" + operatorFormat + valueFormat,
		"tag": "(networks.uuid in (SELECT network_uuid FROM networks_tags WHERE tag_tag" +
			operatorFormat + valueFormat + ") OR devices.uuid in (SELECT device_uuid FROM devices_tags WHERE tag_tag" +
			operatorFormat + valueFormat + ") OR points.uuid in (SELECT point_uuid FROM points_tags WHERE tag_tag" +
			operatorFormat + valueFormat + "))",
		"meta_tag_key": "(networks.uuid in (SELECT network_uuid FROM network_meta_tags WHERE key" +
			operatorFormat + valueFormat + ") OR devices.uuid in (SELECT device_uuid FROM device_meta_tags WHERE key" +
			operatorFormat + valueFormat + ") OR points.uuid in (SELECT point_uuid FROM point_meta_tags WHERE key" +
			operatorFormat + valueFormat + "))",
		"meta_tag_value": "(networks.uuid in (SELECT network_uuid FROM network_meta_tags WHERE value" +
			operatorFormat + valueFormat + ") OR devices.uuid in (SELECT device_uuid FROM device_meta_tags WHERE value" +
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
