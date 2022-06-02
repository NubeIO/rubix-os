package api

import (
	"fmt"
	"github.com/NubeIO/flow-framework/src/system/host"
	"github.com/NubeIO/flow-framework/src/system/ufw"
	"github.com/NubeIO/flow-framework/src/utilstime"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/networking/networking"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/networking/portscanner"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type DeviceInfoDatabase interface {
	GetDeviceInfo() (*model.DeviceInfo, error)
}

type DeviceInfoAPI struct {
	DB   DeviceInfoDatabase
	Nets networking.Nets
}

func (inst *DeviceInfoAPI) GetDeviceInfo(ctx *gin.Context) {
	q, err := inst.DB.GetDeviceInfo()
	ResponseHandler(q, err, ctx)
}

func (inst *DeviceInfoAPI) GetSystemTime(ctx *gin.Context) {
	t := utilstime.SystemTime()
	ResponseHandler(t, nil, ctx)
}

func (inst *DeviceInfoAPI) GetExternalIP(ctx *gin.Context) {
	t, err := inst.Nets.GetInternetIP()
	ResponseHandler(t, err, ctx)
}

func (inst *DeviceInfoAPI) GetNetworks(ctx *gin.Context) {
	t, err := inst.Nets.GetNetworks()
	ResponseHandler(t, err, ctx)
}

func (inst *DeviceInfoAPI) GetInterfacesNames(ctx *gin.Context) {
	t, err := inst.Nets.GetInterfacesNames()
	ResponseHandler(t, err, ctx)
}

func (inst *DeviceInfoAPI) GetOSDetails(ctx *gin.Context) {
	out := host.GetCombinationData(false)
	ResponseHandler(out, nil, ctx)
}

func (inst *DeviceInfoAPI) GetTZoneList(ctx *gin.Context) {
	out, err := utilstime.GetTimeZoneList()
	ResponseHandler(out, err, ctx)
}

func (inst *DeviceInfoAPI) FirewallStatus(ctx *gin.Context) {
	out, err := ufw.FirewallStatus()
	ResponseHandler(out, err, ctx)
}

type Ping struct {
	Ip            string `json:"ip,omitempty"`
	InterfaceName string `json:"interface_name,omitempty"`
}

func (inst *DeviceInfoAPI) RubixNetworkPing(ctx *gin.Context) {
	body, _ := getRubixPingDevice(ctx)
	ip := body.Ip
	interfaceName := body.InterfaceName
	if ip == "" {
		if interfaceName == "" { // if nothing is provided then take a guess of the user network
			gateway, err := inst.Nets.GetNetworksThatHaveGateway()
			if err != nil {
				ResponseHandler(nil, err, ctx)
				return
			}
			for i, net := range gateway {
				if i == 0 {
					interfaceName = net.Interface
				}
			}
			net, err := inst.Nets.GetNetworkByIface(interfaceName)
			if err != nil {
				ResponseHandler(nil, err, ctx)
			}
			ip = net.Gateway
		}
	}

	ip = fmt.Sprintf("%s-254", ip)
	ports := []string{"22", "1313", "1414", "1616", "1615"}

	// IP sequence is defined by a '-' between first and last IP address .
	ipsSequence := []string{ip}

	// result returns a map with open ports for each IP address.
	results := portscanner.IPScanner(ipsSequence, ports, true)
	ResponseHandler(results, nil, ctx)
}
