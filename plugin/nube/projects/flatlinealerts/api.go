package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/plugin"
	modbschema "github.com/NubeIO/rubix-os/schema/modbuschema"
	"github.com/NubeIO/rubix-os/utils/array"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
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
	Network *model.Network
	Device  *model.Device
	Point   *model.Point
	// Operation     `json:"request_body"`
	Scan          `json:"scan"`
	ReturnArray   bool  `json:"return_array"`
	IsSerial      bool  `json:"is_serial"`
	DeviceAddress uint8 `json:"device_address"`
}

func bodyClient(ctx *gin.Context) (dto Body, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

type tmvDeviceJSON struct {
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	SerialPort    string `json:"serial_port"`
	BaudRate      uint   `json:"baud_rate"`
	DeviceAddr    uint   `json:"device_addr"`
	WizardVersion uint   `json:"wizard_version"`
	NameArg       string `json:"name_arg"`
	AddArg        uint   `json:"add_arg"`
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
	inst.flatlinealertsDebugMsg(fmt.Sprintf("RegisterWebhook(): %+v\n", inst))
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
	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, modbschema.GetNetworkSchema())
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, modbschema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, modbschema.GetPointSchema())
	})
}
