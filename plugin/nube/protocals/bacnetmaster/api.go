package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/system/networking"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type Network struct {
	NetworkUUID           string `json:"network_uuid"`
	NetworkName           string `json:"network_name"`
	NetworkIp             string `json:"network_ip"`
	NetworkMask           int    `json:"network_mask"`
	NetworkPort           int    `json:"network_port"`
	NetworkDeviceObjectId int    `json:"network_device_object_id"`
	NetworkDeviceName     string `json:"network_device_name"`
	InterfaceName         string `json:"interface_name"`
}

type Networks struct {
	Networks []struct {
		NetworkUuid string `json:"network_uuid"`
	} `json:"networks"`
}

type Device struct {
	DeviceUUID     string `json:"device_uuid"`
	DeviceName     string `json:"device_name"`
	DeviceIp       string `json:"device_ip"`
	DeviceEnable   bool   `json:"device_enable"`
	DeviceMask     int    `json:"device_mask"`
	DevicePort     int    `json:"device_port"`
	DeviceMac      int    `json:"device_mac"`
	DeviceObjectId int    `json:"device_object_id"`
	NetworkNumber  int    `json:"network_number"`
	TypeMstp       bool   `json:"type_mstp"`
	SupportsRpm    bool   `json:"supports_rpm"`
	SupportsWpm    bool   `json:"supports_wpm"`
	NetworkUuid    string `json:"network_uuid"`
}

type Point struct {
	PointName       string `json:"point_name"`
	PointEnable     bool   `json:"point_enable"`
	PointObjectId   int    `json:"point_object_id"`
	PointObjectType string `json:"point_object_type"`
	DeviceUuid      string `json:"device_uuid"`
}

type Whois struct {
	NetworkNumber         int  `json:"network_number"`
	IsMstp                bool `json:"is_mstp"`
	Whois                 bool `json:"whois"`
	GlobalBroadcast       bool `json:"global_broadcast"`
	FullRange             bool `json:"full_range"`
	RangeStart            int  `json:"range_start"`
	RangeEnd              int  `json:"range_end"`
	ShowSupportedServices bool `json:"show_supported_services"`
	AddDevices            bool `json:"add_devices"`
}

type ReleaseOverride struct {
	Priority int  `json:"priority"`
	Feedback bool `json:"feedback"`
	Timeout  int  `json:"timeout"`
}

type PointWrite struct {
	PointWriteValue int  `json:"point_write_value"`
	Priority        int  `json:"priority"`
	Feedback        bool `json:"feedback"`
	Timeout         int  `json:"timeout"`
}

type PointRead struct {
	GetPriority bool `json:"get_priority"`
	Timeout     int  `json:"timeout"`
}

type PointsDiscover struct {
	AddPoints bool `json:"add_points"`
	Discovery bool `json:"discovery"`
	Timeout   int  `json:"timeout"`
}

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"

	wizard         = "/wizard"
	ping           = "/ping"
	pingBacnet     = "/api/system/ping"
	point          = "/point"
	pointBacnet    = "/api/bm/point"
	points         = "/points"
	pointsBacnet   = "/api/bm/points"
	network        = "/network"
	networkBacnet  = "/api/bm/network"
	networks       = "/networks"
	networksBacnet = "/api/bm/networks/true"
	device         = "/device"
	deviceBacnet   = "/api/bm/device"
)

func getBODYAddNetwork(ctx *gin.Context) (dto *Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveObject(ctx *gin.Context) string {
	return ctx.Param("object")
}

func resolveAddress(ctx *gin.Context) string {
	return ctx.Param("address")
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.POST(wizard, func(ctx *gin.Context) {
		body, err := getBODYAddNetwork(ctx)
		_, err = i.wizard(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, "ok")
			return
		}
	})
	mux.GET(ping, func(ctx *gin.Context) {
		rt.Method = nrest.GET
		rt.Path = pingBacnet
		res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		if err != nil {
			ctx.JSON(code, err)
			return
		} else {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		}
	})
	mux.GET(points, func(ctx *gin.Context) {
		rt.Method = nrest.GET
		rt.Path = pointsBacnet
		res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		if err != nil {
			ctx.JSON(code, err)
			return
		} else {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		}
	})
	mux.GET(networks, func(ctx *gin.Context) {
		rt.Method = nrest.GET
		rt.Path = networksBacnet
		res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		if err != nil {
			ctx.JSON(code, err)
			return
		} else {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		}
	})
	mux.POST(network, func(ctx *gin.Context) {
		body, _ := getBODYAddNetwork(ctx)
		if body.InterfaceName != "" {
			_net, _ := networking.GetInterfaceByName(body.InterfaceName)
			if _net == nil {
				ctx.JSON(http.StatusBadRequest, errors.New("failed to find interface"))
				return
			}
			body.NetworkIp = _net.IP
			body.NetworkMask = _net.NetMaskLength
		}
		rt.Method = nrest.PUT
		rt.Path = networkBacnet
		res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: body})
		if err != nil {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		} else {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		}
	})
	mux.DELETE(networks, func(ctx *gin.Context) {
		//getNetworks
		rt.Method = nrest.GET
		rt.Path = networksBacnet
		res, _, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		nets := new(Networks)
		err = res.ToInterface(nets)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errors.New("failed to find get networks"))
		}
		count := 0
		for _, a := range nets.Networks {
			rt.Method = nrest.DELETE
			rt.Path = networkBacnet + "/" + a.NetworkUuid
			fmt.Println(rt.Path)
			count++
			nrest.DoHTTPReq(rt, &nrest.ReqOpt{})

		}
		r := fmt.Sprintf("deleted %d number of networks", count)
		ctx.JSON(http.StatusOK, r)
	})
}
