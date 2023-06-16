package api

import (
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/lib-systemctl-go/systemctl/properties"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type SystemctlAPI struct {
	SystemCtl *systemctl.SystemCtl
}

func (a *SystemctlAPI) SystemCtlEnable(c *gin.Context) {
	unit := c.Query("unit")
	err := a.SystemCtl.Enable(unit)
	message := interfaces.Message{Message: "enabled successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlDisable(c *gin.Context) {
	unit := c.Query("unit")
	err := a.SystemCtl.Disable(unit)
	message := interfaces.Message{Message: "disabled successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlShow(c *gin.Context) {
	unit := c.Query("unit")
	p := c.Query("property")
	property, err := a.SystemCtl.Show(unit, properties.Property(p))
	property_ := interfaces.SystemCtlProperty{Property: property}
	ResponseHandler(property_, err, c)
}

func (a *SystemctlAPI) SystemCtlStart(c *gin.Context) {
	unit := c.Query("unit")
	err := a.SystemCtl.Start(unit)
	message := interfaces.Message{Message: "started successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlStatus(c *gin.Context) {
	unit := c.Query("unit")
	status, err := a.SystemCtl.Status(unit)
	message := interfaces.SystemCtlStatus{Status: status}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlStop(c *gin.Context) {
	unit := c.Query("unit")
	err := a.SystemCtl.Stop(unit)
	message := interfaces.Message{Message: "stopped successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlResetFailed(c *gin.Context) {
	err := a.SystemCtl.RestartFailed()
	message := interfaces.Message{Message: "reset-failed command executed successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlDaemonReload(c *gin.Context) {
	err := a.SystemCtl.DaemonReload()
	message := interfaces.Message{Message: "daemon-reload command executed successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlRestart(c *gin.Context) {
	unit := c.Query("unit")
	err := a.SystemCtl.Restart(unit)
	message := interfaces.Message{Message: "restarted successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlMask(c *gin.Context) {
	unit := c.Query("unit")
	err := a.SystemCtl.Mask(unit)
	message := interfaces.Message{Message: "masked successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlUnmask(c *gin.Context) {
	unit := c.Query("unit")
	err := a.SystemCtl.Unmask(unit)
	message := interfaces.Message{Message: "unmasked successfully"}
	ResponseHandler(message, err, c)
}

func (a *SystemctlAPI) SystemCtlState(c *gin.Context) {
	unit := c.Query("unit")
	state, err := a.SystemCtl.State(unit)
	ResponseHandler(state, err, c)
}

func (a *SystemctlAPI) SystemCtlIsEnabled(c *gin.Context) {
	unit := c.Query("unit")
	state, err := a.SystemCtl.IsEnabled(unit)
	state_ := interfaces.SystemCtlState{State: state}
	ResponseHandler(state_, err, c)
}

func (a *SystemctlAPI) SystemCtlIsActive(c *gin.Context) {
	unit := c.Query("unit")
	state, status, err := a.SystemCtl.IsActive(unit)
	status_ := interfaces.SystemCtlStateStatus{State: state, Status: status}
	ResponseHandler(status_, err, c)
}

func (a *SystemctlAPI) SystemCtlIsRunning(c *gin.Context) {
	unit := c.Query("unit")
	state, status, err := a.SystemCtl.IsRunning(unit)
	status_ := interfaces.SystemCtlStateStatus{State: state, Status: status}
	ResponseHandler(status_, err, c)
}

func (a *SystemctlAPI) SystemCtlIsFailed(c *gin.Context) {
	unit := c.Query("unit")
	state, err := a.SystemCtl.IsFailed(unit)
	state_ := interfaces.SystemCtlState{State: state}
	ResponseHandler(state_, err, c)
}

func (a *SystemctlAPI) SystemCtlIsInstalled(c *gin.Context) {
	unit := c.Query("unit")
	state, err := a.SystemCtl.IsInstalled(unit)
	state_ := interfaces.SystemCtlState{State: state}
	ResponseHandler(state_, err, c)
}
