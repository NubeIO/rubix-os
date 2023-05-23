package api

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

type ZipApi struct {
	FileMode int
}

func (a *ZipApi) Unzip(c *gin.Context) {
	source := c.Query("source")
	destination := c.Query("destination")
	pathToZip := source
	if source == "" {
		ResponseHandler(nil, errors.New("zip source can not be empty, try /data/zip.zip"), c)
		return
	}
	if destination == "" {
		ResponseHandler(nil, errors.New("zip destination can not be empty, try /data/unzip-test"), c)
		return
	}
	zip, err := fileutils.Unzip(pathToZip, destination, os.FileMode(a.FileMode))
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(zip, err, c)
}

func (a *ZipApi) ZipDir(c *gin.Context) {
	source := c.Query("source")
	destination := c.Query("destination")
	pathToZip := source
	if source == "" {
		ResponseHandler(nil, errors.New("zip source can not be empty, try /data/rubix-os"), c)
		return
	}
	if destination == "" {
		ResponseHandler(nil, errors.New("zip destination can not be empty, try /data/test/rubix-os.zip"), c)
		return
	}
	exists := fileutils.DirExists(pathToZip)
	if !exists {
		ResponseHandler(nil, errors.New("zip source is not found"), c)
		return
	}
	err := os.MkdirAll(filepath.Dir(destination), os.FileMode(a.FileMode))
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	err = fileutils.RecursiveZip(pathToZip, destination)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("zip file is created on: %s", destination)}, nil, c)
}
