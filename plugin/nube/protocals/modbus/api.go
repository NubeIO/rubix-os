package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/mbmodel"
	"github.com/NubeIO/flow-framework/utils/array"
	modbschema "github.com/NubeIO/lib-schema/modbuschema"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uurl"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	help              = "/help"
	listSerial        = "/list/serial"
	schemaNetwork     = "/schema/network"
	schemaDevice      = "/schema/device"
	schemaPoint       = "/schema/point"
	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
)

type Scan struct {
	Start  uint32 `json:"start"`
	Count  uint32 `json:"count"`
	IsCoil bool   `json:"is_coil"`
}

type Body struct {
	Network       *model.Network
	Device        *model.Device
	Point         *model.Point
	Client        `json:"client"`
	Operation     `json:"request_body"`
	Scan          `json:"scan"`
	ReturnArray   bool  `json:"return_array"`
	IsSerial      bool  `json:"is_serial"`
	DeviceAddress uint8 `json:"device_address"`
}

func bodyClient(ctx *gin.Context) (dto Body, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// serialWizard
type wizard struct {
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	SerialPort    string `json:"serial_port"`
	BaudRate      uint   `json:"baud_rate"`
	DeviceAddr    uint   `json:"device_addr"`
	WizardVersion uint   `json:"wizard_version"`
	NameArg       string `json:"name_arg"`
	AddArg        uint   `json:"add_arg"`
}

func bodyWizard(ctx *gin.Context) (dto wizard, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// supportedObjects return all objects that are not bacnet
func supportedObjects() *array.Array {
	out := array.NewArray()
	out.Add(model.ObjTypeAnalogInput)
	out.Add(model.ObjTypeAnalogOutput)
	out.Add(model.ObjTypeAnalogValue)
	out.Add(model.ObjTypeBinaryInput)
	out.Add(model.ObjTypeBinaryOutput)
	out.Add(model.ObjTypeBinaryValue)
	return out
}

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	inst.modbusDebugMsg(fmt.Sprintf("RegisterWebhook(): %+v\n", inst))
	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.addNetwork(body)
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

	mux.GET(listSerial, func(ctx *gin.Context) {
		serial, err := inst.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, serial)
			return
		}
	})
	mux.POST("/modbus/point/operation", func(ctx *gin.Context) {
		body, err := bodyClient(ctx)
		netType := body.Network.TransportType
		mbClient, err := inst.setClient(body.Network, body.Device, false)
		if err != nil {
			inst.modbusErrorMsg(err, "ERROR ON set modbus client")
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if netType == model.TransType.Serial || netType == model.TransType.LoRa {
			if body.Device.AddressId >= 1 {
				mbClient.RTUClientHandler.SlaveID = byte(body.Device.AddressId)
			}
		} else if netType == model.TransType.IP {
			url, err := uurl.JoinIpPort(body.Device.Host, body.Device.Port)
			if err != nil {
				inst.modbusErrorMsg(fmt.Sprintf("failed to validate device IP %s\n", url))
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
			mbClient.TCPClientHandler.Address = url
			mbClient.TCPClientHandler.SlaveID = byte(body.Device.AddressId)
		}
		_, responseValue, err := inst.networkRequest(mbClient, body.Point, false)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		inst.modbusDebugMsg("responseValue", responseValue)
		ctx.JSON(http.StatusOK, responseValue)
		return
	})
	mux.POST("/modbus/wizard/tcp", func(ctx *gin.Context) {
		body, _ := bodyWizard(ctx)
		n, err := inst.wizardTCP(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, n)
			return
		}
	})
	mux.POST("/modbus/wizard/serial", func(ctx *gin.Context) {
		body, _ := bodyWizard(ctx)
		serial, err := inst.wizardSerial(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, serial)
			return
		}
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
		ctx.JSON(http.StatusOK, mbmodel.GetNetworkSchema())
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mbmodel.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mbmodel.GetPointSchema())
	})
	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, modbschema.GetNetworkSchema())
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		fns, err := inst.db.GetFlowNetworks(api.Args{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		deviceSchema := modbschema.GetDeviceSchema()
		deviceSchema.AutoMappingFlowNetworkName.Options = plugin.GetFlowNetworkNames(fns)
		ctx.JSON(http.StatusOK, deviceSchema)
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, modbschema.GetPointSchema())
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
