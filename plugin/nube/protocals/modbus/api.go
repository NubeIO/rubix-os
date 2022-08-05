package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/jsonschema"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/mbmodel"
	"github.com/NubeIO/flow-framework/utils/array"
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
	mux.GET("/modbus/polling/stats", func(ctx *gin.Context) {
		stats, err := inst.getPollingStats()
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
		ctx.JSON(http.StatusOK, jsonschema.GetNetworkSchema())
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, jsonschema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, jsonschema.GetPointSchema())
	})
}

// wizard make a network/dev/pnt
func (inst *Instance) getPollingStats() (result []interface{}, error error) {
	if len(inst.NetworkPollManagers) == 0 {
		return nil, nil
	}
	type Stats struct {
		NetworkName                   string
		MaxPollExecuteTimeSecs        float64 // time in seconds for polling to complete (poll response time, doesn't include the time in queue).
		AveragePollExecuteTimeSecs    float64 // time in seconds for polling to complete (poll response time, doesn't include the time in queue).
		MinPollExecuteTimeSecs        float64 // time in seconds for polling to complete (poll response time, doesn't include the time in queue).
		TotalPollQueueLength          int64   // number of polling points in the current queue.
		TotalStandbyPointsLength      int64   // number of polling points in the standby list.
		TotalPointsOutForPolling      int64   // number of points currently out for polling (currently being handled by the protocol plugin).
		ASAPPriorityPollQueueLength   int64   // number of ASAP priority polling points in the current queue.
		HighPriorityPollQueueLength   int64   // number of High priority polling points in the current queue.
		NormalPriorityPollQueueLength int64   // number of Normal priority polling points in the current queue.
		LowPriorityPollQueueLength    int64   // number of Low priority polling points in the current queue.
		ASAPPriorityAveragePollTime   float64 // average time in seconds between ASAP priority polling point added to current queue, and polling complete.
		HighPriorityAveragePollTime   float64 // average time in seconds between High priority polling point added to current queue, and polling complete.
		NormalPriorityAveragePollTime float64 // average time in seconds between Normal priority polling point added to current queue, and polling complete.
		LowPriorityAveragePollTime    float64 // average time in seconds between Low priority polling point added to current queue, and polling complete.
		TotalPollCount                int64   // total number of polls completed.
		ASAPPriorityPollCount         int64   // total number of ASAP priority polls completed.
		HighPriorityPollCount         int64   // total number of High priority polls completed.
		NormalPriorityPollCount       int64   // total number of Normal priority polls completed.
		LowPriorityPollCount          int64   // total number of Low priority polls completed.
		ASAPPriorityMaxCycleTime      float64 // threshold setting for triggering a lockup alert for ASAP priority.
		HighPriorityMaxCycleTime      float64 // threshold setting for triggering a lockup alert for High priority.
		NormalPriorityMaxCycleTime    float64 // threshold setting for triggering a lockup alert for Normal priority.
		LowPriorityMaxCycleTime       float64 // threshold setting for triggering a lockup alert for Low priority.
		ASAPPriorityLockupAlert       bool    // alert if poll time has exceeded the ASAPPriorityMaxCycleTime
		HighPriorityLockupAlert       bool    // alert if poll time has exceeded the HighPriorityMaxCycleTime
		NormalPriorityLockupAlert     bool    // alert if poll time has exceeded the NormalPriorityMaxCycleTime
		LowPriorityLockupAlert        bool    // alert if poll time has exceeded the LowPriorityMaxCycleTime
		BusyTime                      float64 // percent of the time that the plugin is actively polling.
		EnabledTime                   float64 // time in seconds that the statistics have been running for.
		PortUnavailableTime           float64 // time in seconds that the serial port has been unavailable.
	}

	for _, netPollMan := range inst.NetworkPollManagers {
		var netArg api.Args
		net, _ := inst.db.GetNetwork(netPollMan.FFNetworkUUID, netArg)
		pmResult := Stats{
			NetworkName:                   net.Name,
			MaxPollExecuteTimeSecs:        netPollMan.MaxPollExecuteTimeSecs,
			AveragePollExecuteTimeSecs:    netPollMan.AveragePollExecuteTimeSecs,
			MinPollExecuteTimeSecs:        netPollMan.MinPollExecuteTimeSecs,
			TotalPollQueueLength:          netPollMan.TotalPollQueueLength,
			TotalStandbyPointsLength:      netPollMan.TotalStandbyPointsLength,
			TotalPointsOutForPolling:      netPollMan.TotalPointsOutForPolling,
			ASAPPriorityPollQueueLength:   netPollMan.ASAPPriorityPollQueueLength,
			HighPriorityPollQueueLength:   netPollMan.HighPriorityPollQueueLength,
			NormalPriorityPollQueueLength: netPollMan.NormalPriorityPollQueueLength,
			LowPriorityPollQueueLength:    netPollMan.LowPriorityPollQueueLength,
			ASAPPriorityAveragePollTime:   netPollMan.ASAPPriorityAveragePollTime,
			HighPriorityAveragePollTime:   netPollMan.HighPriorityAveragePollTime,
			NormalPriorityAveragePollTime: netPollMan.NormalPriorityAveragePollTime,
			LowPriorityAveragePollTime:    netPollMan.LowPriorityAveragePollTime,
			TotalPollCount:                netPollMan.TotalPollCount,
			ASAPPriorityPollCount:         netPollMan.ASAPPriorityPollCount,
			HighPriorityPollCount:         netPollMan.HighPriorityPollCount,
			NormalPriorityPollCount:       netPollMan.NormalPriorityPollCount,
			LowPriorityPollCount:          netPollMan.LowPriorityPollCount,
			ASAPPriorityMaxCycleTime:      netPollMan.ASAPPriorityMaxCycleTime.Seconds(),
			HighPriorityMaxCycleTime:      netPollMan.HighPriorityMaxCycleTime.Seconds(),
			NormalPriorityMaxCycleTime:    netPollMan.NormalPriorityMaxCycleTime.Seconds(),
			LowPriorityMaxCycleTime:       netPollMan.LowPriorityMaxCycleTime.Seconds(),
			ASAPPriorityLockupAlert:       netPollMan.ASAPPriorityLockupAlert,
			HighPriorityLockupAlert:       netPollMan.HighPriorityLockupAlert,
			NormalPriorityLockupAlert:     netPollMan.NormalPriorityLockupAlert,
			LowPriorityLockupAlert:        netPollMan.LowPriorityLockupAlert,
			BusyTime:                      netPollMan.BusyTime,
			EnabledTime:                   netPollMan.EnabledTime,
			PortUnavailableTime:           netPollMan.PortUnavailableTime,
		}
		result = append(result, pmResult)
	}
	return result, nil
}
