package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/NubeIO/rubix-os/schema/networklinkerschema"
	"github.com/NubeIO/rubix-os/schema/schema"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/plugin"
	"github.com/NubeIO/rubix-os/plugin/defaults"
	"github.com/NubeIO/rubix-os/plugin/nube/protocals/networklinker/linkmodel"
	"github.com/NubeIO/rubix-os/utils/helpers"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"

	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		body, err := inst.createNetwork(body)
		api.ResponseHandler(body, err, ctx)
	})
	mux.POST(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		body, err := inst.createDevice(body, nil, nil, nil, nil)
		api.ResponseHandler(body, err, ctx)
	})

	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		ok, err := inst.db.DeleteNetwork(body.UUID)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		ok, err := inst.db.DeleteDevice(body.UUID)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		ok, err := inst.db.DeletePoint(body.UUID)
		api.ResponseHandler(ok, err, ctx)
	})

	mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.db.UpdatePoint(body.UUID, body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.PATCH(plugin.PointsWriteURL, func(ctx *gin.Context) {
		inst.handlePointWriteProxy(ctx)
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		fns, err := inst.db.GetFlowNetworks(api.Args{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		fnsNames := make([]string, 0)
		for _, fn := range fns {
			fnsNames = append(fnsNames, fn.Name)
		}
		networkSchema := inst.GetNetworkSchemaOLD()
		networkSchema.AutoMappingFlowNetworkName.Options = fnsNames
		ctx.JSON(http.StatusOK, networkSchema)
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, inst.GetDeviceSchemaOLD())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, inst.GetPointSchemaOLD())
	})

	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, networklinkerschema.GetNetworkSchema())
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		deviceSchema := networklinkerschema.GetDeviceSchema()
		deviceSchema.AddressUUID.Options = inst.GetDeviceAddressUuidOptions()
		ctx.JSON(http.StatusOK, deviceSchema)
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, networklinkerschema.GetPointSchema())
	})
}

func (inst *Instance) GetNetworkSchemaOLD() *linkmodel.SchemaNetwork {
	netSchema := &linkmodel.SchemaNetwork{}
	defaults.Set(netSchema)
	options := inst.GetNetworkAddressUuidOption()
	for i := range options {
		netSchema.AddressUUID.Options = append(netSchema.AddressUUID.Options, options[i].Const)
	}
	return netSchema
}

func (inst *Instance) GetDeviceSchemaOLD() *linkmodel.SchemaDevice {
	devSchema := &linkmodel.SchemaDevice{}
	defaults.Set(devSchema)
	options := inst.GetDeviceAddressUuidOptions()
	for i := range options {
		devSchema.AddressUUID.Options = append(devSchema.AddressUUID.Options, options[i].Const)
	}
	return devSchema
}

func (inst *Instance) GetPointSchemaOLD() *linkmodel.SchemaPoint {
	point := &linkmodel.SchemaPoint{}
	defaults.Set(point)
	return point
}

func (inst *Instance) GetNetworkAddressUuidOption() []schema.OptionOneOf {
	var options []schema.OptionOneOf
	currNets, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	currNetIDs := make([][]string, len(currNets))
	for i, currNw := range currNets {
		currNetIDs[i] = strings.Split(currNw.AddressUUID, INTERNAL_SEPARATOR)
	}
	nets, _ := inst.db.GetNetworks(api.Args{})
	inner := 1
	for _, n := range nets {
		isWriter := inst.networkIsWriter(n)
		if strings.Contains(n.Name, UI_SEPARATOR) && strings.Contains(n.AddressUUID, INTERNAL_SEPARATOR) {
			continue
		}
		for j := inner; j < len(nets); j++ {
			if strings.Contains(nets[j].Name, UI_SEPARATOR) && strings.Contains(nets[j].AddressUUID, INTERNAL_SEPARATOR) {
				continue
			}
			if isWriter && inst.networkIsWriter(nets[j]) {
				continue
			}
			exists := false
			for _, currNw := range currNetIDs {
				if (n.UUID == currNw[0] && nets[j].UUID == currNw[1]) || (n.UUID == currNw[1] && nets[j].UUID == currNw[0]) {
					exists = true
					break
				}
			}
			if exists {
				continue
			}
			options = append(options, schema.OptionOneOf{
				Title: fmt.Sprintf("%s%s%s", n.Name, UI_SEPARATOR, nets[j].Name),
				Const: fmt.Sprintf("%s%s%s", n.UUID, INTERNAL_SEPARATOR, nets[j].UUID),
			})
		}
		inner++
	}
	return options
}

func (inst *Instance) GetDeviceAddressUuidOptions() []schema.OptionOneOf {
	var options []schema.OptionOneOf
	nets, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	for i := range nets {
		netSplit := strings.Split(nets[i].AddressUUID, INTERNAL_SEPARATOR)
		net1, _ := inst.db.GetNetwork(netSplit[0], api.Args{WithDevices: true})
		net2, _ := inst.db.GetNetwork(netSplit[1], api.Args{WithDevices: true})
		for i := range net1.Devices {
			for j := range net2.Devices {
				options = append(options, schema.OptionOneOf{
					Title: fmt.Sprintf("%s%s%s", net1.Devices[i].Name, UI_SEPARATOR, net2.Devices[j].Name),
					Const: fmt.Sprintf("%s%s%s", net1.Devices[i].UUID, INTERNAL_SEPARATOR, net2.Devices[j].UUID),
				})
			}
		}
	}
	return options
}

func (inst *Instance) handlePointWriteProxy(ctx *gin.Context) {
	uuid := plugin.ResolveID(ctx)

	reqCopy := helpers.CloneRequest(ctx)
	pointWriter := model.PointWriter{}
	binding.JSON.Bind(reqCopy, &pointWriter)

	newNet, newPointUUID := inst.getWriterNetworkAndPoint(uuid)
	newPath := strings.Replace(ctx.Request.URL.Path, pluginPath, newNet.PluginPath, 1)
	newPath = strings.Replace(newPath, uuid, *newPointUUID, 1)
	proxy := httputil.NewSingleHostReverseProxy(ctx.Request.URL)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.URL.Scheme = "http"
		req.Host = ctx.Request.Host
		req.URL.Host = ctx.Request.Host
		req.URL.Path = newPath
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			p, _ := inst.db.GetPoint(uuid, api.Args{})
			inst.syncPointSelected(p, *newPointUUID)
		}
		return nil
	}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
