package api

import (
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
	"syscall"
)

type SyscallAPI struct {
}

func (a *SyscallAPI) SyscallUnlink(c *gin.Context) {
	path := c.Query("path")
	err := syscall.Unlink(path)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(interfaces.Message{Message: "unlinked successfully"}, err, c)
}

func (a *SyscallAPI) SyscallLink(c *gin.Context) {
	path := c.Query("path")
	link := c.Query("link")
	err := syscall.Symlink(path, link)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(interfaces.Message{Message: "linked successfully"}, err, c)
}
