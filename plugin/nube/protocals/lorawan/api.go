package main

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	"github.com/NubeIO/lib-schema/lorawanschema"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"

	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
)

func resolveID(ctx *gin.Context) string {
	return ctx.Param("eui")
}

func bodyDevice(ctx *gin.Context) (dto *csrest.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyDeviceAdd(ctx *gin.Context) (dto *csrest.DeviceAdd, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyActivateDevice(ctx *gin.Context) (dto *csrest.DeviceActivation, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyDeviceKey(ctx *gin.Context) (dto *csrest.DeviceKey, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.PATCH(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.db.UpdateNetwork(body.UUID, body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.PATCH(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.db.UpdateDevice(body.UUID, body)
		api.ResponseHandler(device, err, ctx)
	})
	mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.db.UpdatePoint(body.UUID, body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		if inst.enabled {
			err := errors.New("cannot delete lorawan network when plugin is enabled")
			api.ResponseHandler(false, err, ctx)
			return
		}
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

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetNetworkSchema())
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetPointSchema())
	})

	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		fns, err := inst.db.GetFlowNetworks(api.Args{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		networkSchema := lorawanschema.GetNetworkSchema()
		networkSchema.AutoMappingFlowNetworkName.Options = plugin.GetFlowNetworkNames(fns)
		ctx.JSON(http.StatusOK, networkSchema)
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lorawanschema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lorawanschema.GetPointSchema())
	})

	// ---------------------- CS APIs ------------------------------

	mux.GET("/cs/applications", func(ctx *gin.Context) {
		p, err := inst.REST.GetApplications()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/cs/gateways", func(ctx *gin.Context) {
		p, err := inst.REST.GetGateways()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/cs/device-profiles", func(ctx *gin.Context) {
		p, err := inst.REST.GetDeviceProfiles()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/cs/service-profiles", func(ctx *gin.Context) {
		p, err := inst.REST.GetServiceProfiles()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/cs/gateway-profiles", func(ctx *gin.Context) {
		p, err := inst.REST.GetGatewayProfiles()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/cs/devices", func(ctx *gin.Context) {
		p, err := inst.REST.GetDevices()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/cs/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		p, err := inst.REST.GetDevice(eui)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	// add device
	mux.POST("/cs/devices", func(ctx *gin.Context) {
		device, err := bodyDevice(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		p, err := inst.REST.AddDevice(device)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	// edit device
	mux.PUT("/cs/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		device, err := bodyDevice(ctx)
		device.Device.ApplicationID = "1"
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		p, err := inst.REST.EditDevice(eui, device)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	// delete device
	mux.DELETE("/cs/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		p, err := inst.REST.DeleteDevice(eui)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	// DeviceOTAKeysUpdate
	mux.PUT("/cs/devices/keys/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		device, err := bodyDeviceKey(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		p, err := inst.REST.DeviceOTAKeysUpdate(eui, device)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.POST("/cs/devices/keys/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		device, err := bodyDeviceKey(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		p, err := inst.REST.DeviceOTAKeys(eui, device)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	// // activate device
	mux.PUT("/cs/devices/activate/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		device, err := bodyActivateDevice(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		p, err := inst.REST.ActivateDevice(eui, device)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

}
