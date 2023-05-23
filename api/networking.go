package api

import (
	"github.com/NubeIO/lib-dhcpd/dhcpd"
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/rubix-os/services/system"
	"github.com/gin-gonic/gin"
)

var nets = networking.New()

type NetworkingAPI struct {
	System *system.System
}

func (a *NetworkingAPI) Networking(c *gin.Context) {
	data, err := nets.GetNetworks()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) GetInterfacesNames(c *gin.Context) {
	data, err := nets.GetInterfacesNames()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) InternetIP(c *gin.Context) {
	data, err := nets.GetInternetIP()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) RestartNetworking(c *gin.Context) {
	data, err := a.System.RestartNetworking()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) InterfaceUpDown(c *gin.Context) {
	var m system.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.InterfaceUpDown(m)
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) InterfaceUp(c *gin.Context) {
	var m system.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.InterfaceUp(m)
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) InterfaceDown(c *gin.Context) {
	var m system.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.InterfaceDown(m)
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) DHCPPortExists(c *gin.Context) {
	var m system.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.DHCPPortExists(m)
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) DHCPSetAsAuto(c *gin.Context) {
	var m system.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.DHCPSetAsAuto(m)
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) DHCPSetStaticIP(c *gin.Context) {
	var m *dhcpd.SetStaticIP
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.DHCPSetStaticIP(m)
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) UWFActive(c *gin.Context) {
	data, err := a.System.UWFActive()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) UWFEnable(c *gin.Context) {
	data, err := a.System.UWFEnable()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) UWFDisable(c *gin.Context) {
	data, err := a.System.UWFDisable()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) UWFStatus(c *gin.Context) {
	data, err := a.System.UWFStatus()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) UWFStatusList(c *gin.Context) {
	data, err := a.System.UWFStatusList()
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) UWFOpenPort(c *gin.Context) {
	var m system.UFWBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.UWFOpenPort(m)
	ResponseHandler(data, err, c)
}

func (a *NetworkingAPI) UWFClosePort(c *gin.Context) {
	var m system.UFWBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.UWFClosePort(m)
	ResponseHandler(data, err, c)
}
