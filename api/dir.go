package api

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
	"os"
)

type DirApi struct {
	FileMode int
}

func (inst *DirApi) DirExists(c *gin.Context) {
	path := c.Query("path")
	exists := fileutils.DirExists(path)
	dirExistence := interfaces.DirExistence{Path: path, Exists: exists}
	ResponseHandler(dirExistence, nil, c)
}

func (inst *DirApi) CreateDir(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		ResponseHandler(nil, errors.New("path can not be empty"), c)
		return
	}
	err := os.MkdirAll(path, os.FileMode(inst.FileMode))
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("created directory: %s", path)}, err, c)
}
