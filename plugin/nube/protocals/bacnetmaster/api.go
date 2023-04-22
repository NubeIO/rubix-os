package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetmaster/master"
	"github.com/NubeIO/lib-schema/masterschema"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	schemaNetwork     = "/schema/network"
	schemaDevice      = "/schema/device"
	schemaPoint       = "/schema/point"
	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
	whois             = "/whois"
	discoverPoints    = "/device/points"
)

func resolveID(ctx *gin.Context) string {
	return ctx.Param("uuid")
}

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.addNetwork(body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.GET(plugin.NetworksURL, func(ctx *gin.Context) {
		network, err := inst.getNetworks()
		api.ResponseHandler(network, err, ctx)
	})
	mux.POST(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.addDevice(body)
		api.ResponseHandler(device, err, ctx)
	})
	mux.POST(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.addPoint(body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.PATCH(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.updateNetwork(body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.PATCH(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.updateDevice(body)
		api.ResponseHandler(device, err, ctx)
	})
	mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.updatePoint(body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.PATCH(plugin.PointsWriteURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBodyPointWriter(ctx)
		uuid := plugin.ResolveID(ctx)
		point, err := inst.writePoint(uuid, body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		ok, err := inst.deleteNetwork(body)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		ok, err := inst.deleteDevice(body)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		ok, err := inst.deletePoint(body)
		api.ResponseHandler(ok, err, ctx)
	})

	mux.POST(whois+"/:uuid", func(ctx *gin.Context) {
		// body, _ := master.BodyWhoIs(ctx)
		// uuid := resolveID(ctx)
		// addDevices := ctx.Query("add_devices")
		// add, _ := strconv.ParseBool(addDevices)
		// resp, err := inst.whoIs(uuid, body, add)
		// api.ResponseHandler(resp, err, ctx)
	})
	mux.POST("/master/whois", func(ctx *gin.Context) {
		body, _ := bodyMasterWhoIs(ctx)
		resp, err := inst.masterWhoIs(body)
		api.ResponseHandler(resp, err, ctx)
	})
	mux.POST(discoverPoints+"/:uuid", func(ctx *gin.Context) {
		// uuid := resolveID(ctx)
		// addPoints := ctx.Query("add_points")
		// add, _ := strconv.ParseBool(addPoints)
		// makeWriteablePoints := ctx.Query("writeable_points")
		// writeable, _ := strconv.ParseBool(makeWriteablePoints)
		// resp, err := inst.devicePoints(uuid, add, writeable)
		// api.ResponseHandler(resp, err, ctx)
	})

	mux.GET("/polling/stats/network/:name", func(ctx *gin.Context) {
		networkName := ctx.Param("name")
		stats, err := inst.getPollingStats(networkName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, stats)
			return
		}
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, master.GetNetworkSchema())
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, master.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, master.GetPointSchema())
	})
	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		fns, err := inst.db.GetFlowNetworks(api.Args{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		networkSchema := masterschema.GetNetworkSchema()
		networkSchema.AutoMappingFlowNetworkName.Options = plugin.GetFlowNetworkNames(fns)
		ctx.JSON(http.StatusOK, networkSchema)
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, masterschema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, masterschema.GetPointSchema())
	})
}

func (inst *Instance) getPollingStats(networkName string) (result *model.PollQueueStatistics, error error) {
	if len(inst.NetworkPollManagers) == 0 {
		return nil, errors.New("couldn't find any plugin network poll managers")
	}
	for _, netPollMan := range inst.NetworkPollManagers {
		if netPollMan == nil || netPollMan.NetworkName != networkName {
			continue
		}
		result = netPollMan.GetPollingQueueStatistics()
		return result, nil
	}
	return nil, errors.New(fmt.Sprintf("couldn't find network %s for polling statistics", networkName))
}
