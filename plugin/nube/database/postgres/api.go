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
	tagColumns              = []string{"tag"}
	flowNetworkCloneColumns = []string{"global_uuid", "client_id", "client_name", "site_id", "site_name", "device_id",
		"device_name"}
	filterMap = map[string]string{
		"timestamp":          "histories.timestamp",
		"value":              "histories.value",
		"rubix_network_uuid": "networks.uuid",
		"rubix_network_name": "networks.name",
		"rubix_device_uuid":  "devices.uuid",
		"rubix_device_name":  "devices.name",
		"rubix_point_uuid":   "points.uuid",
		"rubix_point_name":   "points.name",
		"tag":                "networks_tags.tag_tag,devices_tags.tag_tag,points_tags.tag_tag",
		"global_uuid":        "flow_network_clones.global_uuid",
		"client_id":          "flow_network_clones.client_id",
		"client_name":        "flow_network_clones.client_name",
		"site_id":            "flow_network_clones.site_id",
		"site_name":          "flow_network_clones.site_name",
		"device_id":          "flow_network_clones.device_id",
		"device_name":        "flow_network_clones.device_name",
	}
)

type Args struct {
	Filter *string
	Limit  *string
	Offset *string
}

var ArgsType = struct {
	Filter string
	Limit  string
	Offset string
}{
	Filter: "filter",
	Limit:  "limit",
	Offset: "offset",
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
